# AWS frontend hosting

This module contains the independent AWS CDK implementation for the repository's static frontend hosting. It was adapted from [`GWC-Hunter-College/Website-Hosting-Iac`](https://github.com/GWC-Hunter-College/Website-Hosting-Iac) at commit `21152e45112f0fe925b57798675650f1027828a7` and checked against the matching frontend stacks in the preserved Event Management System implementation.

## Current stacks

The Go CDK application synthesizes three environment-agnostic stacks:

- `GirlsWhoCodeDomainStack` owns the retained public Route 53 hosted zone and outputs its generated `HostedZoneId` and `NameServers`.
- Girls Who Code at Hunter uses `NewGirlsWhoCodeHostingStack` and `GirlsWhoCodeHostingStackProps`; its historical CDK/CloudFormation stack ID remains `FrontendStack`. It defines the fixed-name `gwc-club-site` bucket and the `gwc-website-ci-deployer` user with the `frontend-gwc-ci-policy` policy.
- Hunter College Clubs / Event Manager uses `NewHunterCollegeClubsHostingStack` and `HunterCollegeClubsHostingStackProps`; its historical CDK/CloudFormation stack ID remains `FrontendHccStack`. It defines the fixed-name `hunter-college-club-event-site` bucket and the `hcc-website-ci-deployer` user with the `frontend-hcc-ci-policy` policy.

## What each stack creates

Both sites use the same independent hosting and deployment model:

- **Private S3 bucket:** stores static build artifacts under separate `/staging` and `/production` prefixes.
- **CloudFront origin access control (OAC):** signs CloudFront requests to the private S3 bucket with SigV4, so the bucket can remain private. Each stack uses one OAC for its staging and production origins.
- **Staging and production CloudFront distributions:** serve their matching S3 prefixes, redirect viewers to HTTPS, and use `index.html` as the default root object.
- **SPA error responses:** translate S3 403 and 404 responses to `/index.html` with status 200 so client-side routes can load directly.
- **CI deployment user and policy:** can upload or delete objects in only that site's bucket, list that bucket, and invalidate only that site's two distributions. This deployment identity is independent of OAC. Neither stack creates an IAM access key or outputs credentials.
- **Outputs:** expose the bucket name, CloudFront URLs and distribution IDs, and staging and production S3 deployment destinations.

Girls Who Code production also uses an ACM certificate and Route 53 A/AAAA aliases for the apex and `www` names. Both names serve the existing `FrontendProduction` distribution and `/production` content; `FrontendMain` keeps its generated CloudFront URL. Hunter College Clubs is unchanged. Frontend delivery automation must build and upload assets to the documented prefixes separately; this module does not create frontend build artifacts or an S3 deployment construct.

The reference repository used `/main` for the Girls Who Code non-production origin. This module intentionally uses `/staging` to match this project's integration branch while retaining the existing `FrontendMain` construct ID and `CloudFront_Main_Info` output ID. Hunter College Clubs uses the corresponding `FrontendStaging` and `CloudFront_Staging_Info` IDs. Existing Hunter College Clubs IAM/output IDs also remain unchanged, while the new Girls Who Code equivalents use explicit application names. These identity differences are intentional; the hosting behavior is parallel. Coordinate the uploader so assets exist under `/staging`, and review the CloudFront origin-path change before any deployment.

## Optional domains and resource ownership

`GirlsWhoCodeHostingStackProps.DomainName` and `HostedZone` are optional as a
pair. Omitting both (including passing `nil` props) creates core hosting only:
the S3 bucket, OAC, staging/production CloudFront distributions, generated
CloudFront URL outputs, and deployment IAM resources. That configuration has no
ACM certificate, Route 53 resources, CloudFront custom aliases, CloudFormation
imports, or dependency on a domain stack. Supplying only one domain prop is an
error so a partial configuration cannot silently detach DNS.

The current `main.go` deliberately still supplies both values for the live GWC
site. Its dependency remains `FrontendStack → GirlsWhoCodeDomainStack` (arrow
means **depends on**), via an automatic hosted-zone ID export/import used by ACM
DNS validation and the four alias records. There is no explicit dependency or
reverse reference. The `DomainName` prop is a plain configuration string and
does not itself create an import.

The domain-stack implementation already takes a caller-supplied domain name and
stack ID, so it can provision a zone for another site without copying code.
Its historical Go names and internal hosted-zone construct ID are retained;
renaming them adds no capability. Hosting accepts an `IPublicHostedZone`, so a
caller can also supply a reference to an existing zone without creating a domain
stack. For the live GWC zone, keep using the existing managed construct in
`main.go`: do not add a second zone stack or switch its ownership to an import.

The certificate and DNS records stay in their existing `FrontendStack` scopes.
Moving them into the zone stack would change CloudFormation ownership and could
create a circular dependency because the records target the frontend's local
production distribution. Optional domain support does not require that migration.

Hunter College Clubs already uses independent `FrontendHccStack` hosting. Its
current constructor receives no domain configuration and creates no hosted zone,
DNS records, or certificate. There is no Hunter domain stack in `main.go`. A
future real domain can use the same zone provisioning implementation with its
own stack identity and explicit hosting integration; no Hunter domain attachment
or placeholder domain is created by this refactor.

The domainless constructor path is for sites configured without a domain. **Do
not remove the existing GWC domain props or domain-stack instance from main.go**:
that would request removal of its working certificate, aliases, and DNS records.
Optional props do not make detaching a live domain a harmless operation. Likewise,
the GWC hosting constructor retains GWC's fixed physical bucket/IAM names; use
the dedicated Hunter constructor for Hunter rather than deploying a second copy
of GWC hosting.

The deployed identities remain:

| Stack | Resource | Logical ID |
| --- | --- | --- |
| `GirlsWhoCodeDomainStack` | Public zone | `GirlsWhoCodeHostedZoneDBA9EB09` |
| `FrontendStack` | S3 bucket | `GwcWebsiteBucketE6A54810` |
| `FrontendStack` | Staging distribution | `FrontendMain4FAF8302` |
| `FrontendStack` | Production distribution | `FrontendProduction57D7F36D` |
| `FrontendStack` | Certificate | `GirlsWhoCodeProductionCertificate70C3275C` |

The four alias-record IDs, zone export, stack environments, and deployment IAM
identities are also preserved. The zone keeps both `DeletionPolicy: Retain` and
`UpdateReplacePolicy: Retain`. Existing bucket auto-delete and IAM user delete
policies remain unchanged; this refactor does not add retention to those resources.

Validation covers GWC with a domain, GWC without a domain (nil and stack-only
props), domainless Hunter hosting, and rejection of partial domain configuration.
For the unchanged live configuration, synthesis should produce the same templates
and `cdk diff GirlsWhoCodeDomainStack FrontendStack --no-change-set` should report
no differences. This diff reads deployed templates without creating change sets.

## Validate locally

From this directory:

```sh
go mod download
go test ./...
go vet ./...
cdk synth --all
```

`cdk.json` runs the Go application for CDK. Synthesis produces CloudFormation templates locally; it does not deploy them.

## Deployment safety and ownership

Girls Who Code `FrontendStack` is already deployed. Deploy its updates to that same account and region to preserve resource ownership. Any separate Hunter College Clubs deployment still requires its own ownership review.

- The extracted stacks deliberately retain the historical stack IDs, construct IDs, and fixed physical bucket names. Hunter College Clubs also retains its existing IAM names; the new fixed Girls Who Code IAM names must be checked for account-level conflicts before deployment. Preserved template identity helps only when the same existing CloudFormation stacks remain the owners. Different stack ownership can cause fixed-name conflicts or replacement/duplicate CloudFront resources.
- Do not deploy the hosting stacks from both this module and `infrastructure/legacy`. Existing resources may need to remain under their current stacks or be explicitly imported/adopted before this module becomes authoritative.
- Both buckets use `RemovalPolicy_DESTROY` with automatic object deletion. Stack removal or a replacement that deletes the old bucket can permanently delete hosted assets.
- Both named CI users are real `AWS::IAM::User` resources with an explicit destroy policy. If either physical user already exists, import that user into its stack before deployment; replacing it with `User_FromUserName` would leave it externally owned, while a normal deploy would fail by trying to create a duplicate.
- The current Girls Who Code user predated `FrontendStack` and was adopted without rotating its active, externally created access key or removing its historical `GWCWebsiteCICDPolicy` attachment. The CDK application does not define or expose that key. The Hunter College Clubs user did not exist when checked and remains a normal stack-created resource.
- Destroying a stack-managed CI user is credential revocation: the current CloudFormation IAM user delete handler can delete attached access keys and detach managed policies before deleting the user. Do not destroy either hosting stack until its deployment credentials are intentionally retired or replaced.
- The implementation uses CloudFront OAC with `S3BucketOrigin`. Before the first deployment, verify that each generated bucket policy grants `cloudfront.amazonaws.com` read access only for that site's two distribution ARNs.
- The application is environment-agnostic, as in the reference repository. Review the intended AWS account, region, stack names, and synthesized change set before any eventual deployment.

The earlier IAM-only adoption has been followed by a Girls Who Code hosting deployment. The domain change updates that existing stack; it does not recreate or migrate its hosting resources.


## Girls Who Code custom domain: two-stage deployment

The single application configuration value is `girlsWhoCodeDomainName` in `main.go` (`girlswhocodehunter.org`). It is passed through stack props. Generated nameservers, hosted zone IDs, certificate ARNs, and distribution IDs are resolved by CDK constructs and CloudFormation outputs, never stored in configuration or a `.env` file.

`FrontendStack` depends on `GirlsWhoCodeDomainStack` through the hosted-zone construct (automatic CloudFormation export/import). The domain stack owns only the long-lived public zone, retained on stack deletion. The frontend stack owns the DNS-validated ACM certificate and all four production alias records, which reference its existing distribution. There is no reverse dependency. A retained zone would need to be imported before recreating its owning stack.

### Confirm the deployment environment

The application still returns `nil` from `environment()`; it has not been moved to a different account or region. Before implementation, AWS `describe-stacks` confirmed the existing `FrontendStack` in `us-east-1` with status `UPDATE_COMPLETE`. Use the same AWS profile/account for both phases. The default profile was used for that check; substitute your existing deployment profile if needed.

```sh
cd hosting/aws
export AWS_PROFILE=default
export AWS_REGION=us-east-1
export AWS_DEFAULT_REGION=us-east-1
aws sts get-caller-identity
aws cloudformation describe-stacks --stack-name FrontendStack   --region us-east-1 --query 'Stacks[0].[StackId,StackStatus]'
```

Confirm the account and existing stack before proceeding. CloudFront requires its ACM certificate in [us-east-1](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/cnames-and-https-requirements.html). If the intended existing stack is in another region, stop and design a dedicated certificate stack; do not move the hosting stack.

### Phase 1: create the authoritative zone

From `hosting/aws`, with the environment above:

```sh
cdk synth
cdk diff GirlsWhoCodeDomainStack --no-change-set
cdk deploy GirlsWhoCodeDomainStack
aws cloudformation describe-stacks --stack-name GirlsWhoCodeDomainStack   --region us-east-1 --query 'Stacks[0].Outputs' --output table
```

Deploy only the domain stack at this point. Synthesis creates local templates for all stacks but does not request a certificate. Copy the four nameservers from the `NameServers` output.

### Manual Namecheap delegation

In Namecheap, open **Domain List → girlswhocodehunter.org → Manage → Nameservers → Custom DNS**. Enter the four Route 53 nameservers from Phase 1 and save. Do not create Namecheap A records pointing at CloudFront. Wait for delegation to propagate and verify that the returned NS set matches the Route 53 output:

```sh
dig NS girlswhocodehunter.org +short
```

Namecheap remains the registrar and handles registration, renewal, and registrant information. This changes authoritative DNS hosting to Route 53; it is not a registrar transfer.

### Phase 2: connect the existing production distribution

Once delegation points to Route 53, from `hosting/aws` with the same profile and region:

```sh
cdk diff FrontendStack --no-change-set
cdk deploy FrontendStack
```

The diff should add one ACM certificate and four Route 53 alias records, and modify only the aliases and viewer certificate of the existing `FrontendProduction` distribution. Stop if it proposes deleting or replacing an existing bucket, distribution, OAC, or IAM resource. `--no-change-set` performs a read-only template diff without creating a CloudFormation change set; review the deployment change set as well.

This deployment requests the apex/`www` certificate, manages its DNS validation through the hosted zone, updates production HTTPS, and creates A and AAAA aliases targeting production. IPv6 is already enabled. Deploying before delegation propagates can leave ACM validation pending. Both names serve the same content, with no redirect between apex and `www`.

Verify after deployment:

```sh
curl -I https://girlswhocodehunter.org
curl -I https://www.girlswhocodehunter.org
```

Route 53 hosts the DNS records, ACM provides TLS, and CloudFront serves the existing S3 production content. Keep the ACM validation records for certificate renewal.

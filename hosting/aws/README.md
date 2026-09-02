# AWS frontend hosting

This module contains the independent AWS CDK implementation for the repository's static frontend hosting. It was adapted from [`GWC-Hunter-College/Website-Hosting-Iac`](https://github.com/GWC-Hunter-College/Website-Hosting-Iac) at commit `21152e45112f0fe925b57798675650f1027828a7` and checked against the matching frontend stacks in the preserved Event Management System implementation.

## Current stacks

The Go CDK application synthesizes two environment-agnostic stacks:

- Girls Who Code at Hunter uses `NewGirlsWhoCodeHostingStack` and `GirlsWhoCodeHostingStackProps`; its historical CDK/CloudFormation stack ID remains `FrontendStack`. It defines the fixed-name `gwc-club-site` bucket and the `gwc-website-ci-deployer` user with the `frontend-gwc-ci-policy` policy.
- Hunter College Clubs / Event Manager uses `NewHunterCollegeClubsHostingStack` and `HunterCollegeClubsHostingStackProps`; its historical CDK/CloudFormation stack ID remains `FrontendHccStack`. It defines the fixed-name `hunter-college-club-event-site` bucket and the `hcc-website-ci-deployer` user with the `frontend-hcc-ci-policy` policy.

## What each stack creates

Both sites use the same independent hosting and deployment model:

- **Private S3 bucket:** stores static build artifacts under separate `/staging` and `/production` prefixes.
- **CloudFront origin access identity (OAI):** lets CloudFront read the private bucket without making its objects public. Both stacks currently use OAI, not OAC.
- **Staging and production CloudFront distributions:** serve their matching S3 prefixes, redirect viewers to HTTPS, and use `index.html` as the default root object.
- **SPA error responses:** translate S3 403 and 404 responses to `/index.html` with status 200 so client-side routes can load directly.
- **CI deployment user and policy:** can upload or delete objects in only that site's bucket, list that bucket, and invalidate only that site's two distributions. Neither stack creates an IAM access key or outputs credentials.
- **Outputs:** expose the bucket name, CloudFront URLs and distribution IDs, and staging and production S3 deployment destinations.

This implementation defines S3, CloudFront, and a scoped deployment identity for each site. It does **not** define Route 53 records, ACM certificates, custom domains, frontend build artifacts, or an S3 deployment construct. Frontend delivery automation must build and upload assets to the documented prefixes separately.

The reference repository used `/main` for the Girls Who Code non-production origin. This module intentionally uses `/staging` to match this project's integration branch while retaining the existing `FrontendMain` construct ID and `CloudFront_Main_Info` output ID. Hunter College Clubs uses the corresponding `FrontendStaging` and `CloudFront_Staging_Info` IDs. Existing Hunter College Clubs IAM/output IDs also remain unchanged, while the new Girls Who Code equivalents use explicit application names. These identity differences are intentional; the hosting behavior is parallel. Coordinate the uploader so assets exist under `/staging`, and review the CloudFront origin-path change before any deployment.

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

Do not deploy this extracted module until the owners of the existing frontend resources and CloudFormation stacks have reviewed a migration plan.

- The extracted stacks deliberately retain the historical stack IDs, construct IDs, and fixed physical bucket names. Hunter College Clubs also retains its existing IAM names; the new fixed Girls Who Code IAM names must be checked for account-level conflicts before deployment. Preserved template identity helps only when the same existing CloudFormation stacks remain the owners. Different stack ownership can cause fixed-name conflicts or replacement/duplicate CloudFront resources.
- Do not deploy the hosting stacks from both this module and `infrastructure/legacy`. Existing resources may need to remain under their current stacks or be explicitly imported/adopted before this module becomes authoritative.
- Both buckets use `RemovalPolicy_DESTROY` with automatic object deletion. Stack removal or a replacement that deletes the old bucket can permanently delete hosted assets.
- The implementation retains the existing CloudFront OAI and `S3Origin` access design. A future move to origin access control (OAC) should be treated as a separate, reviewed infrastructure migration because it changes origin access and bucket-policy resources.
- The application is environment-agnostic, as in the reference repository. Review the intended AWS account, region, stack names, and synthesized change set before any eventual deployment.

No deployment or resource ownership transfer is performed by this module extraction.

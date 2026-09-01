# AWS frontend hosting

This module contains the independent AWS CDK implementation for the repository's static frontend hosting. It was adapted from [`GWC-Hunter-College/Website-Hosting-Iac`](https://github.com/GWC-Hunter-College/Website-Hosting-Iac) at commit `21152e45112f0fe925b57798675650f1027828a7` and checked against the matching frontend stacks in the preserved Event Management System implementation.

## Current stacks

The Go CDK application synthesizes two environment-agnostic stacks:

- `FrontendStack` hosts Girls Who Code at Hunter. Its private, fixed-name `gwc-club-site` S3 bucket supplies separate CloudFront distributions from the `/main` and `/production` prefixes.
- `FrontendHccStack` hosts Hunter College Clubs / Event Manager. Its private, fixed-name `hunter-college-club-event-site` S3 bucket supplies separate CloudFront distributions from the `/staging` and `/production` prefixes. It also creates the fixed-name `hcc-website-ci-deployer` IAM user and `frontend-hcc-ci-policy`, scoped to uploading/deleting site objects, listing that bucket, and invalidating the two HCC distributions. The stack does not create an IAM access key.

All four distributions redirect viewers to HTTPS, serve `index.html` by default, and return that SPA entry point with status 200 for S3 403 and 404 responses. Each bucket is private and grants CloudFront read access through an origin access identity (OAI).

This implementation defines S3, CloudFront, and the HCC deployment identity only. It does **not** define Route 53 records, ACM certificates, custom domains, frontend build artifacts, or an S3 deployment construct. Frontend delivery automation must build and upload assets to the documented prefixes separately.

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

- The extracted stacks deliberately retain the reference stack IDs, construct IDs, fixed physical bucket names, IAM user name, and IAM policy name. That preserves template identity only when the same existing CloudFormation stacks remain the owners. Deploying these definitions under different stack ownership can conflict with globally unique bucket names and account-unique IAM names, or create replacement/duplicate CloudFront resources.
- Do not deploy the hosting stacks from both this module and `infrastructure/legacy`. Existing resources may need to remain under their current stacks or be explicitly imported/adopted before this module becomes authoritative.
- Both buckets use `RemovalPolicy_DESTROY` with automatic object deletion. Stack removal or a replacement that deletes the old bucket can permanently delete hosted assets.
- The implementation retains the existing CloudFront OAI and `S3Origin` design to avoid changing behavior during this structural refactor. A future move to origin access control (OAC) should be treated as a separate, reviewed infrastructure migration because it changes origin access and bucket-policy resources.
- The application is environment-agnostic, as in the reference repository. Review the intended AWS account, region, stack names, and synthesized change set before any eventual deployment.

No deployment or resource ownership transfer is performed by this module extraction.

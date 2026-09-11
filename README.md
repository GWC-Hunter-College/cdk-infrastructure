# Girls Who Code Hunter Website & Event Management System

Website delivery and event-management services for Girls Who Code at Hunter College and the broader Hunter club community. This repository contains the infrastructure that delivers the separately maintained frontends and the backend that stores and serves club and event data.

**[Visit the Girls Who Code Hunter website →](https://girlswhocodehunter.org)**

[![Girls Who Code Hunter homepage](docs/assets/gwc-homepage.png)](https://girlswhocodehunter.org)

## System overview

The system has two major parts:

- **Frontend delivery:** `hosting/` publishes static website artifacts through S3 and CloudFront, with deployment identities and optional custom domains.
- **Event-management backend:** APIs, application logic, and MySQL persistence support club information, student membership, events, and related records. The integrated implementation is retained in `infrastructure/legacy/` as those responsibilities move into dedicated modules.

Frontend application source is maintained separately. This repository connects website delivery with the services used by frontend and administrative clients.

## Website hosting

The [AWS hosting module](hosting/aws/) defines private S3 origins, staging and production CloudFront distributions, origin access control, and deployment IAM resources. Optional custom-domain configuration adds Route 53 DNS and an ACM certificate. The production website URL above comes from [`hosting/aws/main.go`](hosting/aws/main.go).

**[Read the hosting architecture and deployment guide →](hosting/README.md)**

That guide is the authoritative reference for resource ownership, domain configuration, generated URLs, validation, and production safety. `hosting/aws/` remains its own Go/CDK module.

## Event-management backend

The backend stores and serves students, clubs, memberships and role flags, events, descriptions, tags, and image metadata. The schema also includes club verification and administrator records. See the [implemented schema](infrastructure/legacy/docs/database/schema.md) and [API reference](infrastructure/legacy/docs/api/README.md) for the precise data model and route contracts.

![Integrated backend architecture showing frontend clients, Cognito sign-in, API Gateway HTTP APIs, Go Lambda handlers, private S3 image transfers, RDS MySQL, Secrets Manager, and SSM administrative access](infrastructure/legacy/docs/architecture/assets/backend-architecture.svg)

This source-derived diagram shows the integrated implementation retained in [`infrastructure/legacy/`](infrastructure/legacy/README.md), including its original frontend delivery resources. The detailed, independent website-hosting architecture is documented in `hosting/`.

Clients call development or production API Gateway HTTP APIs backed by Go Lambda handlers. Cognito provides sign-in and JWTs; routes configured with a JWT authorizer validate those tokens. Database handlers connect to RDS MySQL in isolated VPC subnets and retrieve credentials through a Secrets Manager endpoint. Image handlers issue signed requests so clients transfer image bytes directly to private S3, while the database holds metadata. Supporting CDK resources include identity synchronization, database initialization, and an SSM-managed EC2 bastion.

[Explore the backend architecture](infrastructure/legacy/docs/architecture/overview.md) · [View the architecture SVG](infrastructure/legacy/docs/architecture/assets/backend-architecture.svg) · [Edit the Excalidraw source](infrastructure/legacy/docs/architecture/assets/backend-architecture.excalidraw)

## Modular backend architecture

The modular design separates application behavior from the AWS services that run it:

| Boundary | Responsibility |
| --- | --- |
| [`database/`](database/README.md) | MySQL schema, migrations, queries, and persistence concerns for club/event data, independent of RDS provisioning. |
| [`api/`](api/README.md) | API and service behavior exposed to frontend/admin clients, including validation and authorization, independent of API Gateway and Lambda deployment. |
| [`infrastructure/`](infrastructure/README.md) | Backend infrastructure and deployment concerns; [`infrastructure/aws/`](infrastructure/aws/README.md) is the AWS module boundary for compute, networking, secrets, and database hosting. |
| [`infrastructure/legacy/`](infrastructure/legacy/README.md) | The retained integrated implementation: Go handlers, SQL, shared utilities, and AWS CDK stacks together during migration. |

### Migration architecture

Database and API logic belong in their application modules; AWS/CDK code belongs in infrastructure. The migration extracts those responsibilities from the integrated implementation while preserving its behavior and resource ownership. The `api/` and `database/` guides map existing handlers and queries to their destination boundaries; executable backend code remains in `infrastructure/legacy/` until migrated. The modular AWS backend boundary is defined in `infrastructure/aws/`.

## Repository structure

```text
hosting/                  Frontend delivery guide and architecture image
└── aws/                  Independent AWS CDK hosting module
database/                 Persistence boundary, schema, and migration map
api/                      API/service boundary and migration map
infrastructure/           Backend deployment concerns
├── aws/                  Modular AWS backend boundary
└── legacy/               Integrated backend, CDK app, and detailed guides
docs/assets/              Repository homepage preview
```

## Development and deployment

Go and the AWS CDK CLI are the primary tools. Run commands from the relevant Go module, not the repository root. For hosting validation:

```sh
cd hosting/aws
go test ./...
go vet ./...
go build ./...
cdk synth --all
```

For backend setup, Go checks, Lambda bundling, and synthesis prerequisites, follow the [legacy local-development guide](infrastructure/legacy/docs/development/local-development.md). Deployment procedures and ownership checks are documented in the [hosting guide](hosting/README.md) and [backend deployment guide](infrastructure/legacy/docs/architectures/serverless/deployment.md).

Synthesis generates templates; it does not deploy resources. Review the intended account, region, and CloudFormation ownership before deployment. The hosting module and legacy application contain overlapping frontend definitions and must not both manage the same resources. Moving code between modules does not migrate AWS resources automatically.

## License

See [LICENSE](LICENSE).

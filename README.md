# Hunter College Event Management System backend

This repository contains the backend and deployment infrastructure for Hunter College's club and event-management platform. It supports frontend applications such as:

- Girls Who Code at Hunter;
- Hunter College Clubs / Event Manager; and
- related administrative and event-management interfaces.

The frontend application source is maintained separately. This repository contains the backend boundaries and the infrastructure definitions used to run or host them.

## Repository areas

This first modularization phase organizes the repository in the following order:

1. **[Hosting](hosting/README.md)** — independent static-frontend hosting. [`hosting/aws/`](hosting/aws/README.md) contains the AWS CDK implementation adapted from the standalone Website-Hosting-Iac repository.
2. **[Database](database/README.md)** — the future provider-independent MySQL module. This phase creates its boundary and documentation only.
3. **[API](api/README.md)** — the future provider-independent backend application used by frontend clients to access data. This phase creates its boundary and documentation only.
4. **[Infrastructure](infrastructure/README.md)** — provider-specific deployment concerns. [`infrastructure/aws/`](infrastructure/aws/README.md) is the future AWS boundary for the API and database, while [`infrastructure/legacy/`](infrastructure/legacy/README.md) preserves the current integrated implementation.

| Area | State after this phase |
| --- | --- |
| `hosting/aws/` | Real, self-contained AWS CDK hosting implementation sourced from the accessible reference repository |
| `database/` | README and schema visual only; database code has not been migrated |
| `api/` | README only; API code has not been migrated |
| `infrastructure/aws/` | README only; new API/database infrastructure has not been implemented |
| `infrastructure/legacy/` | Preserved current all-in-one Event Management System implementation and detailed documentation |

## Current implementation

The following diagram describes the **currently implemented integrated system**, not the future modular architecture:

![Current high-level architecture showing frontend clients, CloudFront and S3 hosting, Cognito, API Gateway, Lambda, image storage, and RDS](infrastructure/legacy/docs/assets/high-level-architecture.png)

Today, the working implementation combines AWS CDK stacks, API Gateway, Lambda handlers, Cognito, RDS MySQL, S3, CloudFront, networking, database initialization, and supporting resources. It remains intact under [`infrastructure/legacy/`](infrastructure/legacy/README.md) while later phases separate what the backend does from how it is hosted.

The intended direction is:

- static frontend deployment remains independent in `hosting/`;
- MySQL schema, migrations, and database logic move toward `database/` without depending on RDS;
- backend application behavior moves toward `api/` without being defined by Lambda or API Gateway; and
- AWS-specific API/database deployment moves toward `infrastructure/aws/`.

No database or API implementation was migrated in this phase. No production deployment or CloudFormation resource migration is implied by this layout.

## Documentation and safety

- Start with the [legacy implementation guide](infrastructure/legacy/README.md) for the current build, architecture, API, database, and deployment documentation.
- Read the [hosting AWS guide](hosting/aws/README.md) before synthesizing the extracted hosting stacks.
- Do not deploy either CDK application solely because its files moved. Existing S3 buckets, CloudFront distributions, IAM resources, DNS, certificates, and CloudFormation stack ownership must be reviewed before any production migration.

## License

See [LICENSE](LICENSE).

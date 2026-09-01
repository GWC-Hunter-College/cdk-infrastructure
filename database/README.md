# Database

This directory is the future boundary for the provider-independent relational database implementation. It is a documented placeholder in this first modularization phase; no database implementation has been moved here yet. The current integrated implementation remains under [`infrastructure/legacy/`](../infrastructure/legacy/).

The intended database technology is MySQL. Amazon RDS is one possible environment for hosting MySQL, not the database implementation itself. The same implementation should eventually be usable with:

- Amazon RDS for MySQL
- MySQL on Amazon EC2
- another compatible MySQL host
- local MySQL

This module may eventually contain the schema, initialization, migrations, database-specific logic, and seed or test data. No schema redesign or database migration is part of this phase.

## Current schema reference

The image below documents the schema used by the current legacy implementation. It is retained as a reference and does not imply that the schema or its implementation has been migrated into this module. See the [detailed legacy database documentation](../infrastructure/legacy/docs/database/README.md) for the current implementation.

![Current legacy database schema](assets/database-schema.png)

# API

This directory is the future boundary for the provider-independent backend API application. It is a documented placeholder in this first modularization phase; no API implementation has been moved here yet. The current integrated implementation remains under [`infrastructure/legacy/`](../infrastructure/legacy/).

The API is the layer through which the separately maintained frontend applications communicate with backend data:

```text
Frontend -> API -> Database
```

This module may eventually contain event, club, membership, image and upload endpoints; authentication integration; authorization logic; and database access. AWS Lambda and Amazon API Gateway should be deployment and hosting choices for the application, not the definition of the API itself.

No endpoints, authentication or authorization behavior, or response formats are being migrated or changed in this phase.

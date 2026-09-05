# Hosting

This area contains infrastructure for hosting the separately maintained frontend applications supported by this backend repository, including:

- Girls Who Code at Hunter
- Hunter College Clubs / Event Manager

AWS currently hosts these static sites. The AWS implementation belongs in [`aws/`](aws/): Amazon S3 stores the deployed static assets and Amazon CloudFront serves them. Route 53 manages the Girls Who Code production domain, with ACM providing its HTTPS certificate. AWS CDK and CloudFormation define the implemented hosting resources.

Frontend hosting is kept separate from the infrastructure used to run the backend API and database.

# Hosting

This area contains infrastructure for hosting the separately maintained frontend applications supported by this backend repository, including:

- Girls Who Code at Hunter
- Hunter College Clubs / Event Manager

AWS currently hosts these static sites. The AWS implementation belongs in [`aws/`](aws/): Amazon S3 stores the deployed static assets, Amazon CloudFront serves them, and Amazon Route 53 manages DNS where applicable. AWS CDK and CloudFormation define the hosting resources.

Frontend hosting is kept separate from the infrastructure used to run the backend API and database.

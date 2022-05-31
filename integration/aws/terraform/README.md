# AWS Integration Test Terraform

There are two environments:
- `dev`
- `prod`

The `dev` environment exists to enable testing of changes to the testing environment without breaking test runs.

Within each environment, there are two Terraform modules:
- `bootstrap`
- `environment`

`bootstrap` contains all the Terraform necessary to supply the GitHub Actions workflow with AWS credentials.
It also contains the IAM policies that restrict what the GitHub actions can do within AWS (to prevent malicious and accidental excess resource usage).
The two bootstrap terraform states must be applied manually, with admin credentials in the respective accounts.

The `environment` state contains the resources being used for testing.
When changes are made to the terraform files, the "Integration Test Setup" GitHub Action Workflow runs to apply those changes.
PRs cause the changes to be applied on the `dev-environment` state, while commits on the `main` branch cause changes to be applied on the `prod-environment` state.
The workflow can also be triggered manually on either environment based on the configuration in the `main` branch.

Resources in the `environment` module should be tagged with `IntegrationTest=true`, as well as a `test` tag that identifies the resource along with its type (the exact value doesn't really matter).
Both of those tags are used by the integration test code (resources without `IntegrationTest=true` are filtered out to reduce issues that may arise with resources not created for testing).

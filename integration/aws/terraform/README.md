# AWS Integration Test Terraform

There are two environments:
- `dev`
- `prod`

The `dev` environment exists to enable testing of changes to the testing environment without breaking test runs.

Within each environment, there are two Terraform modules:
- `bootstrap`
- `environment`

## `bootstrap` Module
`bootstrap` contains all the Terraform necessary to supply the GitHub Actions workflow with AWS credentials.
It also contains the IAM policies that restrict what the GitHub actions can do within AWS (to prevent malicious and accidental excess resource usage).
The two bootstrap terraform states must be applied manually, with admin credentials in the respective accounts.

### Initial setup
Before applying the bootstrap module to a new account, you must create a bucket named `-terraform-state` prefixed with the AWS account ID, and a DynamoDB table named `terraform-locks` with a `string` partition key named `LockID` (on-demand capacity recommended).

### Applying manually
After inital setup is performed, the bootstrap resources can be provisioned with the standard `terraform init`/`terraform plan`/`terraform apply`.
Note that there may be a conflict with existing resources since created resources have a fixed ID, especially the IAM OpenID Connect Provider.

## `environment` Module
The `environment` state contains the resources being used for testing.
When changes are made to the terraform files, the "Integration Test Setup" GitHub Action Workflow runs to apply those changes.
PRs cause the changes to be applied on the `dev-environment` state, while commits on the `main` branch cause changes to be applied on the `prod-environment` state.
The workflow can also be triggered manually on either environment based on the configuration in the `main` branch.

Resources in the `environment` module should be tagged with `IntegrationTest=true` when possible, as the utility functions for testing filter out any resource without that tag to reduce issues that may arise with resources not created for testing.
Resources should also be tagged with a `test` tag with a value specific to the integration test it is for (e.g. `ec2-instance-0`; the numeric suffix makes it easy to support many resources of the same type if the test requires).
This is to prevent a single resource from being matched by tests that aren't expecting it to be present;
for example, the Terraform `aws_instance` and `aws_autoscaling_group` resources will both spawn EC2 instances,
but we don't want the tests for EC2 instances to be affected by the tests for autoscaling groups.

### Applying manually
The `environment` module has almost no dependency on the `bootstrap` module, with the exception of the `test-lambda-execution-role` IAM role.
If you are applying the `environment` module in an account that hasn't had the `bootstrap` module applied, you must first create a `test-lambda-execution-role` IAM role that can be assumed by AWS Lambda (no policies need to be attached to the role).

Like the `bootstrap` module, applying the `environment` module manually is just a standard `terraform init`/`terraform plan`/`terraform apply`.

### Applying Automatically
The `environment` module will be applied automatically on push events to the prod account if the branch is `main`, or otherwise the `dev` account.

# AWS Integration Tests

The AWS integration tests allow us to quickly and easily test the correctness of our AWS provider functions.
Each test runs the fetch functions for each resource type, making API calls to a real AWS account.
Because these tests depend on external state that cannot quickly be setup and torn down for every test,
the integration tests do not always run.
Integration tests only run when either there are AWS credentials present for one of the two known AWS accounts (`316817240772` and `438881294876`; dev and prod, respectively) or if forced to run due to an environment var being set (either `CI` or `CLOUD_INTEGRATION_TESTS`).

## Running Manually
In order to run integration tests manually, you must either set one of the above env vars, or run the tests while the terminal has active credentials for one of the above accounts (or both).
As tests will make use of the credentials provided, the tests are sensitive to the contents of the account on the other end.
The resources expected by tests can be found within the `./terraform` directory.

## Skipping
If you need to avoid the integration tests, but you are unable to dissatisfy the above conditions, you can run `go test` with `-short`,
which will cause all AWS integration tests to be skipped.

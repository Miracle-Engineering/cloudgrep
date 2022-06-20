# Demo

The demo shows some sample data for a medium-sized AWS account (a few hundred cloud resources).

## How can cloudgrep help my organization?

This demo demonstrates how cloudgrep can help with:

- Viewing all your cloud resources for multiple regions in one browser.
- Searching your cloud resources using their tags to measure the progress of your IaC initiative.
- Verifying that your tag values are correct, quickly identifying the misconfigured values.
- Enforcing your tag policies by identifying the resources missing some tags.

## About the demo AWS account

Some interesting facts about the demo account:
- The infrastructure is deployed in 2 AWS regions: `us-east-1` and `eu-west-2`.
- There are some production and development environments. A tag called `env` is used.
- There are many teams that own their infrastructure. A tag called `team` is used.
- Some infrastructure was provisionned with CloudFormation, some with Terraform and some was manual or already existing (ex: default VPC). A tag called `managed-by` is used.

## Run the demo

The cloudgrep demo database has been captured so you can run the cloudgrep demo account without an AWS account.

```
./cloudgrep --config demo/demo.yaml
```

## Interesting use cases to demo

1. Show how to filter on region and type.
    - Show the EKS clusters in each region.
1. Explore the labels: show the most popular labels and explain how they are used:
    - `managed-by`: track the progress of IaC initiative. Internal goal is to migrate from CloudFormation to Terraform.
    - `team`: the goal is to assign and track every infrastructure for each team.
    - `env`: this tag is used to differentiate the developement and production environment.
1. Show all the infrastructure for a specific team.
    - Pick *billing*
    - Then show only the RDS DB instances for this team.
1. Use the tag `managed-by`
    - View all the infrastucture that is not managed by Terraform, or Cloudformation.
    - Then identify the EC2 Instances.
1. Use the `env` tag
    - find the resource that uses "production" instead of "prod".
1. Use the `team` tag
    - Find if there are any RDS DB Instances that is missing this tag.

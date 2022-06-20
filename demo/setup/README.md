# Demo setup

The demo data consists of a few cloudformation stacks and some EKS clusters created with `eksctl`.

## AWS account

Use the *cloudgrep-demo* account.

## EKS clusters

```
eksctl create cluster --config-file=cluster-dev-us-east-1.yaml

eksctl create cluster --config-file=cluster-prod-us-east-1.yaml

eksctl create cluster --config-file=cluster-prod-eu-west-3.yaml
```

## CloudFormation stacks

The stacks can be seen in the AWS console for the supported regions:
- [us-east-1](https://us-east-1.console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks)
- [eu-west-3](https://eu-west-3.console.aws.amazon.com/cloudformation/home?region=eu-west-3#/stacks)

Note: `eksctl` also create some CloudFormation stacks. So everything is in fact in CloudFormation.

## Turning off the demo account

The easiest way is to delete all the CloudFormation stacks.

## Generate the demo DB

1. Recreate the original database (optional)

The only reason to generate the original is to have new resource types that didn't exist before or if the database model has changed.

```shell
# delete the previous dump file
rm demo/setup/original-dump.db

# run cloudgrep once to fetch the database
./cloudgrep --config demo/setup/demo-setup.yaml

```

2. Run the update script

This script will add some resources and update some tags to make it more relevant for a demo purpose.
The input of this script is the original DB, the output is the demo DB.

```shell
go run demo/setup/update_demo.go
```

3. Test the demo data

```shell
# run cloudgrep using the newly created demo db
./cloudgrep --config demo/demo.yaml

```

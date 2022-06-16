# Development

This application provides a single binary that contains the backend and the frontend.

## Project structure

The project structure is inspired by https://github.com/sosedoff/pgweb

    .
    ├── cmd                                  # CLI semi-generated framework (also Golang)
    ├── hack                                 # Tools for development
    ├── pkg                                  # Backend: Golang
    ├── fe                                   # Frontend: React frontend code
    ├── static                               # Frontend: Static assets
    ├── main.go                              # Executable main function
    └── Makefile                             # Defines the tasks to be executed - for local or CI run

## Backend

### Starting the server

```shell
# using local code
make run
```

```shell
# using local binary
make build
./cloudgrep
```

### Running Tests

```shell
make test
```

## Frontend
Checkout the frontend development guide [here](https://github.com/run-x/cloudgrep/blob/main/fe/README.md)

### AWS Resource supported

| Type            |  Status |
|-----------------| ----------- |
| EC2 Instance    |  :white_check_mark: |
| Load Balancer   |  :white_check_mark: |
| S3 Bucket       |  :white_check_mark: |
| EBS             |  :white_check_mark: |
| RDS             |  :white_check_mark: |
| Lambda Function |  :white_check_mark: |

## Design

![design diagram](img/cloudgrep-design.png)

All of these boxes are implemented as distinct Go packages, except for UI which is a JS app.

## API
API design is documented [here](https://github.com/run-x/cloudgrep/blob/main/API.md)

## Configure a new AWS resource
1. If the resource you are adding is for a new, wholly unsupported service, add a new item in the `services` list in the `pkg/provider/aws/config.yaml` file, and then create a new file with `.yaml` appended to the service name in the same directory.
    Use the canonical service initialism or name.
    If the [AWS SDK for Go v2](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2) uses a different name for the service package than what is used within Cloudgrep, add that package name to the `servicePackage` field in the service configuration file (see `elb.yaml` for an example).
2. Add a new item to the `types` list.
    The schema for the definition, along with the documentation for each field, can be found in the `hack/awsgen/config/types.go` file (the file as a whole is the the `Service` struct).
    You can use the existing type definitions in the other adjacent `.yaml` files as a guide.
    Many APIs return tag data directly in the list/describe APIs (configured in the `listApi` field), but if it doesn't,
    you must configure the `getTagsApi` field in the type.
3. \[Optional\] If you need to customize the API call's input, you can use `inputOverrides` to hook the creation of the input struct.
    Using `inputOverrides.fieldFuncs` you can set specific fields, but if you need more control, you can use `inputOverrides.fullFuncs`.
4. Run `make awsgen` to generate the AWS provider resource functions.

## Manually configuring a new AWS resource
If code generation is not sufficient for a specific resource, you can also manually implement the function(s) for a resource using the following instructions.
This method is not preferred, and instead it is preferred to improve `awsgen` as needed to support the resource:
implementing resources manually requires more effort to maintain and improve if we change how the fetch functions work across the board.

1. Implement a new AWS provider method (in pkg/provider/aws) to fetch your resources. It must have the type
   signature of a FetchFunction like so: `type FetchFunc func(context.Context, chan<- model.Resource) error`.
   Please refer to the following example:
   ```go
   func (p *Provider) FetchEC2Instances(ctx context.Context, output chan<- model.Resource) error {
       resourceType := "ec2.Instance"
       ec2Client := ec2.NewFromConfig(p.config)
       input := &ec2.DescribeInstancesInput{}
       paginator := ec2.NewDescribeInstancesPaginator(ec2Client, input)

       resourceConverter := p.converterFor(resourceType)
       for paginator.HasMorePages() {
           page, err := paginator.NextPage(ctx)
           if err != nil {
               return fmt.Errorf("failed to fetch EC2 Instances: %w", err)
           }

           for _, r := range page.Reservations {
               if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, r.Instances); err != nil {
                   return err
               }
           }
       }

       return nil
   }
   ```
2. [Optional] Implement the method to return the tags. Unless there is already a `Tags` field, this method would need
    to be implemented. It should have the type signature `type tagFunc[T any] func(context.Context, T) (model.Tags, error)`
    where T is the aws resource being ingested. Here is an example for Load Balancer:
   ```go
   func (p *Provider) FetchLoadBalancerTag(ctx context.Context, lb types.LoadBalancer) (model.Tags, error) {
       elbClient := elbv2.NewFromConfig(p.config)
       tagsResponse, err := elbClient.DescribeTags(
           ctx,
           &elbv2.DescribeTagsInput{ResourceArns: []string{*lb.LoadBalancerArn}},
       )
       if err != nil {
           return nil, fmt.Errorf("failed to fetch tags for load balancer %v: %w", *lb.LoadBalancerArn, err)
       }
       var tags model.Tags
       for _, tagDescription := range tagsResponse.TagDescriptions {
           for _, tag := range tagDescription.Tags {
               tags = append(tags, model.Tag{Key: *tag.Key, Value: *tag.Value})
           }
       }
       return tags, nil
   }
    ```
3. Add a new mapper struct to the list returned by getTypeMapping in `pkg/provider/aws/types.go`
   This struct has 4 fields:
   * `IdField`: Name of id field in ingested aws resource struct
   * `TagField`: A struct of the TagField type used to dictate where. Leave this empty if passing the tagfields separately like with the load balancer
   * `FetchFunc`: The fetch function created in step 1
   * `IsGlobal`: Set to true if this is a global resource (like a Hosted Zone). Otherwise leave empty.
   Example:
   ```go
   "ec2.Instance": {
       IdField:   "InstanceId",
       Field:  defaultTagField,
       chFunc: p.FetchEC2Instances,
   }
   ```


These methods will be automatically called at startup.
The mapper definition will be used to convert the returned type to some `model.Resource` objects.

## Release

The release process is automated when merging a Pull Request.

### How to trigger a release

1. Create a Pull Request.
1. Attach a label [`bump:patch`, `bump:minor`, or `bump:major`]. Cloudgrep uses [haya14busa/action-bumpr](https://github.com/haya14busa/action-bumpr).
1. [The release workflow](.github/workflows/release.yml) automatically tags a
   new version depending on the label and create a new release on merging the
   Pull Request.

If you do not want to create a release for a given PR, do not attach a bump label to it.

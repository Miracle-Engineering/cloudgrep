# development notes

This application provides a single binary that contains the backend and the frontend.

## Project structure

The project structure is inspired by https://github.com/sosedoff/pgweb

    .
    ├── pkg                                  # Backend: Golang
    └── static                               # Frontend: Static assets
    └── main.go                              # Executable main function
    └── Makefile                             # Defines the tasks to be executed - for local or CI run

## Start the server

```shell
# using local code
make run
```

```shell
# using local binary
make build
./cloudgrep
```

### AWS Resource supported

| Type |  Status |
| ------------- | ------------- |
| EC2 Instance |  :white_check_mark: |
| Load Balancer |  :white_check_mark: |
| S3 Bucket |  :x: |
| EBS |  :x: |
| RDS |  :x: |

## Design

![design diagram](img/cloudgrep-design.png)

- **cli**: this the application entry point, starts all the various components
- **provider**: the provider is responsible for fetching the resource data from the cloud and write it to the datastore. The implementation specific to a cloud provider is done in their own package.
- **datastore**: the datastore provides an interface to read/write/update the collected data. The storage is done in a database.
- **api**: the api is the Gin HTTP server. It defines the routes and implements the api endpoints.
- **config**: all the application configuration is defined in this package. The user can provide a `config.yaml` file or use the default values.
- **model**: these are the base objects to contain *resources* and *tags*.
- **UI**: this is the frontend layer.

All of these boxes are implemented as distinct Go packages, except for UI which is a JS app.

## Configure a new resource

1. Define a mapping in the corresponding provider `mapping.yaml`
    ```yaml
    mappings:
      # the Go type that is returned by the cloud provider
      - type: "github.com/aws/aws-sdk-go-v2/service/ec2/types.Instance"
        # the resource type as displayed by cloudgrep
        resourceType: "ec2.Instance"
        # the name of the id field
        idField: InstanceId
        # the method to call to fetch the resources, it must be implemented
        impl: FetchEC2Instances
        # the method to call to generate the tags only needed if there is not already a "Tags" field
        #tagImpl: FetchEC2Tags
    ```
1. Implement the method define in the mapping
    ```go
    // this method is named after the mapping.impl value and return a slice of the mapping.type value
    func (awsPrv *AWSProvider) FetchEC2Instances(ctx context.Context) ([]types.Instance, error) {
        input := &ec2.DescribeInstancesInput{}
        var instances []types.Instance
        result, err := awsPrv.ec2Client.DescribeInstances(ctx, input)
        if err != nil {
            return nil, err
        }

        for _, r := range result.Reservations {
            instances = append(instances, r.Instances...)
        }
        return instances, nil
    }
    ```
1. Implement the method to return the tags. Unless there is already a `Tags` field, this method would need to be implemented. Here is an example for Load Balancer.
    ```go
    // this method is named after the mapping.tagImpl value and return some model.Tags
    // The ELB doesn't have a Tags field so this method calls `elasticloadbalancingv2.DescribeTags`
    func (p *AWSProvider) FetchLoadBalancerTag(ctx context.Context, lb types.LoadBalancer) (model.Tags, error) {
        tagsResponse, err := p.elbClient.DescribeTags(
            ctx,
            &elbv2.DescribeTagsInput{ResourceArns: []string{*lb.LoadBalancerArn}},
        )
        if err != nil {
            return nil, fmt.Errorf("failed to fetch tags for load balancer %v: %w", &lb.LoadBalancerArn, err)
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


This method will be automatically called at startup.  
The mapping definition will be used to convert the returned type to some `model.Resource` objects.

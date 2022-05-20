

# API

The backend exposes an API at `http://localhost:8080/api`.

## List resources

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resources](http://localhost:8080/api/resources)  | GET  | Return list of cloud resources |  :white_check_mark: |

| Parameters | Description |  Examples |
| ------------- | ------------- | ------------- |
| tags  | return resource(s) with the provided tag  | `tags[team]=infra` return resources with the tag `team=infra`, meaning "team" with value "infra" <br />`tags[team]=infra&tags[env]=prod` return resources with the tags `team=infra` **and** `env=dev` <br />`tags[env]=prod,dev` return resources with the tags `env=prod` **and** `env=dev` <br />`tags[team]=*` return all the resources with the tag `team` defined|
| exclude-tags  | return resource(s) without the provided tag  | `exclude-tags=team` return resources without the tag `team`<br />`exclude-tags=team,env` return resources without the tag `team` **and** `env`

## Get a resource

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resource](http://localhost:8080/api/resource)  | GET  | Return a resource |  :white_check_mark: |

| Parameters | Description |  Examples |
| ------------- | ------------- | ------------- |
| id  | the resource id  | `id=i-024c4971f7f510c8f` return resource with the id `i-024c4971f7f510c8f`

## List fields

Return the list of fileds available for filtering the resources.
The fields can be used when querying the resources. (details TBD)

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/fields](http://localhost:8080/api/fields)  | GET  |  |  :white_check_mark: |

Example of response:
```js
[
  {
    "name":"region",
    "count":16,
    "values":[
      {
        "value":"us-east-1",
        "count":8
      },
      {
        "value":"us-west-2",
        "count":8
      }
    ]
  },
  {
    "name":"type",
    "count":16,
    "values":[
      {
        "value":"ec2.Instance",
        "count":3
      },
      {
        "value":"ec2.Volume",
        "count":6
      },
      {
        "value":"elb.LoadBalancer",
        "count":1
      },
      {
        "value":"s3.Bucket",
        "count":6
      }
    ]
  },
  {
    "name":"cluster",
    "count":6,
    "values":[
      {
        "value":"prod",
        "count":2
      },
      {
        "value":"dev",
        "count":2
      },
      {
        "value":"stage",
        "count":2
      }
    ]
  }
]
```


## Mocked data

There is also a mocked API at `http://localhost:8080/mock/api`.  
The mocked api serves static data, it doesn't handle any query parameters.

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resources](http://localhost:8080/mock/api/resources)  | GET  |  Example of "resources" reponse | :white_check_mark: |
| [/tags](http://localhost:8080/mock/api/tags)  | GET  |  Example of "tags" ressponse |  :white_check_mark: |


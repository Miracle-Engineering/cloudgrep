

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

Example of queries: (not yet supported)
```js

//default return all the resources (no payload)
{}

//return resources of type "ec2.Instance" with the tag "team" equals "marketplace"
{
  "filter":{
    "type": "ec2.Instance",
	  "team": "marketplace"
  }
}

//return resources with the tag "team"
{
  "filter":{
	  "team": { "$neq": "" }
  }
}

//return resources missing the tag "team"
{
  "filter":{
	  "team": { "$eq": "" }
  }
}

//filter with more than one value for a field
// will return resources with type=ec2.Volume AND team IN ("marketplace", "shipping")
{
  "filter":{
    "type":"ec2.Volume",
    "$or": [
      { "team": "marketplace" },
      { "team": "shipping" }
    ]
  }
}

//sort by a field
{
  "filter":{
    "type": "s3.Bucket"
  },
  "sort": ["region"]
}

//The default order for column is ascending order but you can control it with an optional prefix: + or -. + means ascending order, and - means descending order.
//sort by region desc
{
  "filter":{
    "type": "s3.Bucket"
  },
  "sort": ["-region"]
}

//Set a limit: default 25, Max is 100
//return the ec2.Instance with a limit of 10 results
{
  "limit": 10,
  "filter":{
    "type": "ec2.Instance"
  }
}

//used with limit, the offset paramerter specifies the number of rows to skip before any rows are retrieved
//first page: first 10 results
{
  "limit": 10,
  "offset": 0,
  "filter":{
    "type": "ec2.Instance"
  }
}
//second page: next 10 results
{
  "limit": 10,
  "offset": 10,
  "filter":{
    "type": "ec2.Instance"
  }
}

```


## Get a resource

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resource](http://localhost:8080/api/resource)  | GET  | Return a resource |  :white_check_mark: |

| Parameters | Description |  Examples |
| ------------- | ------------- | ------------- |
| id  | the resource id  | `id=i-024c4971f7f510c8f` return resource with the id `i-024c4971f7f510c8f`

## List fields

Return the list of fields available for filtering the resources.

The fields can be presented to the user to enable the user to build a search query using these field.

A field can be:
- a resource type
- a region
- a tag key

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/fields](http://localhost:8080/api/fields)  | GET  | Return the fields available for the stored resources |  :white_check_mark: |

Example of response:
```js
[
  {
    "name":"region",
    "group": "core",
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
    "group": "core",
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
    "group": "tags",
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


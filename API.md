

# API

The backend exposes an API at `http://localhost:8080/api`.

<details>
<summary>List resources</summary>

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resources](http://localhost:8080/api/resources)  | POST  | Return list of cloud resources |  :white_check_mark: |

To filter the resources, send a body containing a query.
The name of the fields used in `filter` and `sort` is the `field.name` returned in the `/fields/` API.

```js
{
  // set a limit, default is 25, max is 100 
  "limit": 25,
  //specifies the number of rows to skip before any rows are retrieved
  "offset": 0,
  //filter the resources
  "filter": {
    "type": "ec2.Instance"
  }
  //optional sort
  "sort": ["type"]
}
```

```shell
curl --location --request POST 'http://localhost:8080/api/resources' \
--header 'Content-Type: application/json' \
--data-raw '{
    "filter": {
        "kubernetes.io/created-for/pv/name": "opta-persistent-0-hellopv-hellopv-k8s-service-0"
    }
}'
```

Example of queries:
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

//return resources with the tag "team" defined
{
  "filter":{
	  "team": "(not null)"
  }
}

//return resources missing the tag "team"
{
  "filter":{
	  "team": "(null)"
  }
}

//filter with more than one value for a field using a OR
// will return resources with type=ec2.Volume AND (team="marketplace" OR team="shipping")
{
  "filter":{
    "type":"ec2.Volume",
    "$or": [
      { "team": "marketplace" },
      { "team": "shipping" }
    ]
  }
}

//Using multiple OR sections
// will return resources with (team="marketplace" OR team="shipping") AND (cluster="dev" OR cluster="prod")  AND (size="large" OR size="medium") 
{
  "filter":{
    "$or": [
      { "team": "marketplace" },
      { "team": "shipping" }
    ],
    "$and": [
      { "$or": [
        { "cluster": "dev" },
        { "cluster": "prod" }
      ] },
      { "$or": [
        { "size": "large" },
        { "size": "medium" }
      ] }
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
</details>
<details>
<summary>Get a resource</summary>

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
</details>
<details>
<summary>Get Engine Status</summary>

Returns the Status of the Cloudgrep run.

| Route                                                   | Method | Description              |  Status |
|---------------------------------------------------------| ------------- |--------------------------| ------------- |
| [/enginestatus](http://localhost:8080/api/enginestatus) | GET  | Return the Engine status |  :white_check_mark: |

Sample Responses:
```js
// Engine completed successfully
{
    "runId": "6fd67489-d852-4962-95bc-eea01159993f",
    "eventType": "engine",
    "status": "success",
    "providerName": "",
    "resourceType": "",
    "error": "",
    "createdAt": "2022-06-22T02:54:12.727066+05:30",
    "updatedAt": "2022-06-22T02:54:25.458235+05:30",
    "childEvents": [
        {
            "runId": "6fd67489-d852-4962-95bc-eea01159993f",
            "eventType": "provider",
            "status": "success",
            "providerName": "aws",
            "resourceType": "",
            "error": "",
            "createdAt": "2022-06-22T02:54:12.727395+05:30",
            "updatedAt": "2022-06-22T02:54:13.979699+05:30",
            "childEvents": null
        },
        {
            "runId": "6fd67489-d852-4962-95bc-eea01159993f",
            "eventType": "resource",
            "status": "success",
            "providerName": "AWS Provider for account 693658092572, region us-east-2",
            "resourceType": "ec2.Volume",
            "error": "",
            "createdAt": "2022-06-22T02:54:13.980207+05:30",
            "updatedAt": "2022-06-22T02:54:16.658743+05:30",
            "childEvents": null
        }
    ]
}

// Engine is currently running
{
    "runId": "6fd67489-d852-4962-95bc-eea01159993f",
    "eventType": "engine",
    "status": "failed",
    "providerName": "",
    "resourceType": "",
    "error": "1 error message\n error message",
    "createdAt": "2022-06-22T02:54:12.727066+05:30",
    "updatedAt": "2022-06-22T02:54:25.458235+05:30",
    "childEvents": [
    {
        "runId": "6fd67489-d852-4962-95bc-eea01159993f",
        "eventType": "provider",
        "status": "success",
        "providerName": "aws",
        "resourceType": "",
        "error": "",
        "createdAt": "2022-06-22T02:54:12.727395+05:30",
        "updatedAt": "2022-06-22T02:54:13.979699+05:30",
        "childEvents": null
    },
    {
        "runId": "6fd67489-d852-4962-95bc-eea01159993f",
        "eventType": "resource",
        "status": "failed",
        "providerName": "AWS Provider for account 693658092572, region us-east-2",
        "resourceType": "ec2.Volume",
        "error": "error message",
        "createdAt": "2022-06-22T02:54:13.980207+05:30",
        "updatedAt": "2022-06-22T02:54:16.658743+05:30",
        "childEvents": null
    }
]
}

```

If you need to know when the engine is done running, keep pulling this endpoint until the status is no longer **fetching**.

</details>
<details>
<summary>Refresh the resources</summary>

Trigger the engine to refresh the cloud resources.
Calling this endpoint will returns immediately, the engine will start fetching the resources async.

| Route                                                   | Method | Description              |  Status |
|---------------------------------------------------------| ------------- |--------------------------| ------------- |
| [/refresh](http://localhost:8080/api/refresh) | POST  | Refresh the cloud resources |  :white_check_mark: |

Sample Responses:
```js
// Refresh request acknowledged, the refresh has started.
code: 200
body: {}

// The refresh has already been triggered and is in progress
code: 202
{
  "status":"202",
  "error":"engine is already running"
}

// There was an error
code: 400
{
  "status":"400",
  "error":"can't connect to datastore"
}

```

Once the refreshed is triggered, call **Get Engine Status** API to know if the refresh is done.
</details>

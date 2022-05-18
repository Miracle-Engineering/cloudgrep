

# API

The backend exposes an API at `http://localhost:8080/api`.

## List resources

List the resources for a given filter.

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resources](http://localhost:8080/api/resources)  | GET  | Return list of cloud resources |  :white_check_mark: |

| Parameters | Description |  Examples |
| ------------- | ------------- | ------------- |
| tags  | return resource(s) with the provided tag  | `tags[team]=infra` return resources with the tag `team=infra`, meaning "team" with value "infra" <br />`tags[team]=infra&tags[env]=prod` return resources with the tags `team=infra` **and** `env=dev` <br />`tags[env]=prod,dev` return resources with the tags `env=prod` **and** `env=dev` <br />`tags[team]=*` return all the resources with the tag `team` defined|
| exclude-tags  | return resource(s) without the provided tag  | `exclude-tags=team` return resources without the tag `team`<br />`exclude-tags=team,env` return resources without the tag `team` **and** `env`|

Example of response:
```js
[
  {
    "id":"arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/net/staging-ingress/14522ba1bd959dd6",
    "region":"us-east-1",
    "type":"elb.LoadBalancer",
    "tags":[
      {
        "key":"cluster",
        "value":"staging"
      },
      {
        "key":"service.k8s.aws/resource",
        "value":"LoadBalancer"
      },
      {
        "key":"service.k8s.aws/stack",
        "value":"ingress-nginx/ingress-nginx-controller"
      }
    ],
    "properties":[
      {
        "name":"CreatedTime",
        "value":"2022-04-07 22:27:39.35 +0000 UTC"
      },
      {
        "name":"DNSName",
        "value":"staging-ingress-14522ba1bd959dd6.elb.us-east-1.amazonaws.com"
      },
      {
        "name":"IpAddressType",
        "value":"ipv4"
      },
      {
        "name":"LoadBalancerArn",
        "value":"arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/net/staging-ingress/14522ba1bd959dd6"
      },
      {
        "name":"LoadBalancerName",
        "value":"staging-ingress"
      }
    ]
  }
]
```

## List tags

List the tags associated with the resources for a given filter.  
This list is sorted to return the most popular tags first (highest count).

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/tags](http://localhost:8080/api/tags)  | GET  | Return list of tags |  :white_check_mark: |

| Parameters | Description |  Examples |
| ------------- | ------------- | ------------- |
| tags  | return tags for the resources with the provided tag  | `tags[team]=infra` return the tags associated with the resources with the tag `team=infra`" <br />`tags[team]=infra&tags[env]=prod` return the tags associated with the resources with the tags `team=infra` **and** `env=dev` <br />`tags[env]=prod,dev` return the tags associated with the resources with the tags `env=prod` **and** `env=dev` <br />`tags[team]=*` return the tags associated with the resources with the tag `team` defined|
| exclude-tags  | return the tags associated with the resources without the provided tag  | `exclude-tags=team` return resources without the tag `team`<br />`exclude-tags=team,env` return the tags associated with the resources without the tag `team` **and** `env`
| limit  | the number of tags to return.<br />Default: 10, Maximum value: 100.

Example of response:
```js
[
  {
    "key":"cluster",
    //the tag values found for this key
    "values":[
      "dev-cluster",
      "prod-cluster",
    ],
    //there are 6 resource with this tag key
    "count":6,
    "ResourceIds":[
      "i-024c4971f7f510c8f",
      "i-046b8584f97edce25",
      "i-0ce5fd258122ee34b",
      "vol-06690829257c1451a",
      "vol-0a6cf8e1480199cb3",
      "vol-0effb20041bde4898"
    ]
  }
]
```

## Get a resource

Return a resource by id.

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resource](http://localhost:8080/api/resource)  | GET  | Return a resource |  :white_check_mark: |

| Parameters | Description |  Examples |
| ------------- | ------------- | ------------- |
| id  | the resource id  | `id=i-024c4971f7f510c8f` return resource with the id `i-024c4971f7f510c8f`

Example of response:
```js
{
  "id":"arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/net/staging-ingress/14522ba1bd959dd6",
  "region":"us-east-1",
  "type":"elb.LoadBalancer",
  "tags":[
    {
      "key":"cluster",
      "value":"staging"
    },
    {
      "key":"service.k8s.aws/resource",
      "value":"LoadBalancer"
    },
    {
      "key":"service.k8s.aws/stack",
      "value":"ingress-nginx/ingress-nginx-controller"
    }
  ],
  "properties":[
    {
      "name":"CreatedTime",
      "value":"2022-04-07 22:27:39.35 +0000 UTC"
    },
    {
      "name":"DNSName",
      "value":"staging-ingress-14522ba1bd959dd6.elb.us-east-1.amazonaws.com"
    },
    {
      "name":"IpAddressType",
      "value":"ipv4"
    },
    {
      "name":"LoadBalancerArn",
      "value":"arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/net/staging-ingress/14522ba1bd959dd6"
    },
    {
      "name":"LoadBalancerName",
      "value":"staging-ingress"
    }
  ]
}
```

## Mocked data

There is also a mocked API at `http://localhost:8080/mock/api`.  
The mocked api serves static data, it doesn't handle any query parameters.

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resources](http://localhost:8080/mock/api/resources)  | GET  |  Example of "resources" reponse | :white_check_mark: |
| [/tags](http://localhost:8080/mock/api/tags)  | GET  |  Example of "tags" ressponse |  :white_check_mark: |


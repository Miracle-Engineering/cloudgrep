

# API

The backend exposes an API at `http://localhost:8080/api`.

## List resources

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resources](http://localhost:8080/api/resources)  | GET  | Return list of cloud resources |  :white_check_mark: |

| Parameters | Description |  Examples |
| ------------- | ------------- | ------------- |
| tags  | return resource(s) with the provided tag  | `tags[team]=infra` return resources with the tag `team=infra`, meaning "team" with value "infra" <br />`tags[team]=infra&tags[env]=prod` return resources with the tags `team=infra` **and** `env=dev` <br />`tags[env]=prod,dev` return resources with the tags `env=prod` **and** `env=dev` <br />`tags[team]=*` return all the resources with the tag `team` defined|
| tags  | return resource(s) without the provided tag  | `exclude-tags=team` return resources without the tag `team`<br />`exclude-tags=team,env` return resources without the tag `team` **and** `env`


## Mocked data

There is also a mocked API at `http://localhost:8080/mock/api`.  
The mocked api serves static data, it doesn't handle any query parameters.

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resources](http://localhost:8080/mock/api/resources)  | GET  |  Example of "resources" reponse | :white_check_mark: |
| [/tags](http://localhost:8080/mock/api/tags)  | GET  |  Example of "tags" ressponse |  :white_check_mark: |


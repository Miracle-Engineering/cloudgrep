# cloudgrep

## How to use

```shell
# TODO we should not need to import the dependencies each time

# generate the graphql classes
go get github.com/99designs/gqlgen@v0.17.3 && go get github.com/99designs/gqlgen/internal/imports@v0.17.3 && go get github.com/99designs/gqlgen/codegen/config@v0.17.3 && go get github.com/99designs/gqlgen/internal/imports@v0.17.3 && go run github.com/99designs/gqlgen generate
```

```shell
# run the server
go run server.go
```

## Create a resource

```graphql
mutation createResource ($input: NewResource!) {
  createResource(input: $input) {
    id
    region
    type
    tags {
      key
      value
    }
  }
}
```


```json
//TODO set id to be string
{
  "input": {
    "id": 4,
    "region": "us-east-1",
    "type": "EC2-Instance",
    "tags": [
      {
      "key": "env",
      "value": "dev"
      },
      {
      "key": "team",
      "value": "infra"
      }
    ]
  }
}
```



## View a resource

```graphql
query resources {
  resources {
    id
    type
    region
    tags {
      key
      value
    }
  }
}
```

```json
//TODO check all tags returned
{
  "data": {
    "resources": [
      {
        "id": "4",
        "type": "EC2-Instance",
        "region": "us-east-1",
        "tags": [
          {
            "key": "env",
            "value": "dev"
          }
        ]
      }
    ]
  }
}
```
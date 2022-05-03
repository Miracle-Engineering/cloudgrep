# development notes

This application provides a single binary that contains the backend and the frontend.

## Project structure

The project structure is inspired by https://github.com/sosedoff/pgweb

    .
    ├── pkg                                  # Backend: Golang
    └── static                               # Frontend: Static assets
    └── main.go                              # Executable main function
    └── Makefile                             # Defines the tasks to be executed - for local or CI run

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

## API

The backend exposes an API at `http://localhost:8080/api`.

### Routes

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/info](http://localhost:8080/api/info)  | GET  | API information |  :white_check_mark: |
| [/resources](http://localhost:8080/api/resources)  | GET  | Return list of cloud resources |  :x: |
| [/tags](http://localhost:8080/api/tags)  | GET  | Return list of tags |  :x: |

### Mocked data

There is also a mocked API at `http://localhost:8080/mock/api`.  
The mocked api serves static data, it doesn't handle any query parameters.

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resources](http://localhost:8080/mock/api/resources)  | GET  |  Example of "resources" reponse | :white_check_mark: |
| [/tags](http://localhost:8080/mock/api/tags)  | GET  |  Example of "tags" ressponse |  :white_check_mark: |


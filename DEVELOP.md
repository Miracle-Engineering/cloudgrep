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


# mininal-dropbox

This is a minimal dropbox that implements a server rendered static page that is meant to be used when you quickly want to pass files around.

The idea is for this to not be a complex data storage, but just a simple and minimal "dropbox". As it's written in go, it's Docker image size is also small.

## Table of Contents

- [Project Description](#project-description)
- [Features](#features)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Project Description

This project uses Gin to implement a REST API for uploading and downloading files, acting much similar to a dropbox. Deleting the files is only possible by using the endpoint, as the home page is meant to be simple and uses pure HTML.

It also uses Zerolog for structured logging.

## Features

It features three types of object storage, defined under `./storage`:

- `memory`: This storage system saves the files on memory. Be careful as there's no limit on the memory usage.
- `file_system`: This storage system saves the files on the current file system under the folder path defined on `STORAGE_MANAGED_DIR`.
- `nop`: This is a no-operation storage system meant for tests.

If you wish to implement a new storage, all you need is to add a new system definition for `storage.System` and implement the `Storage` interface.

It also features two types of rest system, defined under `./rest`:

- `gin`: This uses the Gin framework to route the endpoints.
- `nop`: This is a no-operation router that is meant for tests.

If you wish to implement a new rest system, all you need is to add a new system definition for `rest.System` and implement the `Server` interface.

## Usage

```shell
$ git clone https://github.com/CG-SS/minimal-dropbox.git
$ cd minimal-dropbox
$ go mod download
$ go mod run .
```

Then access the home page at `http://localhost:12345/`.

You can also run the Dockerfile:

```shell
$ docker build -t minimal-dropbox .
$ docker run -dp 12345:12345 minimal-dropbox
```

## API Endpoints

| Method | Path            | Description                                     |
|--------|-----------------|-------------------------------------------------|
| GET    | /               | Home page for the application. Can be disabled. |
| GET    | /health         | Health route.                                   |
| GET    | /file/all       | Lists all the current files.                    |
| GET    | /file/:filename | Gets the file defined by `filename`.            |
| POST   | /file/upload    | Multi-part form upload for files.               |
| DELETE | /file/:filename | Deletes file defined by `filename`.             |

## Configuration

All the configuration pieces are available using environment variables:

- `LOGGING_ENABLED`: Enables logging.
- `LOGGING_LEVEL`: If logging is enabled, can be set based on `ZerologLevel`. Possible values: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic` and `disabled`.
- `STORAGE_SYSTEM`: The storage system used to store the files. Possible values: `nop`, `file_system` and `memory`.
- `STORAGE_MANAGED_DIR`: If the storage system is `file_system`, this will become the dir where the files will be stored.
- `CORS_ENABLED`: If set to `true`, will use CORS. Possible values: `true` or `false`.
- `CORS_ALLOWED_ORIGINS`: If CORS is enabled, will be used to define the origins.
- `REST_SYSTEM`: The rest system to be used. Possible values: `nop` and `gin`.
- `REST_HOST`: The host to be used by the rest server.
- `REST_PORT`: The port to be used by the rest server.
- `REST_HOME_ROUTE_ENABLED`: If set to `true`, will make a home page available on `/`. Possible values: `true` or `false`.
- `GIN_MODE`: This is an environment variable used by Gin in order to set the debug mode. Set to `release` in order to not have extra logging.

## Contributing

Push a PR, and follow Golang standards.

## License

This project is licensed under the Apache License.
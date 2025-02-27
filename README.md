# Beetle Rest API

This repo represents the REST api for the Beetle platform

- [Beetle REST API](#beetle-rest-api)
  - [Use](#use)
    - [Pre-requisites](#pre-requisites)
    - [Installation](#installation)
    - [Running Locally (for front-end development, etc)](#running-locally-for-front-end-development-etc)
    - [API Documentation](#api-documentation)
    - [Database Administration](#database-administration)
  - [Development](#development)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation-1)
    - [Running Locally](#running-locally)
    - [Commands](#commands)
    - [Configuration](#configuration)
    - [Swagger Documentation](#swagger-documentation)
    - [Authentication](#authentication)
    - [Hot Reload](#hot-reload)
    - [Structure](#structure)
    - [Validation](#validation)
    - [Database Migrations](#database-migrations)

## Use

This section is geared towards someone simply running this backend locally in docker. For installation and running instructions for actual development see the [#Development](#development) section.

### Pre-requisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop)

### Installation

1. Clone this **repository**
2. Make a copy of `.env.example` called `.env` and substitute in relevant entries.
3. Run `source ./.env` to load the environment variables from this file
4. That's it!

### Running Locally (for front-end development, etc)

1. Run `make docker-build`. Anytime you pull changes you'll need to do this.
2. Run `make docker-up`
3. Your server is running and up to date

### API Documentation

Once running locally, the API documentation can be found at [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Database Administration

Once running locally, you can use the Adminer tool [http://localhost:18080](http://localhost:18080) with the following settings (everything except password will be pre-filled if you use [this link](http://localhost:18080/?pgsql=database&username=postgres&db=beetle&ns=public))

- System: **PostgreSQL**
- Server: **database**
- Username: **postgres**
- Password: **postgres**
- Database: **beetle**

## Development

### Prerequisites

- [Go v1.24+](https://golang.org/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)

### Installation

1. Clone this **repository**
2. Make a copy of `.env.example` called `.env` and substitute in relevant entries (Note: jwt-secret can probably stay the same for local dev).
3. Run `source ./.env` to load the environment variables from this file
4. Run `make install`
5. Install Goose by running `go install github.com/pressly/goose/v3/cmd/goose`
6. Install Swag by running `go install github.com/swaggo/swag/cmd/swag@latest`
7. Install Air by running `go install github.com/air-verse/air@latest`
8. Install StaticCheck by running `go install honnef.co/go/tools/cmd/staticcheck@latest`
9. That's it!

### Running Locally

1. Run required services (the database) via `make docker-up-db`
2. Run the server in watch mode (automatically detect changes and restart) locally by running `make watch`

### Commands

See [Makefile](./Makefile) for additional commands for development

### Swagger Documentation

Our API is documented via inline comments using swagger annotation. All endpoints must have this documentation and it should follow the [Declarative Comments Format](https://github.com/swaggo/swag#declarative-comments-format) as outlined in the [swag](https://github.com/swaggo/swag) documentation.

### Authentication

#### User Auth with JWT

Many endpoints require that the user be authenticated. This is done via sending a JWT token (acquired via POST to `/tokens`) in the `Authorization` header in the following format:
`Bearer TOKEN_GOES_HERE`.

### Hot Reload

For an easier development experience this repo support hot reloading via [air](https://github.com/cosmtrek/air). To use simply install air and run `make watch` in the project folder instead of `make`

### Structure

- This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).
- In addition we try to follow a modified version of [Ben Johnson's "Standard Package Layout"](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
  - Specfically we use the `domain` package as our `root` by his definition, and follow his dependency base subpackage rule for implementing services, as well as mocking rules closely.
- It also uses [Go Modules](https://github.com/golang/go/wiki/Modules) which means it **does not** need to reside in your GOPATH.

### Validation

This repo uses the [go-playground/validator](https://github.com/go-playground/validator) package for validation. For a list of struct tags available by default for validation see [here](https://github.com/go-playground/validator#baked-in-validations). For our custom tags see the [internal/validation package](./internal/validation/validator.go).

### Database Migrations

We use the [pressly fork of goose](https://github.com/pressly/goose#goose--) for database migrations.
Our migrations are stored in [assets/db/migrations](./assets/db/migrations). Our test data is stored in [assets/db/seeds](./assets/db/seeds). `goose-up` is automatically run when the server starts, but should you need to manually change the goose state we have the following `make` command wrappers for goose which make it easier to use:

- `make goose-create`
- `make goose-status`
- `make goose-up` and alias `make migrate`
- `make goose-up-one`
- `make goose-down`

To run seeding, which applies migrations with no versioning (should be used for test data):

- `make goose-create-seed`
- `make goose-seed-up`
- `make goose-seed-down`

#### Resetting the database

WARNING: this will erase everything!

If `make goose-down` isn't working and you need a clean slate with a fresh database, you can run the following

```sh
# install postgresql first with `brew install postgresql`
PGPASSWORD=postgres psql -h localhost -p 15432 -U postgres -d beetle -c 'DROP SCHEMA public CASCADE; CREATE SCHEMA public;'
```

Then, set up the database again:

```sh
make goose-up
make goose-seed-up
```

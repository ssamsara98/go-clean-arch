# Go Clean Architecture

Clean Architecture with

- [dipeshdulal](https://github.com/dipeshdulal/clean-gin)
- [wesionaryTEAM](https://github.com/wesionaryTEAM/go_clean_architecture)
- [Gin Web Framework](https://github.com/gin-gonic/gin)

## Linter setup

Need [Python3](https://www.python.org/) to setup linter in git pre-commit hook.

```zsh
make lint-setup
```

---

## Run application

- Setup environment variables

```zsh
cp template.env .env
```

- Update your database credentials environment variables in `.env` file

### Locally

- Run `go run main.go app:serve` to start the server.
- There are other commands available as well. You can run `go run main.go -help` to know about other commands available.

### Using `Docker`

> Ensure Docker is already installed in the machine.

- Start server using command `docker-compose up -d` or `sudo docker-compose up -d` if there are permission issues.

---

## Migration Commands

⚓️ &nbsp; If you want to run the migration runner from the host environment instead of the docker environment; ensure that `sql-migrate` is installed on your local machine.

### Install `sql-migrate`

> You can skip this step if `sql-migrate` has already been installed on your local machine.

**Note:** Starting in Go 1.17, installing executables with `go get` is deprecated. `go install` may be used instead. [Read more](https://go.dev/doc/go-get-install-deprecation)

```zsh
go install github.com/rubenv/sql-migrate/...@latest
```

If you're using Go version below `1.18`

```zsh
go get -v github.com/rubenv/sql-migrate/...
```

### Running migration

Add argument `p=host` after `make` command to run migration commands on local environment

<b>Example:</b>

```zsh
make p=host migrate-up
```

<details>
    <summary>Available migration commands</summary>

| Command               | Desc                                                       |
| --------------------- | ---------------------------------------------------------- |
| `make migrate-status` | Show migration status                                      |
| `make migrate-up`     | Migrates the database to the most recent version available |
| `make migrate-down`   | Undo a database migration                                  |
| `make redo`           | Reapply the last migration                                 |
| `make create`         | Create new migration file                                  |

</details>

---

## Update Dependencies

<details>
    <summary><b>Steps to Update Dependencies</b></summary>
    
1. `go get -u`
2. Remove all the dependencies packages that has `// indirect` from the modules
3. `go mod tidy`
</details>

<details>
    <summary><b>Discovering available updates</b></summary>
    
List all of the modules that are dependencies of your current module, along with the latest version available for each:
```zsh 
go list -m -u all
```

Display the latest version available for a specific module:

```zsh
go list -m -u example.com/theirmodule
```

<b>Example:</b>

```zsh
go list -m -u cloud.google.com/go/firestore
cloud.google.com/go/firestore v1.2.0 [v1.6.1]
```

</details>

<details>
    <summary><b>Getting a specific dependency version</b></summary>
    
To get a specific numbered version, append the module path with an `@` sign followed by the `version` you want:

```zsh
go get example.com/theirmodule@v1.3.4
```

To get the latest version, append the module path with @latest:

```zsh
go get example.com/theirmodule@latest
```

</details>

<details>
    <summary><b>Synchronizing your code’s dependencies</b></summary>
 
```zsh
go mod tidy
```
</details>

---

Hierarchy

```
main
bootstrap
commands
api
  routes
  middleware
  controller
  service
helpers
infrastructure
model
lib
utils
constants
```

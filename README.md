# esquema

[![Go Reference](https://img.shields.io/badge/Go_Reference-pkg.go.dev-007D9C?style=flat-square)](https://pkg.go.dev/github.com/ju4n97/esquema)
[![Release](https://img.shields.io/github/v/release/ju4n97/esquema?style=flat-square\&label=Release)](https://github.com/ju4n97/esquema/releases/latest)
[![CI](https://img.shields.io/github/actions/workflow/status/ju4n97/esquema/ci.yaml?style=flat-square\&label=CI)](https://github.com/ju4n97/esquema/actions/workflows/ci.yaml)

esquema is a declarative API runtime powered by HCL.

Define HTTP routes, validation, SQL, Valkey, Starlark, Go callbacks, and OpenAPI documentation without generating application code.

Manifests are loaded, validated, and compiled at startup, then executed directly at request time.

[Documentation](https://ju4n97.github.io/esquema/) · [Examples](./examples)

## Example

```hcl
server {
  host = "0.0.0.0"
  port = 8080
}

connection "sql" "main" {
  engine = "postgres"
  source = env("DATABASE_URL")
}

route "POST /users" {
  request {
    body {
      field "email" {
        type     = "string"
        format   = "email"
        required = true
      }

      field "name" {
        type     = "string"
        required = true
      }
    }
  }

  step "sql" "create" {
    connection = "main"

    query = <<-SQL
      INSERT INTO users (email, name)
      VALUES (@email, @name)
      RETURNING id, email, name
    SQL

    args = {
      email = ctx.request.body.email
      name  = ctx.request.body.name
    }
  }

  respond {
    status = 201
    body   = steps.create.row
  }
}
```

## Go

esquema is also embeddable:

```go
config, err := esquema.Load("routes/*.hcl")
if err != nil {
    log.Fatal(err)
}

engine, err := esquema.New(config)
if err != nil {
    log.Fatal(err)
}
defer engine.Close()

http.ListenAndServe(":8080", engine)
```

Custom Go behavior can be registered with `esquema.WithStep`. More information available in the [Go integration guide](https://ju4n97.github.io/esquema/guides/go).

## Install

```bash
go install github.com/ju4n97/esquema/cmd/esquema@latest
```

Or use the release binaries and container images documented in the [installation guide](https://ju4n97.github.io/esquema/installation).

## Documentation

See [esquema documentation](https://ju4n97.github.io/esquema/).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)

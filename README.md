# es-go

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/pascalallen/es-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pascalallen/es-go)](https://goreportcard.com/report/github.com/pascalallen/es-go)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/pascalallen/es-go/go.yml)
![GitHub](https://img.shields.io/github/license/pascalallen/es-go)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/pascalallen/es-go)

![Logo](web/static/logo.svg)

es-go is open source software offering a set of guidelines to use Event Sourcing and Domain-Driven Design in Go. 
This is a living project and is subject to change.

## Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Development Environment Setup

### Clone Repository

```bash
cd <projects-parent-directory> && git clone https://github.com/pascalallen/es-go.git
```

### Copy & Modify `.env` File

```bash
cp .env.example .env
```

### Bring Up Environment

```bash
bin/up
``` 

### Take Down Environment

```bash
bin/down
```

## Testing

Run tests and create coverage profile:

```bash
bin/exec go test ./... -covermode=count -coverprofile=coverage.out
```

Generate HTML file to view test coverage profile:

```bash
bin/exec go tool cover -html=coverage.out -o coverage.html
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](LICENSE)

# dotenv

![build status](https://travis-ci.org/fairyhunter13/dotenv.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/fairyhunter13/dotenv)](https://goreportcard.com/report/github.com/fairyhunter13/dotenv)
![godocs](https://godoc.org/github.com/fairyhunter13/dotenv?status.svg)

A Go (golang) implementation of dotenv _(inspired by: [https://github.com/joho/godotenv](https://github.com/joho/godotenv))_.

## Installation

As a **Library**:

```sh
go get github.com/fairyhunter13/dotenv
```

## Usage

In your environment file (canonically named `.env`):

```sh
S3_BUCKET=YOURS3BUCKET
SECRET_KEY=YOURSECRETKEYGOESHERE

MESSAGE="A message containing important spaces."
ESCAPED='You can escape you\'re strings too.'

# A comment line that will be ignored
GIT_PROVIDER=github.com
LIB=${GIT_PROVIDER}/fairyhunter13/dotenv # variable interpolation (plus ignored trailing comment)
```

In your application:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/fairyhunter13/dotenv"
)

func main() {
  err := dotenv.Load()
  if err != nil {
    log.Fatalf("Error loading .env file: %v", err)
  }

  s3Bucket := os.Getenv("S3_BUCKET")
  secretKey := os.Getenv("SECRET_KEY")

  fmt.Println(os.Getenv("MESSAGE"))
}
```

## Documentation

[https://godoc.org/github.com/fairyhunter13/dotenv](https://godoc.org/github.com/fairyhunter13/dotenv)

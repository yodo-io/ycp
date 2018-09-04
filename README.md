# Yodo Cloud Platform

Your next best cloud provider.

## Prerequisites

- Go and [dep](https://golang.github.io/dep/)
- Make (optional)

## Usage

- Run all tests: `make test`
- Clean artifacts and build: `make clean build`
- Start server with sample data: `make clean run`

## API

```sh
# Will fail with 401 Unauthorized
curl -H"Token: $TOKEN" localhost:9000/v1/resources/1

# Get token (for jq, see see https://stedolan.github.io/jq/)
# Without jq, just copy & paste the token from the JSON response
TOKEN=`curl localhost:9000/auth/token \
  -XPOST \
  -H 'Content-type: application/json' \
  -d '{"email":"joe@example.org","password":"secret"}' | jq -r '.token'`

# Get users
curl -H"Token: $TOKEN" localhost:9000/v1/users/1

# Get resources for user with id 1
curl -H"Token: $TOKEN" localhost:9000/v1/resources/1

# Get resources for user with id 2 - will fail with 403 Forbidden
curl -H"Token: $TOKEN" localhost:9000/v1/resources/2

# List catalog
curl -H"Token: $TOKEN" localhost:9000/v1/catalog

# List quotas for user with id 1
curl -H"Token: $TOKEN" localhost:9000/v1/quotas/1
```

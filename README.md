# 7 Solutions Backend Challenge

## Project Setup

To start the project, run the following command:

```bash
docker compose up
```

## JWT Token Usage Guide

### 1. Obtain Tokens

- **Endpoint:** `POST /sessions`
- **Response:** Returns `access_token` and `refresh_token`.

### 2. Accessing and Updating Data

- Use the `access_token` in the request header to access or update data.
- Example header:

    ```json
    {
        "Authorization": "Bearer <access_token>"
    }
    ```

### 3. Refreshing Tokens

- **Endpoint:** `POST /token/refresh`
- Use the `refresh_token` to obtain a new `access_token` when the current one is expired.

## Sample API requests/responses

After running `docker compose up` you can access `http://localhost:8888/docs/index.html` for Swagger/OpenAPI

## Sample gRPC requests/response

### Sample gRPC Service

The sample gRPC service documentation can be found in the `grpc-docs` directory. Refer to the provided `.proto` files and examples for details on available services, methods, and message structures.

To test the gRPC service, you can use tools like `grpcurl` or any gRPC client library in your preferred programming language.

## Mocking

Run command inside `makefile` which is `make gen_mocks`

## Unit Testing

Run command  `make unit_test`. If all test pass you can see percent coverage at `coverage.html`.

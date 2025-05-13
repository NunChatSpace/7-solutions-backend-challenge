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

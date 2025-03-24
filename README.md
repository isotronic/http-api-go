# HTTP API Server

This project is an HTTP API server written in Go. It provides a simple API for managing resources.

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

1. Clone the repository:

```sh
git clone https://github.com/isotronic/http-api-go.git
```

2. Navigate to the project directory:

```sh
cd http-api-go
```

3. Create a `.env` file with the following environment variables:

```sh
DB_URL="postgres://postgres:password@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
TOKEN_SECRET="your_token_secret"
POLKA_KEY="your_polka_key"
```

4. Build the project:

```sh
go build
```

5. Run the server:

```sh
./http-api-go
```

## API Documentation

### Base URL

```
http://localhost:8080/
```

### Endpoints

#### Health Check

- **URL:** `/api/healthz`
- **Method:** `GET`
- **Description:** Checks the health of the API.
- **Response:**
  - **Status:** `200 OK`
  - **Body:** `OK`

#### Get All Chirps

- **URL:** `/api/chirps`
- **Method:** `GET`
- **Description:** Retrieves a list of all chirps.
- **Response:**

  ```json
  [
    {
      "id": "uuid",
      "user_id": "uuid",
      "body": "Chirp body",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
  ```

#### Get Chirp by ID

- **URL:** `/api/chirps/{chirpID}`
- **Method:** `GET`
- **Description:** Retrieves a single chirp by its ID.
- **Response:**

  ```json
  {
    "id": "uuid",
    "user_id": "uuid",
    "body": "Chirp body",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```

#### Create Chirp

- **URL:** `/api/chirps`
- **Method:** `POST`
- **Description:** Creates a new chirp.
- **Request Body:**

  ```json
  {
    "body": "Chirp body"
  }
  ```

- **Headers:**
  - `Authorization: Bearer <access_token>`
- **Response:**

  ```json
  {
    "id": "uuid",
    "user_id": "uuid",
    "body": "Chirp body",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```

#### Delete Chirp

- **URL:** `/api/chirps/{chirpID}`
- **Method:** `DELETE`
- **Description:** Deletes a chirp by its ID.
- **Headers:**
  - `Authorization: Bearer <access_token>`
- **Response:**
  - **Status:** `204 No Content`

#### Create User

- **URL:** `/api/users`
- **Method:** `POST`
- **Description:** Creates a new user.
- **Request Body:**

  ```json
  {
    "email": "user@example.com",
    "password": "password"
  }
  ```

- **Response:**

  ```json
  {
    "id": "uuid",
    "email": "user@example.com",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "is_chirpy_red": false
  }
  ```

#### Update User

- **URL:** `/api/users`
- **Method:** `PUT`
- **Description:** Updates an existing user.
- **Request Body:**

  ```json
  {
    "email": "newemail@example.com",
    "password": "newpassword"
  }
  ```

- **Headers:**
  - `Authorization: Bearer <access_token>`
- **Response:**

  ```json
  {
    "id": "uuid",
    "email": "newemail@example.com",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "is_chirpy_red": false
  }
  ```

#### Login

- **URL:** `/api/login`
- **Method:** `POST`
- **Description:** Logs in a user and returns access and refresh tokens.
- **Request Body:**

  ```json
  {
    "email": "user@example.com",
    "password": "password"
  }
  ```

- **Response:**

  ```json
  {
    "id": "uuid",
    "email": "user@example.com",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "token": "access_token",
    "refresh_token": "refresh_token",
    "is_chirpy_red": false
  }
  ```

#### Refresh Token

- **URL:** `/api/refresh`
- **Method:** `POST`
- **Description:** Refreshes the access token using a refresh token.
- **Headers:**
  - `Authorization: Bearer <refresh_token>`
- **Response:**

  ```json
  {
    "token": "new_access_token"
  }
  ```

#### Revoke Token

- **URL:** `/api/revoke`
- **Method:** `POST`
- **Description:** Revokes a refresh token.
- **Headers:**
  - `Authorization: Bearer <refresh_token>`
- **Response:**
  - **Status:** `204 No Content`

#### Polka Webhooks

- **URL:** `/api/polka/webhooks`
- **Method:** `POST`
- **Description:** Handles Polka webhooks.
- **Request Body:**

  ```json
  {
    "event": "user.upgraded",
    "data": {
      "user_id": "uuid"
    }
  }
  ```

- **Headers:**
  - `Authorization: ApiKey <polka_key>`
- **Response:**
  - **Status:** `204 No Content`

#### Admin Metrics

- **URL:** `/admin/metrics`
- **Method:** `GET`
- **Description:** Retrieves admin metrics.
- **Response:**
  - **Status:** `200 OK`
  - **Body:** HTML content with metrics

#### Admin Reset

- **URL:** `/admin/reset`
- **Method:** `POST`
- **Description:** Resets the database (only in dev mode).
- **Response:**
  - **Status:** `200 OK`

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.

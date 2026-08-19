# API Documentation

Base URL (local): `http://localhost:8080`
Base URL (production): set via your deployed backend's URL.

All request/response bodies are JSON. Authenticated endpoints rely on an
`auth_token` httpOnly cookie set by `/api/register` or `/api/login` --
no `Authorization` header is used.

---

## Auth

### `POST /api/register`

Creates a new user and logs them in (sets the `auth_token` cookie).

**Request body**
```json
{
  "first_name": "Alex",
  "last_name": "N",
  "email": "alex@example.com",
  "password": "secret123"
}
```

**Response `200`**
```json
{
  "id": 1,
  "first_name": "Alex",
  "last_name": "N",
  "email": "alex@example.com",
  "created_at": "2026-08-19T10:00:00Z"
}
```

**Errors**
- `400` -- invalid request body or password hashing failure
- `409` -- a user with this email already exists

---

### `POST /api/login`

Verifies email/password and logs the user in (sets the `auth_token` cookie).

**Request body**
```json
{
  "email": "alex@example.com",
  "password": "secret123"
}
```

**Response `200`** -- same shape as register.

**Errors**
- `400` -- invalid request body
- `401` -- invalid email or password (deliberately the same message for
  both cases, so the endpoint can't be used to enumerate registered emails)

---

### `POST /api/logout`

Clears the `auth_token` cookie. No request body.

**Response `200`** -- empty body.

---

## Shortener

### `POST /api/shorten`

Creates a shortened URL. Works for both anonymous and logged-in
requests -- if a valid `auth_token` cookie is present, the created URL
is associated with that user; otherwise it's created with no owner.

**Request body**
```json
{ "url": "https://example.com/some/long/path" }
```

**Response `200`**
```json
{
  "short_code": "a1b2c3",
  "short_url": "http://localhost:8080/a1b2c3",
  "original_url": "https://example.com/some/long/path"
}
```

**Errors**
- `400` -- missing/invalid URL (must start with `http://` or `https://`)

---

### `GET /{code}`

Redirects to the original URL for a given short code. Logs a click
(asynchronously) against that URL.

**Response `302`** with a `Location` header pointing at the original URL.

**Errors**
- `404` -- short code not found

---

### `GET /api/my/urls`

Returns every shortened URL belonging to the logged-in user. **Requires
authentication** (a valid `auth_token` cookie) -- unlike `/api/shorten`,
this endpoint rejects anonymous requests.

**Response `200`**
```json
[
  {
    "short_code": "a1b2c3",
    "short_url": "http://localhost:8080/a1b2c3",
    "original_url": "https://example.com/some/long/path",
    "created_at": "2026-08-19T10:00:00Z"
  }
]
```

Returns `[]` (not `null`) if the user has no URLs yet.

**Errors**
- `401` -- no valid session

---

## Analytics

### `GET /api/stats/{code}`

Returns click statistics for a given short code.

**Response `200`**
```json
{
  "short_code": "a1b2c3",
  "total_clicks": 12,
  "timeline": [
    { "date": "2026-08-18", "count": 5 },
    { "date": "2026-08-19", "count": 7 }
  ]
}
```

`timeline` is `[]` (not `null`) if there are no clicks yet.

---

## Health

### `GET /health`

Plain liveness check, returns `200` with body `ok`. No database access
-- useful for confirming the server process itself is up, separate
from whether the database connection is working.

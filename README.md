# URL Shortener

A full-stack URL shortener built as a solo capstone project for the
Chingu Voyage prerequisite. Backend in Go, frontend in React + TypeScript.

I actually build this project as a challenge after learning how to design a url shortening service. The shortening service uses hashing with collision resolution strategy and not base62 encode conversion due to a number of limitations.

Paste a link, get a short one back -- no account required. Sign up to
see a dashboard of your own shortened links and how many clicks each
one has gotten over time.

## Features

- Shorten a URL, with or without an account
- Anonymous or logged-in shortening (URLs are attributed to your
  account only if you're logged in when you create them)
- Register / log in / log out, with hashed passwords and JWT-based
  sessions (httpOnly cookie, not stored in localStorage)
- Per-user dashboard listing your shortened URLs
- Click tracking per URL: total count + a per-day timeline
- Share a shortened link directly to WhatsApp, X, Facebook, email, or
  your device's native share sheet

## Tech stack

**Backend:** Go, [chi](https://github.com/go-chi/chi) (routing),
[pgx](https://github.com/jackc/pgx) (Postgres driver), PostgreSQL,
[golang-migrate](https://github.com/golang-migrate/migrate) (schema
migrations), [golang-jwt](https://github.com/golang-jwt/jwt) (JWT),
bcrypt (password hashing)

**Frontend:** React, TypeScript, Vite, Tailwind CSS v4,
[Bun](https://bun.sh) (package manager / dev runtime),
react-router-dom

**Infra:** Docker + Docker Compose (local dev), GitHub Actions (CI/CD),
Render (backend hosting), Netlify (frontend hosting), Neon (production
Postgres)

## Architecture

The backend is organized by domain rather than by technical layer --
each domain (`shortener`, `auth`, `analytics`) owns its own model,
repository, service, and HTTP handler:

```
backend/
├── cmd/api/main.go              # wires everything together
├── internal/
│   ├── config/                  # env var loading
│   ├── database/
│   │   └── migrations/          # golang-migrate .up/.down SQL files
│   ├── httpserver/               # shared router + CORS middleware
│   ├── shortener/                # url creation, redirect, per-user url list
│   ├── auth/                     # register, login, logout, JWT, middleware
│   └── analytics/                # click logging, stats aggregation
```

Within each domain:
- **model** -- plain data structs, no logic
- **repository** -- the only place that touches SQL; exposed as an
  interface so it can be swapped/mocked
- **service** -- business rules (validation, hashing, orchestration)
- **handler** -- HTTP-only concerns: parsing requests, picking status
  codes, shaping JSON responses

See [`docs/API.md`](docs/API.md) for the full endpoint reference.

## Local development

### Prerequisites
- Go 1.23+
- [Bun](https://bun.sh)
- Docker + Docker Compose

### Setup

1. Clone the repo.
2. Create a `.env` file at the repo root (see `.env.example`) with:
   ```
   PORT=8080
   DB_USER=...
   DB_PASSWORD=...
   DB_NAME=url-shortener
   DATABASE_URL=postgres://...@postgres_db:5432/url-shortener?sslmode=disable
   JWT_SECRET=...
   BASE_URL=http://localhost:8080
   FRONTEND_URL=http://localhost:5173
   APP_ENV=development
   ```
3. Create `frontend/.env` with:
   ```
   VITE_API_URL=http://localhost:8080
   ```
4. Start Postgres + run migrations + start the backend:
   ```
   docker compose up --build
   ```
   Migrations run automatically via a dedicated `migrate` service
   before the backend starts.
5. In a separate terminal, start the frontend:
   ```
   cd frontend
   bun install
   bun run dev
   ```
6. Visit `http://localhost:5173`.

## Deployment

- **Backend** deploys to Render as a Docker service. Database
  migrations run via GitHub Actions against the production database
  before each deploy is triggered (see `.github/workflows/deploy.yml`).
- **Frontend** deploys to Netlify as a static build, triggered via a
  Netlify build hook in the same workflow.
- Production Postgres is hosted on [Neon](https://neon.tech). Neon's
  pooled connection strings require an `options=endpoint=<endpoint-id>`
  query parameter for compatibility with the `migrate` CLI's Postgres
  driver (see their [SNI docs](https://neon.tech/sni) if you hit
  `Endpoint ID is not specified`).

## Known limitations

Given the project timeline, a few things were deliberately scoped out
or simplified rather than left half-built:

- **No session-check endpoint.** There's no `GET /api/me`, so the
  frontend can't verify an existing session on page load. A logged-in
  user who refreshes the page will appear logged out (even though
  their cookie is still valid) until they log in again.
- **Short code collisions** are handled with a single prefixed retry
  (not a full retry loop), so a second-level collision is theoretically
  possible, if very unlikely at this scale.
- **No password reset / email verification** -- out of scope for the
  capstone timeline.

## License

MIT

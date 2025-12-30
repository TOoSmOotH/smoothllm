# Repository Guidelines

This guide summarizes how to work in this repo without re-reading the full docs.

## Project Structure & Module Organization
- `backend/` is the Go API. Entry point: `backend/cmd/api/main.go`. Core logic lives in `backend/internal`, configs in `backend/configs`, and persistent data in `backend/data` and `backend/uploads`.
- `frontend/` is the Vue 3 + TypeScript app. Source code is under `frontend/src`, static assets in `frontend/public` and `frontend/src/assets`, and tests in `frontend/tests`.
- `docker-compose.yml` wires backend and frontend services for local development.

## Build, Test, and Development Commands
- `docker-compose up --build`: run full stack locally with Docker (see `DOCKER.md`).
- `cd backend && go run ./cmd/api`: run the API directly without Docker.
- `cd frontend && npm run dev`: start the Vite dev server.
- `cd frontend && npm run build`: type-check and build production assets.
- `cd frontend && npm run test` / `npm run test:e2e`: run Vitest unit tests or Playwright E2E tests.
- `cd backend && go test ./...`: run Go unit and integration tests.

## Coding Style & Naming Conventions
- `.editorconfig` sets 2-space indentation by default; Go uses tabs.
- Go: format with `gofmt`, follow Effective Go, keep exported symbols documented.
- Vue: use `<script setup lang="ts">`, Composition API, and PascalCase component names in `frontend/src/components`.
- Composables should be `useXxx` in `frontend/src/composables`.

## Testing Guidelines
- Backend tests live in `backend/tests/unit` and `backend/tests/integration`; use `*_test.go` and `go test ./...`.
- Frontend tests live in `frontend/tests/unit` (Vitest) and `frontend/tests/e2e` (Playwright).
- Target coverage in `CONTRIBUTING.md`: backend >= 85%, frontend >= 80%.

## Commit & Pull Request Guidelines
- No git history is available in this workspace; follow Conventional Commits per `CONTRIBUTING.md`.
- PRs should include a clear description, linked issues, and screenshots for UI changes.
- Use templates if present (see `GITHUB_TEMPLATE_GUIDE.md`).

## Security & Configuration Tips
- Copy `.env.example` to `.env` in `backend/` and `frontend/` before running locally.
- Keep secrets out of git; use environment variables or Docker Compose overrides.

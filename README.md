# SmoothWeb Template

SmoothWeb is a production-ready full-stack template with Go + Vue, user management, and RBAC out of the box.

## Quick Start

```sh
docker-compose up --build
```

Or run services directly:

```sh
cd backend && go run ./cmd/api
cd frontend && npm run dev
```

## Customize Your App

- Frontend configuration: `frontend/src/custom/appConfig.ts`
- Frontend routes: `frontend/src/custom/routes.ts`
- Backend routes: `backend/internal/custom/routes.go`
- Customization guide: `docs/CUSTOMIZATION.md`

## Project Structure

- `backend/` Go API (`backend/cmd/api/main.go`)
- `frontend/` Vue 3 + TypeScript (`frontend/src`)
- `docker-compose.yml` local dev services

## Template Notes

- `node_modules/` is ignored in `.gitignore` and should not be committed.
- Use `CHANGELOG.md` to record template updates.

## Update Flow (Downstream Projects)

Keep your project in sync with the template by adding the template repo as a remote and merging updates:

```sh
git remote add template https://github.com/USERNAME/smoothweb.git
git fetch template main
git merge template/main
```

If you use a different default branch, replace `main` with your branch name.

# Mini CI - Build Orchestrator (WIP)

[![CI](https://github.com/H3nSte1n/ci-orchestrator/actions/workflows/ci.yml/badge.svg)](https://github.com/H3nSte1n/ci-orchestrator/actions/workflows/ci.yml)
[![Coverage Status](https://codecov.io/gh/H3nSte1n/ci-orchestrator/branch/main/graph/badge.svg)](https://codecov.io/gh/H3nSte1n/ci-orchestrator)

A minimal build orchestration system in Go: submit a build, queue it, execute it on workers, persist logs, and track results.

This repository explores practical problems in build infrastructure. Job lifecycle, concurrency control, isolation boundaries, caching opportunities, and operational visibility, through a small, end-to-end implementation.

> **Status:** Work in progress. Core API + database model are implemented. Worker executes commands and persists logs. Git clone/checkout and SSE log streaming are in progress.

---

## Communication flow
![mini-ci flow](docs/flow.png)

---

## Overview

`mini-ci` consists of two services:

- **API / Scheduler**  
  Accepts build requests, persists job state, supports cancellation, and exposes build status.

- **Worker**  
  Claims queued builds, creates a workspace, executes the build command (host runner for now), and persists logs + results.

---

## Features

### Implemented
- Build lifecycle persisted in Postgres (UUID IDs)
- Endpoints:
  - `POST /api/v1/builds` — create a build job
  - `GET /api/v1/builds/:id` — fetch job state
  - `POST /api/v1/builds/:id/cancel` — request cancellation
  - `PATCH /api/v1/builds/:id/status` — update status *(development endpoint, will be restricted/removed)*

- Migrations:
  - `builds` table (job state + locking fields)
  - `build_logs` table (persistent logs per build)

- Worker: claim + execute (host runner) + complete builds
- Persist logs (stdout/stderr) to build_logs

### In progress
- Git clone + checkout ref (workspace from repo)
- Client log streaming (SSE)
- Container runner adapter (Docker/Podman) + resource limits
- Artifact upload (local -> S3/MinIO)
- Cache restore/save with content-addressed keys
- Reliability: heartbeats + stuck-job recovery + retries

---

## Architecture

The codebase follows a **ports & adapters (hexagonal)** layout:

- **Core (domain/use-cases)** is infrastructure-agnostic
- **Adapters** implement HTTP, Postgres, queue/locking, runners, and log streaming

This keeps the job lifecycle logic testable and allows swapping infrastructure (e.g., DB polling -> Redis/NATS, host exec -> Docker/Podman) without rewriting core behavior.

---

## Running locally

### Requirements
- docker
- or run services locally with Go (requires Postgres)

### Start
```
docker-compose up --build
```

Services:
 - API: `http://localhost:8000`
 - Worker: `http://localhost:8001`
 - Postgres: `localhost:5432`

> Development containers run with live reload (air).

## Roadmap
- [x] Worker: claim queued jobs safely and execute commands (host runner)
- [x] Persist logs to DB
- [ ] Stream logs (SSE)
- [ ] Container runner adapter (Docker/Podman) with resource limits
- [ ] Artifact upload (local -> S3)
- [ ] Cache restore/save (content-addressed keys)
- [ ] Heartbeats, retries, stuck-job recovery

## License
MIT License. See [LICENSE](LICENSE) for details.


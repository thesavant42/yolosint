# Docker Dorker - OCI ICU 

Forked from and Inspired by `oci.dag.dev`

## Docker Dorker
This is a web server for exploring the contents of an OCI registry, FS style output.

## Quick Start (Docker)

```bash
git pull
docker compose up --build
```

Service runs at `http://localhost:8042`

### When to use `--build`

| Scenario | Command |
|----------|---------|
| First run | `docker compose up --build` |
| After pulling changes | `docker compose up --build` |
| Restarting (no code changes) | `docker compose up` |
| Stop | `docker compose down` |

## Running Locally (without Docker)

```bash
git pull
rm ./oci
go build ./cmd/oci
./oci -v
```

Opens a listener on port localhost:8080 but I forward this to port 8042 in Docker. That aligns with the url on line 16 of [DockerHubOCI Explorer-1.1.user.js](/DockerHubOCI Explorer-1.1.user.js)
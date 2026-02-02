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

Opens a listener on port *localhost:8080*, **but I forward this to port 8042** in Docker.
    - That aligns with the url on line 16 of [DockerHubOCI yolosint.user.js](/DockerHubOCI Explorer-1.1.user.js)

## Cached Files
Stored in the [./cache/](./cache/) folder.

## User Script
- `yolosint.user.js` is a user script written for Tampermonkey, ymmv.
- Adds an overlay when viewing docker hub to quickly view a container image in docker dorker. 
- To use after installing, simply browse an image in Docker Hub and click the "Tags" tab to view details about recent tags for that Repository.
- The overlay will be displayed nexted to the Tag ID.



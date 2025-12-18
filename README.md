# YOLOSINT
DO NOT DELETE TASKS.MD
A rewrite of gitsome-ng using the TUIOS windowing system.

## Purpose

This project solves the problems of the previous version: the UI/UX was buggy and full of contradictions. We are rebuilding with modular coding principles and a proper windowing system.

## Task List

See **[TASKS.md](TASKS.md)** for the complete project task list with subtasks.

## Development Environment

- GoLang on Windows via PowerShell Terminal
- Use `;` semicolons to chain commands (NOT `&&`)

## Project Structure

```
yolosint/
├── cmd/
│   └── yolosint/
│       └── main.go                 # Entry point, TUIOS bootstrap
│
├── internal/
│   ├── config/
│   │   └── env.go                  # .env loading, secrets
│   │
│   ├── style/
│   │   └── stylebible.go           # TUIOS theme constants
│   │
│   ├── app/
│   │   ├── app.go                  # TUIOS app shell, window manager
│   │   └── navigation.go           # Cross-feature navigation
│   │
│   ├── project/
│   │   ├── window.go               # Project selector window
│   │   ├── db.go                   # Project CRUD
│   │   ├── models.go               # Project model
│   │   │
│   │   ├── github/                 # GitHub (peer of domain, poi)
│   │   │   ├── window.go           # GitHub window
│   │   │   ├── db.go               # GitHub data storage
│   │   │   ├── models.go           # Commit, user, repo models
│   │   │   ├── commits/
│   │   │   ├── profile/
│   │   │   └── repos/
│   │   │
│   │   ├── domain/                 # Domain (peer of github, docker, poi, report)
│   │   │   ├── window.go           # Domain management window
│   │   │   ├── db.go               # Domain CRUD
│   │   │   ├── models.go           # Domain model
│   │   │   ├── wayback/            # Wayback (child of domain)
│   │   │   └── subdomain/          # Subdomain (child of domain)
│   │   │       ├── virustotal/
│   │   │       └── crtsh/
│   │   │
│   │   ├── poi/                    # Person of Interest (peer of domain)
│   │   │   ├── window.go           # POI management window
│   │   │   ├── db.go               # POI CRUD
│   │   │   ├── models.go           # POI model
│   │   │   ├── email/              # Emails (child of POI)
│   │   │   ├── social/             # Social media profiles (child of POI)
│   │   │   ├── links/              # Personal links (child of POI)
│   │   │   └── notes/              # Notes (child of POI)
│   │   │
│   │   ├── docker/                 # Docker (peer of domain)
│   │   │   ├── window.go           # LayerSlayer window
│   │   │   ├── db.go               # Docker data storage
│   │   │   ├── models.go           # Image, tag models
│   │   │   ├── hub/
│   │   │   ├── registry/
│   │   │   └── layers/
│   │   │
│   │   └── report/                 # Report (peer of domain)
│   │       ├── window.go           # Report export window
│   │       └── export.go           # Report generation
│   │
│   └── (no orphan modules - all features under project)
│
├── docs/
│   ├── includes/
│   │   ├── gitsome-ng/             # Reference code (DATA ONLY)
│   │   └── tuios/                  # TUIOS documentation
│   └── testing/
│       └── fixtures/               # Test data files
│
├── .env                            # Secrets (gitignored)
├── env.example
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── TASKS.md
```

## Data Hierarchy

```
Project
  ├── GitHub              (peer of Domain, POI - repos exist independently)
  │     ├── Commits
  │     ├── Profile
  │     └── Repos
  │
  ├── Domain
  │     ├── Wayback       (child of Domain - requires domain)
  │     └── Subdomain     (child of Domain - requires domain)
  │           ├── VirusTotal
  │           └── crt.sh
  │
  ├── POI                 (peer of Domain - Person of Interest)
  │     ├── Email         (child of POI - can be any domain)
  │     ├── Social        (child of POI - social media profiles)
  │     ├── Links         (child of POI - personal websites, portfolios)
  │     └── Notes         (child of POI - freeform investigation notes)
  │
  ├── Docker              (peer of Domain)
  │     ├── Hub
  │     ├── Registry
  │     └── Layers
  │
  └── Report              (peer of Domain - stored in project DB)
```

## Documentation Reference

| Source | Path | Purpose |
|--------|------|---------|
| TUIOS Docs | `docs/includes/tuios/docs/` | Windowing system documentation |
| TUIOS Library | `docs/includes/tuios/pkg/tuios/` | Library API reference |
| gitsome-ng API | `docs/includes/gitsome-ng/internal/api/` | API clients to port |
| gitsome-ng DB | `docs/includes/gitsome-ng/internal/db/` | Database layer to port |
| gitsome-ng Models | `docs/includes/gitsome-ng/internal/models/` | Data models to port |
| Test Fixtures | `docs/testing/fixtures/` | Sample data for testing |

## Environment Variables

Required secrets in `.env`:

| Variable | Purpose |
|----------|---------|
| `GITHUB_TOKEN` | GitHub API authentication |
| `VIRUS_TOTAL_TOKEN` | VirusTotal API key |
| `DOCKERHUB_USER` | Docker Hub username |
| `DOCKERHUB_TOKEN` | Docker Hub access token |

## Constraints

### Style Rules
- No emojis anywhere (not in code, not in communications)
- Succinct code is best code
- No file should exceed 100 lines of code (refactor into submodules)
- Must use named constants (no magic numbers)

### Linting Rules
- Never finish a task with linter messages
- Do not delete linter errors; fix them
- Unused variables indicate incomplete refactoring; investigate before removing

### Task Rules
- Never complete a task that has not been tested
- Compilation does not mean it works
- Only the user can mark a task as complete
- Testing must occur in collaboration with the user

### What NOT to Port
- UI/UX elements from gitsome-ng (rebuild from scratch with TUIOS)
- Command line arguments (focus on TUI experience)

## Project Philosophy

This project follows the "fail fast" methodology:
- Quickly test ideas to identify flaws early
- Learn from mistakes cheaply
- Pivot or refine before significant investment
- Fast iteration and hypothesis testing

Elegant code is required:
- Modular and reusable
- Functional and legible
- No shortcuts

## Data Storage

- SQLite databases created in and loaded from current working directory
- Multiple project variations can run concurrently
- Reports exported to current working directory
- Report naming format: `PROJECT-Reportname-DDMMYYYY-HMS`

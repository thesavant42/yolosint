# YOLOSINT Task List DO NOT DELETE

Port gitsome-ng functionality to yolosint using TUIOS windowing system.

---

## CRITICAL: UI/UX CONTAMINATION WARNING

When porting from gitsome-ng:
1. Extract ONLY the pure data logic (structs, API calls, DB queries)
2. Strip ALL rendering, styling, and user interaction code
3. Build fresh UI using TUIOS patterns

**FORBIDDEN:** `docs/includes/gitsome-ng/internal/ui/` - Do not reference or import

---

## Project Structure

Hierarchical layout reflecting data model.

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
│   │   │   │
│   │   │   ├── commits/
│   │   │   │   ├── window.go       # Commits view
│   │   │   │   ├── api.go          # Commits API
│   │   │   │   └── db.go           # Commits queries
│   │   │   │
│   │   │   ├── profile/
│   │   │   │   ├── window.go       # User profile view
│   │   │   │   ├── api.go          # User API
│   │   │   │   └── db.go           # Profile queries
│   │   │   │
│   │   │   └── repos/
│   │   │       ├── window.go       # Repos list view
│   │   │       ├── api.go          # Repos API
│   │   │       └── db.go           # Repos queries
│   │   │
│   │   ├── domain/                 # Domain (peer of github, docker, poi, report)
│   │   │   ├── window.go           # Domain management window
│   │   │   ├── db.go               # Domain CRUD
│   │   │   ├── models.go           # Domain model
│   │   │   │
│   │   │   ├── wayback/            # Wayback (child of domain)
│   │   │   │   ├── window.go       # Wayback Machine window
│   │   │   │   ├── api.go          # CDX API client
│   │   │   │   ├── db.go           # Wayback storage
│   │   │   │   └── models.go       # CDX record models
│   │   │   │
│   │   │   └── subdomain/          # Subdomain (child of domain)
│   │   │       ├── window.go       # Subdomonster window
│   │   │       ├── db.go           # Subdomain storage
│   │   │       ├── models.go       # Subdomain models
│   │   │       │
│   │   │       ├── virustotal/
│   │   │       │   └── api.go      # VirusTotal client
│   │   │       │
│   │   │       └── crtsh/
│   │   │           └── api.go      # crt.sh client
│   │   │
│   │   ├── poi/                    # Person of Interest (peer of domain)
│   │   │   ├── window.go           # POI management window
│   │   │   ├── db.go               # POI CRUD
│   │   │   ├── models.go           # POI model
│   │   │   │
│   │   │   ├── email/              # Emails (child of POI)
│   │   │   │   ├── window.go       # Email list window
│   │   │   │   ├── db.go           # Email storage
│   │   │   │   └── models.go       # Email model (can be any domain)
│   │   │   │
│   │   │   ├── social/             # Social media profiles (child of POI)
│   │   │   │   ├── window.go       # Social profiles window
│   │   │   │   ├── db.go           # Social profile storage
│   │   │   │   └── models.go       # Platform, handle, URL, etc.
│   │   │   │
│   │   │   ├── links/              # Personal links (child of POI)
│   │   │   │   ├── window.go       # Links window
│   │   │   │   ├── db.go           # Links storage
│   │   │   │   └── models.go       # URL, title, description
│   │   │   │
│   │   │   └── notes/              # Notes (child of POI)
│   │   │       ├── window.go       # Notes window
│   │   │       ├── db.go           # Notes storage
│   │   │       └── models.go       # Freeform notes, timestamps
│   │   │
│   │   ├── docker/                 # Docker (peer of domain)
│   │   │   ├── window.go           # LayerSlayer window
│   │   │   ├── db.go               # Docker data storage
│   │   │   ├── models.go           # Image, tag models
│   │   │   │
│   │   │   ├── hub/
│   │   │   │   └── api.go          # Docker Hub API
│   │   │   │
│   │   │   ├── registry/
│   │   │   │   └── api.go          # Container registry API
│   │   │   │
│   │   │   └── layers/
│   │   │       ├── window.go       # Layer inspection view
│   │   │       └── api.go          # Layer fetch logic
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

Data hierarchy:
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

POI workflow:
```
Discover email in GitHub commit
  -> Create new POI or add to existing
     -> Search email for social profiles
        -> Add social media accounts
           -> Discover personal links from profiles
              -> Add notes during investigation
```

---

## Phase 1: Foundation Setup

### 1.1 Project Initialization
- [ ] Initialize Go module for yolosint workspace
- [ ] Add TUIOS dependency
- [ ] Add SQLite dependency (modernc.org/sqlite)
- [ ] Create `internal/config/env.go` to load .env secrets
- [ ] Verify all dependencies compile

### 1.2 Style Bible Creation
- [ ] Create `internal/style/stylebible.go` with TUIOS-compliant theme constants
- [ ] Define color palette constants (no magic numbers)
- [ ] Define border style constants
- [ ] Define font/typography constants
- [ ] Document naming conventions for UI elements

---

## Phase 2: Core Infrastructure

**WARNING: Extract ONLY pure data logic. Strip all UI/UX code.**

### 2.1 Project Module
- [ ] Create `internal/project/models.go` - Project struct
- [ ] Create `internal/project/db.go` - Project CRUD operations

### 2.2 GitHub Module (child of Project, peer of Domain)
Port from gitsome-ng, strip UI contamination:
- [ ] Create `internal/project/github/models.go` - Commit, user, repo structs
- [ ] Create `internal/project/github/db.go` - GitHub data storage
- [ ] Create `internal/project/github/commits/api.go` - Commits API
- [ ] Create `internal/project/github/commits/db.go` - Commits queries
- [ ] Create `internal/project/github/profile/api.go` - User API
- [ ] Create `internal/project/github/profile/db.go` - Profile queries
- [ ] Create `internal/project/github/repos/api.go` - Repos API
- [ ] Create `internal/project/github/repos/db.go` - Repos queries

### 2.3 Domain Module (child of Project, peer of GitHub)
- [ ] Create `internal/project/domain/models.go` - Domain struct
- [ ] Create `internal/project/domain/db.go` - Domain CRUD operations

### 2.4 Wayback Module (child of Domain)
Port from gitsome-ng, strip UI contamination:
- [ ] Create `internal/project/domain/wayback/models.go` - CDX record structs
- [ ] Create `internal/project/domain/wayback/api.go` - CDX API client
- [ ] Create `internal/project/domain/wayback/db.go` - Wayback storage

### 2.5 Subdomain Module (child of Domain)
Port from gitsome-ng, strip UI contamination:
- [ ] Create `internal/project/domain/subdomain/models.go` - Subdomain structs
- [ ] Create `internal/project/domain/subdomain/db.go` - Subdomain storage
- [ ] Create `internal/project/domain/subdomain/virustotal/api.go` - VT client
- [ ] Create `internal/project/domain/subdomain/crtsh/api.go` - crt.sh client

### 2.6 POI Module (child of Project, peer of Domain)
- [ ] Create `internal/project/poi/models.go` - Person of Interest struct
- [ ] Create `internal/project/poi/db.go` - POI CRUD operations

### 2.7 POI Email Submodule (child of POI)
- [ ] Create `internal/project/poi/email/models.go` - Email struct
- [ ] Create `internal/project/poi/email/db.go` - Email storage

### 2.8 POI Social Submodule (child of POI)
- [ ] Create `internal/project/poi/social/models.go` - Platform, handle, URL
- [ ] Create `internal/project/poi/social/db.go` - Social profile storage

### 2.9 POI Links Submodule (child of POI)
- [ ] Create `internal/project/poi/links/models.go` - URL, title, description
- [ ] Create `internal/project/poi/links/db.go` - Links storage

### 2.10 POI Notes Submodule (child of POI)
- [ ] Create `internal/project/poi/notes/models.go` - Freeform notes, timestamps
- [ ] Create `internal/project/poi/notes/db.go` - Notes storage

### 2.11 Docker Module (child of Project, peer of Domain)
Port from gitsome-ng, strip UI contamination:
- [ ] Create `internal/project/docker/models.go` - Image, tag structs
- [ ] Create `internal/project/docker/db.go` - Docker data storage
- [ ] Create `internal/project/docker/hub/api.go` - Docker Hub API
- [ ] Create `internal/project/docker/registry/api.go` - Registry API
- [ ] Create `internal/project/docker/layers/api.go` - Layer fetch

### 2.12 Report Module (child of Project, peer of Domain)
- [ ] Create `internal/project/report/export.go` - Report generation
- [ ] Report naming: PROJECT-Reportname-DDMMYYYY-HMS

---

## Phase 3: TUIOS Integration

### 3.1 Application Shell
- [ ] Create `internal/app/app.go` - TUIOS app shell, window manager
- [ ] Create `internal/app/navigation.go` - Cross-feature navigation
- [ ] Update `cmd/yolosint/main.go` with TUIOS bootstrap
- [ ] Configure TUIOS options (theme, borders, workspaces)
- [ ] Implement graceful shutdown and cleanup

### 3.2 Window Architecture
- [ ] Study TUIOS window management API (use Deepwiki)
- [ ] Define window interface for all feature windows
- [ ] Implement window lifecycle hooks

---

## Phase 4: Feature Windows

Build ALL windows fresh using TUIOS.

### 4.1 Project Window
- [ ] `internal/project/window.go` - Project selector

### 4.2 GitHub Windows (peer of Domain)
- [ ] `internal/project/github/window.go` - GitHub main
- [ ] `internal/project/github/commits/window.go` - Commits view
- [ ] `internal/project/github/profile/window.go` - User profile
- [ ] `internal/project/github/repos/window.go` - Repos list

### 4.3 Domain Window (peer of GitHub)
- [ ] `internal/project/domain/window.go` - Domain management

### 4.4 Wayback Window (child of Domain)
- [ ] `internal/project/domain/wayback/window.go` - Wayback browser

### 4.5 Subdomain Window (child of Domain)
- [ ] `internal/project/domain/subdomain/window.go` - Subdomonster

### 4.6 POI Windows (peer of Domain)
- [ ] `internal/project/poi/window.go` - POI management (list, create, link)
- [ ] `internal/project/poi/email/window.go` - Email list
- [ ] `internal/project/poi/social/window.go` - Social profiles
- [ ] `internal/project/poi/links/window.go` - Personal links
- [ ] `internal/project/poi/notes/window.go` - Investigation notes

### 4.7 Docker Windows (peer of Domain)
- [ ] `internal/project/docker/window.go` - LayerSlayer main
- [ ] `internal/project/docker/layers/window.go` - Layer inspection

### 4.8 Report Window (peer of Domain)
- [ ] `internal/project/report/window.go` - Report export

---

## Phase 5: Polish and Testing

### 5.1 Keybindings
- [ ] Document all keybindings
- [ ] Ensure TUIOS-consistent navigation
- [ ] Implement help overlay

### 5.2 Testing
- [ ] Test with fixtures in `docs/testing/fixtures/`
- [ ] Test against live APIs
- [ ] Verify all windows render correctly
- [ ] Test project creation/loading flow

### 5.3 Error Handling
- [ ] All errors fail explicitly (no fallbacks)
- [ ] User-friendly error messages
- [ ] No linter warnings

---

## Constraints Checklist

Before marking any task complete:
- [ ] No UI/UX code imported from gitsome-ng
- [ ] No file exceeds 100 lines (refactor into submodule)
- [ ] No magic numbers (use named constants)
- [ ] No linter errors
- [ ] No unused variables
- [ ] No emojis anywhere
- [ ] Tested in collaboration with user
- [ ] Code is modular and reusable

---

## Documentation Reference

| Source | Path | Purpose |
|--------|------|---------|
| TUIOS Docs | `docs/includes/tuios/docs/` | Windowing documentation |
| TUIOS Library | `docs/includes/tuios/pkg/tuios/` | Library API |
| gitsome-ng API | `docs/includes/gitsome-ng/internal/api/` | API reference (DATA ONLY) |
| gitsome-ng DB | `docs/includes/gitsome-ng/internal/db/` | DB reference (DATA ONLY) |
| gitsome-ng Models | `docs/includes/gitsome-ng/internal/models/` | Model reference (DATA ONLY) |
| Test Fixtures | `docs/testing/fixtures/` | Sample data |

---

## Environment Variables

Required secrets in `.env`:
- `GITHUB_TOKEN` - GitHub API authentication
- `VIRUS_TOTAL_TOKEN` - VirusTotal API key
- `DOCKERHUB_USER` - Docker Hub username
- `DOCKERHUB_TOKEN` - Docker Hub access token

---

## Notes

- **UI Code from gitsome-ng**: FORBIDDEN. Do not port.
- **Charmbracelet imports**: Strip them. They indicate UI contamination.
- **CLI Arguments**: Not implementing. Focus on TUI.
- **Fallbacks**: Never implement. Fail fast.
- **Documentation**: Always check context7/Deepwiki before writing code.


# PROJECT KNOWLEDGE BASE

**Generated:** 2026-03-01T03:26:39Z
**Commit:** N/A (uncommitted)
**Branch:** master

## OVERVIEW
Go-based CLI tool (`skillmgr`) for registering and managing local skill repositories. Built with `spf13/cobra`.

## STRUCTURE
```
.
├── cmd/          # CLI command definitions (Cobra)
├── internal/     # Core business logic (encapsulated)
│   ├── config/   # Config file management (~/.skillmgr/config.json)
│   ├── link/     # Symlinking logic
│   ├── move/     # Move/migrate logic
│   ├── repo/     # Repository registration
│   └── status/   # Status reporting
└── main.go       # Entry point
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| CLI Commands | `cmd/` | `root.go` is the entry point for all commands |
| Config Logic | `internal/config/` | Handles JSON persistence in home dir |
| Core Logic | `internal/` | Packages are highly granular (one per concern) |

## CODE MAP
| Symbol | Type | Location | Role |
|--------|------|----------|------|
| `rootCmd` | Variable | `cmd/root.go` | Root Cobra command |
| `Config` | Struct | `internal/config/config.go` | Top-level config structure |
| `Load` | Function | `internal/config/config.go` | Loads/initializes config |

## CONVENTIONS
- **CLI Framework**: Always use `spf13/cobra`.
- **Logic Isolation**: All core logic must reside in `internal/`.
- **Encapsulation**: Prefer deep fragmentation in `internal/` (one package per feature).
- **Persistence**: Use `~/.skillmgr/config.json` for all state.

## ANTI-PATTERNS (THIS PROJECT)
- **No Tests**: Currently lacks automated tests (Significant deviation).
- **Public API**: Do not expose packages in `pkg/`; keep everything in `internal/`.

## COMMANDS
```bash
# Build
go build -o skillmgr main.go

# Run
./skillmgr
```

## NOTES
- **Home Directory**: Depends on `os.UserHomeDir()` for config location.
- **Atomicity**: `config.Save` uses temp file + rename for atomic writes.

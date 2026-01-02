# SwiftStack

SwiftStack — a high-speed project scaffolder

SwiftStack is a command-line tool, implemented in Go, that assembles new projects by stitching together pre-built "slices" (packaged project pieces). Rather than installing dependencies from package registries for every fresh project, SwiftStack composes projects from compressed, versioned slice artifacts (.tar.zst) that live in a local cache or a remote registry. Every slice includes a SHA-256 hash so consumers can verify integrity.

This README documents the repository as implemented in source (commands, engine behaviour, and internals).

Highlights

- Create new projects from pre-built slices (bases + addons) quickly.
- Slices are compressed as .tar.zst and include SHA-256 hashes.
- Local cache for slices in the OS user cache directory.
- CLI commands: `create`, `build`, `sync`, `ui` (interactive wizard).
- Merges slice `package.json` into base project and runs an npm lockfile update for consistent installs.
- Written in Go with minimal runtime dependencies.

Table of contents

- [Quickstart](#quickstart)
- [Install](#install)
- [CLI reference](#cli-reference)
- [CLI examples (with full session outputs)](#cli-examples-with-full-session-outputs)
- [How SwiftStack works (internals)](#how-swiftstack-works-internals)
- [Slice format & registry](#slice-format--registry)
- [Development](#development)
- [Troubleshooting & tips](#troubleshooting--tips)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

Quickstart

1. Install SwiftStack (see Install).
2. Create a new project using a base slice and optional addons:

```bash
swiftstack create --name my-app --base next-base --addons tailwind-ui,auth
```

The `create` command will:
- Resolve slice aliases via the local/remote manifest,
- Ensure slices are cached locally (downloads if necessary),
- Verify slice integrity using SHA-256,
- Extract and merge slices into the target directory,
- Merge slice `package.json` files into the base, and
- Run an npm lockfile update to finalize dependency references.

Install

From prebuilt releases
- Download the appropriate binary for your OS from the Releases page and put it on your PATH.

From source (recommended for contributors)
- Requires Go (see `go.mod` for minimum required version). Build/install:

```bash
# Install the CLI into $GOBIN (or $GOPATH/bin)
go install github.com/004Ongoro/swiftstack/cmd/swiftstack@latest

# or clone and build locally
git clone https://github.com/004Ongoro/swiftstack.git
cd swiftstack
go build -o swiftstack ./cmd/swiftstack
```

CLI reference

- `swiftstack create --name <project> --base <base-alias> [--addons <a,b>]`
  - Create a new project from a base slice and zero or more addon slices.
  - Flags:
    - `--name`, `-n` (required) — project name
    - `--base`, `-b` (required) — base slice alias (e.g., `next-base`)
    - `--addons`, `-a` — comma-separated addon aliases

- `swiftstack build [source_dir] [output_file.tar.zst]`
  - Pack a directory into a `.tar.zst` slice and print its SHA-256. Use this when producing slices to publish to a registry/manifest.

- `swiftstack sync`
  - Update the local registry/manifest from the remote source (used to resolve aliases to slice URLs).

- `swiftstack ui` (or `swiftstack wizard`)
  - Start the interactive terminal wizard (TUI) to assemble a project using a guided flow.

CLI examples (with full session outputs)

Below are example sessions that show typical output you can expect from each command. These are realistic simulated outputs reflecting messages emitted by code in this repository.

1) Create a project (base + addons)

```bash
$ swiftstack create --name example-app --base next-base --addons tailwind-ui,auth
🚀 Starting SwiftStack assembly for 'example-app'...

Resolving slice 'next-base' from registry...
Downloading next-base...                         # only shown if the slice is not cached
Verifying integrity of next-base...
Extracting next-base -> ./example-app

Processing addon: tailwind-ui
Downloading tailwind-ui...
Verifying integrity of tailwind-ui...
Extracting tailwind-ui -> ./example-app/.swiftstack_temp
Merging package.json from tailwind-ui into base package.json
Moving files from addon into project (with backup)

Processing addon: auth
Downloading auth...
Verifying integrity of auth...
Extracting auth -> ./example-app/.swiftstack_temp
(Merging package.json from auth)
Moving files from addon into project (with backup)

Finalizing project structure...
Running npm lockfile update (RunNpmLockUpdate) to ensure consistent lock entries

✨ Successfully assembled 'example-app' in record time!
```

Notes:
- If a required slice alias is not found in the manifest, the command exits with an error like:
  - `engine: slice <alias> not found in registry`
- If integrity verification fails, SwiftStack deletes the corrupted cache entry and returns an error:
  - `security alert: verification failed for <slice>`

2) Build a slice artifact (pack a directory)

```bash
$ swiftstack build ./component ./component@v1.0.0.tar.zst
Building slice from ./component...
Successfully created slice!
Location: ./component@v1.0.0.tar.zst
SHA-256:  9f2c7d4b6a5e2b... (full 64-char hex)

# After building, add the SHA to your registry/manifest so consumers can verify:
# (manual step) add slice metadata (id, url, version, hash) to registry.json / manifest
```

3) Sync registry (update manifest)

```bash
$ swiftstack sync
Fetching remote manifest from registry...
Downloaded manifest -> updated local manifest cache
Registry updated successfully!
```

4) Interactive wizard (TUI)

The `ui` command launches an interactive text UI (uses Bubble Tea). A sample non-interactive transcript:

```text
$ swiftstack ui
[ TUI opens in your terminal ]

Welcome to SwiftStack — Create a project in minutes
> Project name: [ my-app               ]
> Choose base slice:
  • next-base
  • react-base
  • minimal-base

[Select next-base]

> Choose addons:
  [ ] tailwind-ui
  [ ] auth
  [ ] analytics

[Select tailwind-ui, auth]

[Press Enter to assemble]

[Status]
Resolving slices...
Downloading next-base (if needed)
Verifying integrity...
Extracting...
Merging addons...
Finalizing...

[Success] Project 'my-app' created at ./my-app
```

How SwiftStack works (internals)

- Slices
  - Packaged as `.tar.zst` files created by the `builder` package.
  - Each slice should have an associated SHA-256 hash in the registry manifest so consumers can verify integrity.

- Cache
  - Local cache directory is based on the OS user cache dir (`os.UserCacheDir()`), under `swiftstack`.
  - Slice filename format: `<id>@<version>.tar.zst` (e.g., `next-base@1.0.0.tar.zst`).

- Project generation (`internal/engine`)
  - `GenerateProject` orchestrates the flow:
    1. Resolve slice metadata (URL & hash) from the manifest (`cache.LoadManifest()`).
    2. Ensure the slice is present in the local cache (download if missing).
    3. Verify file SHA-256 integrity via `utils.VerifyFileHash`.
    4. Extract base slice into the target project directory.
    5. For each addon:
       - Extract the addon into a `.swiftstack_temp` directory inside the project.
       - If the addon contains `package.json`, merge it into the base project's `package.json` (via `MergePackageJSON`).
       - Move addon files into the project using `utils.MoveWithBackup`.
    6. Run `RunNpmLockUpdate(fullPath)` to update the lockfile (ensures dependency references are coherent).
  - On errors during assembly, the engine attempts to clean up the partially created project directory.

- Builder (`internal/builder`)
  - Walks a source directory, creates a tar archive and compresses it using zstd (`klauspost/compress/zstd`), producing `.tar.zst` files.
  - The `build` CLI prints the SHA-256 so the artifact author can add it to the registry manifest.

Slice format & registry

- Manifest model (`internal/models/manifest.go`):
  - `SliceMetadata`:
    - `id`, `title`, `description`, `url`, `version`, `hash` (SHA-256)
  - `RemoteManifest`:
    - `bases` (array), `addons` (array)

- Registry operations
  - `ResolveAlias(alias)` searches the manifest for an alias and returns its URL (suggests running `swiftstack sync` if not found).
  - `sync` fetches the remote manifest and replaces/updates the local manifest cache.

Development

- Run tests:

```bash
go test ./...
```

- Formatting & linting:

```bash
gofmt -w .
# optional: golangci-lint run
```

- Build locally:

```bash
go build -o swiftstack ./cmd/swiftstack
```

- Install:

```bash
go install ./cmd/swiftstack
```

Repository notes

- CLI is built with Cobra (see `cmd/swiftstack/*.go`).
- Interactive UI uses Charmbracelet Bubble Tea (`internal/ui`).
- Archives use zstd compression (via `github.com/klauspost/compress/zstd`).
- The engine merges `package.json` files and relies on Node tooling (npm lockfile update) to finalize dependency references in the produced project.
- There's a PowerShell release helper `release.ps1` to automate tagging and pushing to origin.

Troubleshooting & tips

- If alias resolution fails:
  - Run `swiftstack sync` to refresh the manifest.
- If integrity verification fails:
  - The cached slice will be removed; re-run the command to redownload.
- If addon merging clobbers files, inspect the temporary extraction directory `.swiftstack_temp` created during assembly.
- The engine attempts to back up files when moving to avoid accidental data loss — verify backups in case of unexpected changes.
- Because SwiftStack modifies `package.json` and the lockfile, ensure you review generated package metadata before shipping.

Contributing

- Fork the repository, create a branch, and open a pull request.
- Run tests and ensure `gofmt` and `go vet` pass.
- Add unit tests for behavioural changes (especially for `MergePackageJSON`, `MoveWithBackup`, and `RunNpmLockUpdate`).
- Describe behavioral changes clearly in PR descriptions.

License

See LICENSE in the repository root. If none exists, add an appropriate open-source license (MIT / Apache-2.0 are common choices).

Acknowledgements

- Built with Go, Cobra for CLI, Bubble Tea for TUI, and klauspost zstd for compression.
- Designed to make project scaffolding fast by re-using pre-built slice artifacts rather than re-installing dependencies on every new project.

Appendix: Example manifest (schema)

```json
{
  "bases": [
    {
      "id": "next-base",
      "title": "Next.js Base",
      "description": "Opinionated Next.js base project",
      "url": "https://cdn.example.com/slices/next-base@1.0.0.tar.zst",
      "version": "1.0.0",
      "hash": "9f2c7d4b6a5e2b... (64 hex chars)"
    }
  ],
  "addons": [
    {
      "id": "tailwind-ui",
      "title": "Tailwind UI Setup",
      "description": "Tailwind + config + sample components",
      "url": "https://cdn.example.com/slices/tailwind-ui@0.4.2.tar.zst",
      "version": "0.4.2",
      "hash": "3a4b5c6d7e8f... (64 hex chars)"
    }
  ]
}
```

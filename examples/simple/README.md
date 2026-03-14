# Simple Example

This example demonstrates using `qtinstaller` to generate a Qt Installer Framework
package structure. It uses a dummy executable and placeholder assets.

## Files

- `.env` — Configuration for the installer
- `myapp` — Placeholder executable (dummy file)
- `logo.png` — Placeholder logo (1x1 pixel PNG)
- `licence.txt` — Example licence file

## Usage

```bash
cd examples/simple

# Build qtinstaller (from project root)
go build -o ../../qtinstaller-cli ../../main.go

# Generate a fresh .env template (optional — one is already provided)
../../qtinstaller-cli init

# Run the installer packaging (will fail at deploy tool step since
# myapp is not a real Qt binary, but the directory structure and
# XML files will be generated successfully up to that point)
../../qtinstaller-cli
```

## What it tests

1. `.env` parsing with `godotenv`
2. Config validation (`AssertFilesExist`)
3. Directory structure creation (`config/`, `packages/<pkg>/data/`, `packages/<pkg>/meta/`)
4. Template rendering (`config.xml`, `package.xml`, `installscript.qs`)

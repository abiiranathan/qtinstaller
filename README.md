# qtinstaller

[![Go](https://img.shields.io/badge/--00ADD8?logo=go&logoColor=ffffff)](https://golang.org/)

Cross-platform wrapper around the [Qt Installer Framework](https://doc.qt.io/qtinstallerframework/) for easily generating offline installers for your [Qt](https://doc.qt.io/qt.html) applications on **Windows** and **Linux**.

On Windows it uses [**windeployqt**](https://doc.qt.io/qt-6/windows-deployment.html) to gather DLLs. On Linux it includes a built-in deployer that copies Qt shared libraries, plugins, and symlinks using `ldd` and `qtpaths6` — no external deploy tool required. [**binarycreator**](https://doc.qt.io/qtinstallerframework/ifw-tools.html#binarycreator) (from the Qt Installer Framework) is required on both platforms.

> NB: This assumes Qt6. It has not been tested on Qt5 or lower versions.

## Prerequisites

|                   | Windows                  | Linux                        |
| ----------------- | ------------------------ | ---------------------------- |
| **Qt 6**          | Required                 | Required                     |
| **windeployqt**   | Required (ships with Qt) | N/A                          |
| **binarycreator** | Required (Qt IFW)        | Required (Qt IFW)            |
| **linuxdeployqt** | N/A                      | Optional (built-in fallback) |

## Installation

```bash
go install github.com/abiiranathan/qtinstaller@latest
```

Or build from source:

```bash
git clone https://github.com/abiiranathan/qtinstaller.git
cd qtinstaller
go build -o qtinstaller .
```

## Usage

### 1. Generate a config file

Run `init` in your project directory to generate a template `.env` file:

```bash
# Windows
qtinstaller.exe init

# Linux
qtinstaller init
```

This creates a `.env` file with platform-appropriate defaults.

### 2. Edit the `.env` file

Fill in all required variables:

```env
# Application display name
DISPLAY_NAME=MyApp

# Company or developer name
PUBLISHER=Acme Software

# Short description of the application
DESCRIPTION=A Qt desktop application

# Release date (YYYY-MM-DD)
RELEASE_DATE=2026-03-14

# Semantic version
VERSION=1.0.0

# Globally unique package identifier
PACKAGE_NAME=com.acme.myapp

# Path to the built Qt executable
EXECUTABLE=build/release/myapp        # Linux
# EXECUTABLE=build\release\myapp.exe  # Windows

# Output installer filename
INSTALLER_NAME=MyApp-Installer        # Linux
# INSTALLER_NAME=MyApp-Installer.exe  # Windows

# License name shown in the installer wizard
LICENSE_NAME=MIT License Agreement

# Path to the license file
LICENSE_FILE=licence.txt

# Logo in PNG format (used in the installer wizard)
LOGO=logo.png

# Custom installer icon (required on Windows, optional on Linux)
# Windows: .ico format   Linux: .png format (or leave empty)
INSTALLER_APPLICATION_ICON=favicon.ico  # Windows
# INSTALLER_APPLICATION_ICON=           # Linux (optional)
```

### 3. Build the installer

```bash
qtinstaller
```

This will:

1. **Parse** the `.env` configuration
2. **Validate** that all referenced files exist
3. **Create** the Qt Installer Framework directory structure:
   ```
   config/
     config.xml
     logo.png
     installerLogo.png
   packages/<PACKAGE_NAME>/
     data/
       <executable>
       <Qt libraries and plugins>   (Linux)
       <Qt DLLs>                    (Windows)
     meta/
       package.xml
       licence.txt
       installscript.qs
   ```
4. **Deploy Qt dependencies**:
   - **Windows**: runs `windeployqt` to copy required DLLs
   - **Linux**: uses the built-in deployer (copies Qt `.so` files, plugins, symlinks, and `qt.conf`)
5. **Run `binarycreator`** to produce the final offline installer

The generated installer will be in the current directory.

## How the Linux deployer works

When no external deploy tool (`linuxdeployqt`, `cqtdeployer`) is found, the built-in deployer:

1. Runs `ldd` on the executable to discover linked Qt libraries
2. Copies Qt6, ICU, XCB, EGL, OpenGL, and xkbcommon libraries with proper symlinks
3. Uses `qtpaths6` (or falls back to common paths) to find the Qt plugins directory
4. Copies platform plugins (xcb, wayland, eglfs), theme plugins, and image format plugins
5. Writes a `qt.conf` so the application finds bundled plugins at runtime

## Installer features

The generated installer includes:

- **Windows**: Start Menu shortcut + Desktop shortcut (via `installscript.qs`)
- **Linux**: `.desktop` file creation (via Qt IFW's `CreateDesktopEntry` operation)
- Offline-only installation (no network required)
- License agreement page
- Configurable target directory (`@ApplicationsDir@/<DisplayName>`)

## Example

See the [examples/simple](examples/simple/) directory for a complete working example using a Qt6 Hello World application.

```bash
cd examples/simple
cat .env               # review the configuration
../../qtinstaller      # build the installer
./HelloQt-Installer    # run the generated installer
```

## Environment variables

| Variable                     | Required     | Description                                              |
| ---------------------------- | ------------ | -------------------------------------------------------- |
| `DISPLAY_NAME`               | Yes          | Application name shown in the installer                  |
| `PUBLISHER`                  | Yes          | Company or developer name                                |
| `DESCRIPTION`                | Yes          | Short application description                            |
| `PACKAGE_NAME`               | Yes          | Unique identifier (e.g. `com.company.app`)               |
| `EXECUTABLE`                 | Yes          | Path to the built Qt binary                              |
| `LICENSE_NAME`               | Yes          | License name for the wizard                              |
| `LICENSE_FILE`               | Yes          | Path to the license text file                            |
| `LOGO`                       | Yes          | PNG logo for the installer wizard                        |
| `INSTALLER_APPLICATION_ICON` | Windows only | `.ico` file for the installer (optional on Linux)        |
| `VERSION`                    | No           | Version string (default: `0.1.0-1`)                      |
| `RELEASE_DATE`               | No           | Release date (default: today)                            |
| `INSTALLER_NAME`             | No           | Output filename (default: `Installer.exe` / `Installer`) |

## License

See [LICENCE](LICENCE).
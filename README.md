# qtinstaller

[![Go](https://img.shields.io/badge/--00ADD8?logo=go&logoColor=ffffff)](https://golang.org/)

Wrapper around [**windeployqt**](https://doc.qt.io/qt-6/windows-deployment.html) and [**binarycreator**](https://doc.qt.io/qtinstallerframework/ifw-tools.html#binarycreator) for easily generating windows installers for your [Qt](https://doc.qt.io/qt.html) applications.

*The two binaries must be installed and already in PATH.* These are normally installed together with Qt.

> NB: This assumes Qt6. This has not been tested on Qt5 or lower versions. Feel free to try it out and share your feedback.

## Installation
```bash
go install github.com/abiiranathan/qtinstaller@latest
```

### Usage

- Generate .env file in the same directory.
   This will store installer configuration.
   ```bash
   qtinstaller.exe init
   ```
- Edit the .env file and set all the required variables then run ```qtinstaller.exe``` to generate the installer.

> Note that the *.exe* extension can be ommitted.

By default, this generates an offline-only installer and sets up a start menu shortcut and desktop icon for your application using an install script.

### Example env file

```txt

# Application display name
DISPLAY_NAME=HMIS

# Company name
PUBLISHER=Acme Software Technologies

# Description of what the software does.
DESCRIPTION=Health information management system

# Release date
RELEASE_DATE=2023-01-01

# Version number of the software
VERSION=0.1.0-1

# Package name for the software
PACKAGE_NAME=com.acme.hmis

# The path to the release Qt application
EXECUTABLE=C:\Qt\Projects\Releases\HMIS.exe

# Output name for the installer
INSTALLER_NAME=Installer.exe

# Name of the license
LICENSE_NAME=MIT LICENCE

# Path to the license file
LICENSE_FILE=license.txt

# Filename for a logo in PNG format used as QWizard::LogoPixmap.
LOGO=logo.png

# Filename for a custom installer icon.
INSTALLER_APPLICATION_ICON=favicon.ico

```
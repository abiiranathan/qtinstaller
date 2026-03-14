package qtinstaller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const linuxDeployQtURL = "https://github.com/probonopd/linuxdeployqt/releases/download/continuous/linuxdeployqt-continuous-x86_64.AppImage"

// toolCacheDir returns the directory used to cache downloaded tools.
// Defaults to ~/.local/share/qtinstaller/bin.
func toolCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	dir := filepath.Join(home, ".local", "share", "qtinstaller", "bin")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("could not create cache directory: %w", err)
	}
	return dir, nil
}

// downloadFile downloads a file from url and saves it to destPath, then
// makes it executable.
func downloadFile(url, destPath string) error {
	log.Printf("Downloading %s ...\n", url)

	resp, err := http.Get(url) //#nosec G107 -- URL is a compile-time constant, not user input
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d for %s", resp.StatusCode, url)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("failed writing to %s: %w", destPath, err)
	}

	if err := out.Close(); err != nil {
		return err
	}

	// Make the file executable (necessary for AppImage on Linux)
	return os.Chmod(destPath, 0755)
}

// downloadLinuxDeployQt downloads linuxdeployqt AppImage into the tool
// cache directory and returns the path to the binary.
func downloadLinuxDeployQt() (string, error) {
	cacheDir, err := toolCacheDir()
	if err != nil {
		return "", err
	}

	dest := filepath.Join(cacheDir, "linuxdeployqt")

	// If already cached, reuse it.
	if info, err := os.Stat(dest); err == nil && info.Mode().IsRegular() {
		log.Printf("Using cached linuxdeployqt at %s\n", dest)
		return dest, nil
	}

	if err := downloadFile(linuxDeployQtURL, dest); err != nil {
		return "", err
	}

	log.Printf("linuxdeployqt downloaded to %s\n", dest)
	return dest, nil
}

// findWinDeployQt searches common Qt installation paths for windeployqt.exe.
func findWinDeployQt() (string, error) {
	searchDirs := []string{
		`C:\Qt`,
		`C:\Qt6`,
		`C:\Qt5`,
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Qt"),
		filepath.Join(os.Getenv("ProgramFiles"), "Qt"),
	}

	for _, root := range searchDirs {
		if root == "" {
			continue
		}
		matches, _ := filepath.Glob(filepath.Join(root, "*", "*", "bin", "windeployqt.exe"))
		if len(matches) > 0 {
			return matches[len(matches)-1], nil // use latest version found
		}
		// Also check directly under bin/
		matches, _ = filepath.Glob(filepath.Join(root, "bin", "windeployqt.exe"))
		if len(matches) > 0 {
			return matches[0], nil
		}
	}

	return "", fmt.Errorf("windeployqt.exe not found in common Qt installation paths (%v). "+
		"Ensure Qt is installed and windeployqt.exe is in your PATH", searchDirs)
}

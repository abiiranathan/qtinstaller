package qtinstaller

import (
	"fmt"
	"os"
	"path/filepath"
)

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

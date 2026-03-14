package qtinstaller

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DeployQtLibs copies the Qt shared libraries needed by the given binary
// into destDir. This is a native Go fallback for when linuxdeployqt is
// unavailable or incompatible with the host system.
func DeployQtLibs(binary, destDir string) error {
	output, err := exec.Command("ldd", binary).CombinedOutput()
	if err != nil {
		return fmt.Errorf("ldd failed on %s: %w\n%s", binary, err, output)
	}

	var copied []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)

		// ldd output format: "libFoo.so.6 => /usr/lib/libFoo.so.6 (0x...)"
		parts := strings.SplitN(line, " => ", 2)
		if len(parts) != 2 {
			continue
		}

		libPath := strings.TrimSpace(strings.SplitN(parts[1], " (", 2)[0])
		if libPath == "" || libPath == "not found" {
			continue
		}

		libName := strings.TrimSpace(parts[0])

		// Copy Qt libraries and their common dependencies
		if isQtRelatedLib(libName) {
			if err := copyLibWithSymlinks(libPath, destDir); err != nil {
				return fmt.Errorf("failed to copy %s: %w", libPath, err)
			}
			copied = append(copied, libName)
		}
	}

	if len(copied) == 0 {
		return fmt.Errorf("no Qt libraries found for %s — is it a Qt application?", binary)
	}

	log.Printf("Copied %d Qt libraries to %s\n", len(copied), destDir)

	// Also deploy Qt plugins (platforms, etc.)
	return deployQtPlugins(destDir)
}

// isQtRelatedLib returns true if the library name is Qt-related or a
// common dependency that should be bundled.
func isQtRelatedLib(name string) bool {
	prefixes := []string{
		"libQt6", "libQt5",
		"libicu",    // ICU libs needed by Qt
		"libxcb",    // X11/XCB libs
		"libxkb",    // keyboard libs
		"libEGL",    // graphics
		"libGLX",    // graphics
		"libOpenGL", // graphics
	}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

// copyLibWithSymlinks copies a library and resolves symlinks so both the
// symlink name and the real file end up in destDir.
func copyLibWithSymlinks(libPath, destDir string) error {
	// Read the real file
	realPath, err := filepath.EvalSymlinks(libPath)
	if err != nil {
		realPath = libPath
	}

	data, err := os.ReadFile(realPath)
	if err != nil {
		return err
	}

	// Copy the real file
	realBase := filepath.Base(realPath)
	destReal := filepath.Join(destDir, realBase)
	if err := os.WriteFile(destReal, data, 0755); err != nil {
		return err
	}

	// If the original was a symlink, create the symlink too
	linkBase := filepath.Base(libPath)
	if linkBase != realBase {
		destLink := filepath.Join(destDir, linkBase)
		os.Remove(destLink) // remove if exists
		if err := os.Symlink(realBase, destLink); err != nil {
			return err
		}
	}

	return nil
}

// deployQtPlugins copies essential Qt platform plugins to destDir/plugins/.
func deployQtPlugins(destDir string) error {
	// Find Qt plugin directory
	qtPluginDir := findQtPluginDir()
	if qtPluginDir == "" {
		log.Println("Warning: Qt plugins directory not found, skipping plugin deployment")
		return nil
	}

	// Essential plugin subdirectories
	pluginDirs := []string{"platforms", "platformthemes", "imageformats", "wayland-shell-integration"}

	for _, subdir := range pluginDirs {
		srcDir := filepath.Join(qtPluginDir, subdir)
		if _, err := os.Stat(srcDir); err != nil {
			continue
		}

		dstDir := filepath.Join(destDir, "plugins", subdir)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return err
		}

		entries, err := os.ReadDir(srcDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".so") {
				continue
			}
			src := filepath.Join(srcDir, entry.Name())
			data, err := os.ReadFile(src)
			if err != nil {
				continue
			}
			if err := os.WriteFile(filepath.Join(dstDir, entry.Name()), data, 0755); err != nil {
				return err
			}
		}

		log.Printf("Deployed %s plugins from %s\n", subdir, srcDir)
	}

	// Write qt.conf so the app finds the plugins
	qtConf := filepath.Join(destDir, "qt.conf")
	return os.WriteFile(qtConf, []byte("[Paths]\nPlugins = plugins\n"), 0644)
}

// findQtPluginDir locates the Qt plugins directory on the system.
func findQtPluginDir() string {
	// Try qtpaths first (Qt6)
	if out, err := exec.Command("qtpaths6", "--plugin-dir").CombinedOutput(); err == nil {
		if dir := strings.TrimSpace(string(out)); dir != "" {
			return dir
		}
	}

	// Try qtpaths (Qt5/6)
	if out, err := exec.Command("qtpaths", "--plugin-dir").CombinedOutput(); err == nil {
		if dir := strings.TrimSpace(string(out)); dir != "" {
			return dir
		}
	}

	// Common paths
	candidates := []string{
		"/usr/lib/qt6/plugins",
		"/usr/lib/qt/plugins",
		"/usr/lib64/qt6/plugins",
		"/usr/lib/x86_64-linux-gnu/qt6/plugins",
	}
	for _, d := range candidates {
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			return d
		}
	}
	return ""
}

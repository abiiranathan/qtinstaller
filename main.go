package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"qtinstaller/qtinstaller"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	log.SetPrefix("qtinstaller: ")
	log.SetFlags(log.Lshortfile)

	// If init argument, generate config file and exit.
	if len(os.Args) > 1 && os.Args[1] == "init" {
		err := qtinstaller.GenerateConfigFile()
		if err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}

	// Make sure windeployqt and binarycreator are installed and in PATH
	windeployqt, binarycreator, err := qtinstaller.GetQtBinaries()
	if err != nil {
		log.Fatalln(err)
	}

	// Create default configuration
	config := qtinstaller.NewConfig()

	// Assert that all expected files exist
	qtinstaller.AssertFilesExist(config)

	// Get installation directories
	dirs := qtinstaller.GetInstallerDirs(config)
	err = qtinstaller.CreateDirectoryStructure(dirs)
	if err != nil {
		log.Fatalln(err)
	}

	// write the installer files to disk.
	err = qtinstaller.WriteFiles(config, dirs)
	if err != nil {
		log.Fatalln(err)
	}

	// get path to target exe in data directory.
	targetExe := filepath.Join(dirs.Data, filepath.Base(config.Executable))

	// Run windeployqt on the executable to gather neccessary dlls.
	output, err := exec.Command(windeployqt, "--no-translations", "--dir", dirs.Data, targetExe).CombinedOutput()
	if err != nil {
		log.Fatalln(string(output))
	}

	// Run binary creator to generate the installer.
	output, err = exec.Command(
		binarycreator, "--offline-only", "-c", "config/config.xml", "-p", "packages",
		config.InstallerName).CombinedOutput()

	if err != nil {
		log.Fatalln(string(output))
	}

	wd, _ := os.Getwd()
	log.Printf("Installer created at path %q\n", filepath.Join(wd, config.InstallerName))
}

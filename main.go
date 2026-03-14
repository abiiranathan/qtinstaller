package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/abiiranathan/qtinstaller/qtinstaller"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

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

	// Make sure deploy tool and binarycreator are installed and in PATH
	deployTool, binarycreator, err := qtinstaller.GetQtBinaries()
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

	// Run deploy tool to gather necessary libraries.
	var output []byte
	if runtime.GOOS == "windows" {
		output, err = exec.Command(deployTool, "--no-translations", "--dir", dirs.Data, targetExe).CombinedOutput()
		if err != nil {
			log.Fatalln(string(output))
		}
	} else if deployTool != "" {
		output, err = exec.Command(deployTool, targetExe, "-no-translations", "-always-overwrite").CombinedOutput()
		if err != nil {
			log.Printf("External deploy tool failed: %s\n", string(output))
			log.Println("Falling back to built-in Qt library deployer...")
			if err := qtinstaller.DeployQtLibs(targetExe, dirs.Data); err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		// No external tool — use built-in deployer
		if err := qtinstaller.DeployQtLibs(targetExe, dirs.Data); err != nil {
			log.Fatalln(err)
		}
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

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	repoURL      = "https://github.com/AchoraSoft/kissc.at/tree/master/"
	denoInstallURL = "https://deno.land/install.sh"
	envContent   = `PORT=8000
				BASE_PATH=./snaps
				HOME_PATH=home
				ALLOWED_METHODS=GET,POST,DELETE
				VIEWS_BASE=./snaps
				COMPONENTS_DIR=components
				PUBLIC_PATH=public
				`
	)

var templateFiles = []string{
	"core/Controller.ts",
	"core/Router.ts",
	"core/types.ts",
	"core/Views.ts",
	"deno.json",
	"package.json",
	"server.ts",
	"snaps/home/get.ts",
	"snaps/layout.eta",
}

// Check if Deno is installed
func isDenoInstalled() bool {
	_, err := exec.LookPath("deno")
	return err == nil
}

// Install Deno using the official installer
func installDeno() error {
	fmt.Println("Deno not found. Installing Deno...")
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// PowerShell command for Windows
		psCmd := `irm https://deno.land/install.ps1 | iex`
		cmd = exec.Command("powershell", "-Command", psCmd)
	} else {
		// Shell command for Unix-like systems
		cmd = exec.Command("sh", "-c", "curl -fsSL https://deno.land/x/install/install.sh | sh")
	}
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Deno: %v", err)
	}

	// Add Deno to PATH if not already there
	if runtime.GOOS != "windows" {
		home, _ := os.UserHomeDir()
		denoBin := filepath.Join(home, ".deno", "bin")
		if !strings.Contains(os.Getenv("PATH"), denoBin) {
			fmt.Printf("\nPlease add Deno to your PATH:\n")
			fmt.Printf("export PATH=\"%s:$PATH\"\n", denoBin)
			fmt.Printf("Then run this command again.\n")
			os.Exit(1)
		}
	}
	
	fmt.Println("Deno installed successfully!")
	return nil
}

func downloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: %s", url, resp.Status)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func createProject(projectName string) error {
	// Check for Deno
	if !isDenoInstalled() {
		if err := installDeno(); err != nil {
			return err
		}
	}

	// Verify Deno is now available
	if !isDenoInstalled() {
		return fmt.Errorf("Deno installation failed - please install it manually from https://deno.land")
	}

	// Create project directory
	if err := os.Mkdir(projectName, 0755); err != nil {
		return err
	}

	// Create public directory
	if err := os.Mkdir(filepath.Join(projectName, "public"), 0755); err != nil {
		return err
	}

	// Download template files
	for _, file := range templateFiles {
		url := repoURL + file
		dest := filepath.Join(projectName, file)
		fmt.Printf("Downloading %s...\n", file)
		if err := downloadFile(url, dest); err != nil {
			return fmt.Errorf("error downloading %s: %v", file, err)
		}
	}

	// Create .env file
	envPath := filepath.Join(projectName, ".env")
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		return fmt.Errorf("error creating .env file: %v", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: my-cli <project-name>")
		os.Exit(1)
	}

	projectName := os.Args[1]
	fmt.Printf("Creating project %s...\n", projectName)

	if err := createProject(projectName); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(`
Project %s created successfully!

Next steps:
1. cd %s
2. deno run --allow-net --allow-read server.ts
`, projectName, projectName)
}
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

	"gopkg.in/yaml.v3"
)

const (
	repoURL       = "https://raw.githubusercontent.com/AchoraSoft/tofo/master/"
	templateRepo  = repoURL + "templates/"
	denoInstallURL = "https://deno.land/install.sh"
)

var coreFiles = []string{
    "core/Controller.ts",
    "core/Router.ts",
    "core/types.ts",
    "core/Views.ts",
    "deno.json",
    "package.json",
    "server.ts",
}
type RouteConfig struct {
	ViewName string `yaml:"viewName"`
}

type ProjectConfig struct {
	Project struct {
		Structure map[string]map[string]RouteConfig `yaml:"structure"`
		Env       map[string]string                 `yaml:"env"`
	} `yaml:"project"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	if os.Args[1] == "create" {
		handleCreateCommand()
	} else {
		// Legacy mode
		projectName := os.Args[1]
		if err := createProject(projectName, ""); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		printSuccess(projectName)
	}
}

func handleCreateCommand() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	projectName := os.Args[2]
	var configPath string

	// Proper argument parsing
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "-f" {
			if i+1 < len(os.Args) {
				configPath = os.Args[i+1]
				break
			} else {
				fmt.Println("Error: missing config file path after -f flag")
				os.Exit(1)
			}
		}
	}

	if err := createProject(projectName, configPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	printSuccess(projectName)
}

func printUsage() {
	fmt.Println(`Usage:
  tofo <project-name>                  # Basic project (legacy)
  tofo create <project-name>           # Basic project
  tofo create <project-name> -f <config.yaml>  # With config`)
}

func printSuccess(projectName string) {
	fmt.Printf(`
Project %s created successfully!

Next steps:
1. cd %s
2. deno run --allow-net --allow-read server.ts
`, projectName, projectName)
}

func createProject(projectName, configPath string) error {
    // 1. Deno check
    if err := checkDeno(); err != nil {
        return err
    }

    // 2. Create project directory
    if err := os.Mkdir(projectName, 0755); err != nil {
        return err
    }

    // 3. Download CORE files (always needed)
    for _, file := range coreFiles {
        if err := downloadAndSave(repoURL+file, filepath.Join(projectName, file)); err != nil {
            return fmt.Errorf("failed to download %s: %v", file, err)
        }
    }

    // 4. Process config if provided
    var envVars map[string]string
    var routes map[string]map[string]RouteConfig
    usingCustomConfig := configPath != ""

    if usingCustomConfig {
        cfg, err := loadConfig(configPath)
        if err != nil {
            return fmt.Errorf("config error: %v", err)
        }
        envVars = cfg.Project.Env
        routes = cfg.Project.Structure
    } else {
        // Download DEFAULT structure files when not using custom config
        defaultFiles := []string{
            "snaps/home/get.ts",
        }
        for _, file := range defaultFiles {
            if err := downloadAndSave(repoURL+file, filepath.Join(projectName, file)); err != nil {
                return fmt.Errorf("failed to download %s: %v", file, err)
            }
        }
    }

    // 5. Create .env file (merge defaults with custom)
    envContent := buildEnvContent(envVars)
    if err := os.WriteFile(filepath.Join(projectName, ".env"), []byte(envContent), 0644); err != nil {
        return err
    }

    // 6. Always download layout.eta to the routes directory
    basePath := "./snaps"
    if envVars != nil {
        if bp, ok := envVars["BASE_PATH"]; ok {
            basePath = bp
        }
    }
    layoutDest := filepath.Join(projectName, basePath, "layout.eta")
    if err := downloadAndSave(repoURL+"snaps/layout.eta", layoutDest); err != nil {
        return fmt.Errorf("failed to download layout.eta: %v", err)
    }

    // 7. Generate routes (only for custom config mode)
    if usingCustomConfig && routes != nil {
        for route, methods := range routes {
            for method, config := range methods {
                if err := createRouteFromTemplate(projectName, basePath, route, method, config.ViewName); err != nil {
                    return fmt.Errorf("failed to create route %s %s: %v", route, method, err)
                }
            }
        }
    }

    return nil
}


func buildEnvContent(customVars map[string]string) string {
    defaultVars := map[string]string{
        "PUBLIC_PATH":     "public",
        "PORT":           "8000",
        "BASE_PATH":      "./snaps",  // Consistent base path
        "HOME_PATH":      "home",
        "ALLOWED_METHODS": "GET,POST",
        "VIEWS_BASE":     "./snaps",  // Consistent with BASE_PATH
        "COMPONENTS_DIR": "components",
    }

    // Merge custom variables
    if customVars != nil {
        for k, v := range customVars {
            defaultVars[k] = v
        }
    }

    // Build env content
    var builder strings.Builder
    for k, v := range defaultVars {
        builder.WriteString(fmt.Sprintf("%s=%s\n", k, v))
    }
    return builder.String()
}

func createRouteFromTemplate(projectName, basePath, route, method, viewName string) error {
    fullPath := filepath.Join(projectName, basePath, route)
    routeFile := filepath.Join(fullPath, strings.ToLower(method)+".ts")

    // Determine which template to use
    templateFile := "snaps/home/get.ts"
    if viewName != "" {
        templateFile = "snaps/profile/login/get.ts"
        
        // Create view directory and template
        if err := createViewTemplate(projectName, basePath, route, viewName); err != nil {
            return err
        }
    }

    // Download and process route template
    content, err := downloadContent(repoURL + templateFile)
    if err != nil {
        return fmt.Errorf("failed to download route template: %v", err)
    }

    contentStr := string(content)
    if viewName != "" {
        contentStr = strings.ReplaceAll(contentStr, `"login"`, fmt.Sprintf(`"%s"`, viewName))
        contentStr = strings.ReplaceAll(contentStr, `"/profile/login"`, fmt.Sprintf(`"/%s"`, route))
    } else {
        contentStr = strings.ReplaceAll(contentStr, `"/home"`, fmt.Sprintf(`"/%s"`, route))
    }

    // Create route file
    if err := os.MkdirAll(filepath.Dir(routeFile), 0755); err != nil {
        return err
    }
    return os.WriteFile(routeFile, []byte(contentStr), 0644)
}

func createViewTemplate(projectName, basePath, route, viewName string) error {
    // Download the ETA template
    templateContent, err := downloadContent(repoURL + "snaps/posts/[id]/views/template.eta")
    if err != nil {
        return fmt.Errorf("failed to download view template: %v", err)
    }

    // Create views directory
    viewsDir := filepath.Join(projectName, basePath, route, "views")
    if err := os.MkdirAll(viewsDir, 0755); err != nil {
        return fmt.Errorf("failed to create views directory: %v", err)
    }

    // Create view file
    viewFile := filepath.Join(viewsDir, viewName+".eta")
    return os.WriteFile(viewFile, templateContent, 0644)
}

func downloadContent(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to download %s: %s", url, resp.Status)
    }

    return io.ReadAll(resp.Body)
}

// func createRouteFile(path, route, method, viewName string) error {
// 	routeFile := filepath.Join(path, strings.ToLower(method)+".ts")
// 	relPath := strings.Repeat("../", strings.Count(route, "/")+2) + "core"

// 	content := fmt.Sprintf(`import { Router } from "%s/Router"

// export default {
// 	path: "/%s",
// 	handler: (req: Request) => {`, relPath, route)

// 	if viewName != "" {
// 		content += fmt.Sprintf(`
// 		return Router.view("%s/%s")`, route, viewName)
// 	} else {
// 		content += `
// 		return new Response("Hello from ` + method + ` ` + route + ` route")`
// 	}

// 	content += `
// 	}
// }`

// 	return os.WriteFile(routeFile, []byte(content), 0644)
// }

func loadConfig(path string) (*ProjectConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %v", err)
    }

    var cfg ProjectConfig
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("invalid YAML: %v", err)
    }

    return &cfg, nil
}

// func generateRouteFiles(projectName, route string, methods map[string]RouteConfig) error {
// 	// Get BASE_PATH from env
// 	basePath, err := getEnvValue(filepath.Join(projectName, ".env"), "BASE_PATH")
// 	if err != nil || basePath == "" {
// 		basePath = "./snaps" // Default if not set
// 	}

// 	routePath := filepath.Join(projectName, basePath, route)
// 	if err := os.MkdirAll(routePath, 0755); err != nil {
// 		return fmt.Errorf("failed to create route directory: %v", err)
// 	}

// 	for method, config := range methods {
// 		// Create route file
// 		filePath := filepath.Join(routePath, strings.ToLower(method)+".ts")
// 		content := fmt.Sprintf(`import { Router } from "../core/Router"

// export default {
// 	path: "/%s",
// 	handler: (req: Request) => {
// 		return new Response("Hello from %s %s route")
// 	}
// }`, route, method, route)
// 		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
// 			return fmt.Errorf("failed to create route file: %v", err)
// 		}

// 		// Generate view if specified
// 		if config.ViewName != "" {
// 			viewPath := filepath.Join(routePath, "views", config.ViewName+".eta")
// 			if err := generateView(viewPath, config.ViewName); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

func generateView(path, viewName string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create views directory: %v", err)
	}

	content := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<title>%s</title>
</head>
<body>
	<h1>%s View</h1>
</body>
</html>`, viewName, viewName)
	return os.WriteFile(path, []byte(content), 0644)
}

func getEnvValue(envPath, key string) (string, error) {
	content, err := os.ReadFile(envPath)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, key+"=") {
			return strings.TrimPrefix(line, key+"="), nil
		}
	}
	return "", nil
}

func updateEnv(projectName string, newEnv map[string]string) error {
	envPath := filepath.Join(projectName, ".env")
	
	// 1. Read current .env
	currentContent, err := os.ReadFile(envPath)
	if err != nil {
		return err
	}
	
	// 2. Parse into map
	envMap := make(map[string]string)
	for _, line := range strings.Split(string(currentContent), "\n") {
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	
	// 3. Merge new values (overwriting defaults)
	for k, v := range newEnv {
		envMap[k] = v
	}
	
	// 4. Rebuild .env content
	var builder strings.Builder
	for k, v := range envMap {
		builder.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}
	
	return os.WriteFile(envPath, []byte(builder.String()), 0644)
}

func checkDeno() error {
	if !isDenoInstalled() {
		if err := installDeno(); err != nil {
			return err
		}
	}
	if !isDenoInstalled() {
		return fmt.Errorf("Deno installation failed - please install manually from https://deno.land")
	}
	return nil
}

func isDenoInstalled() bool {
	if _, err := exec.LookPath("deno"); err == nil {
		return true
	}
	
	// Additional check for Windows install location
	if runtime.GOOS == "windows" {
		home, err := os.UserHomeDir()
		if err != nil {
			return false
		}
		if _, err := os.Stat(filepath.Join(home, ".deno", "bin", "deno.exe")); err == nil {
			return true
		}
	}
	return false
}

func installDeno() error {
	fmt.Println("Deno not found. Installing Deno...")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		installScript := `
		[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
		$ProgressPreference = 'SilentlyContinue'
		$denoPath = "$env:USERPROFILE\.deno\bin"
		
		# Install Deno
		iwr https://deno.land/install.ps1 -useb | iex
		
		# Ensure it's in PATH
		if (-not ($env:Path -like "*$denoPath*")) {
			[System.Environment]::SetEnvironmentVariable(
				"Path",
				"$env:Path;$denoPath",
				[System.EnvironmentVariableTarget]::User
			)
			$env:Path += ";$denoPath"
		}
		`
		cmd = exec.Command("powershell", "-Command", installScript)

	case "darwin", "linux":
		cmd = exec.Command("sh", "-c", `
			set -e
			curl -fsSL https://deno.land/x/install/install.sh | sh
			echo 'export PATH="$HOME/.deno/bin:$PATH"' >> $HOME/.bashrc
			export PATH="$HOME/.deno/bin:$PATH"
		`)

	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if isDenoInstalled() {
			fmt.Println("Deno installed successfully (some non-critical errors ignored)")
			return nil
		}
		return fmt.Errorf("failed to install Deno: %v", err)
	}

	// Final verification
	if !isDenoInstalled() {
		return fmt.Errorf("Deno installation failed - please install manually from https://deno.land")
	}

	fmt.Println("Deno installed and added to PATH successfully!")
	return nil
}

func downloadAndSave(url, dest string) error {
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
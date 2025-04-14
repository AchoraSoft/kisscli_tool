<p align="center">
<img src="https://github.com/AchoraSoft/tofo/raw/master/logo.png" alt="TOFO Logo" width="250"/>
</p>

# TOFO | Think Once, Folw On

A minimalist CLI tool for quickly scaffolding Deno projects with sensible defaults.

## Features

- 🚀 One-command project initialization
- ⚡ Deno-powered backend
- 📦 Includes all necessary project structure
- 🔧 Configurable via environment variables

## Quick Install

### Linux/macOS (bash)

```bash
curl -fsSL https://github.com/AchoraSoft/tofocli_tool/releases/download/v1.3.1/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://github.com/AchoraSoft/tofocli_tool/releases/download/v1.3.1/install.ps1 | iex
```

## Manual Installation

1. Download the appropriate binary from our [releases page](https://github.com/AchoraSoft/tofocli_tool/releases)
2. Make it executable: `chmod +x tofo`
3. Move it to your PATH: `mv tofo /usr/local/bin/`

## Usage

```bash
# Create a new project by default
tofo create my-project

OR

tofo create my-project -f <your-yaml-configuration>

# Example yaml file content

project:
  structure:
    home:
      get:
        viewName: "index" # if yout endpoint is a page with view
      post: {} # in case if your endpoint is API call with JSON result
    api:
      get: {}
  env:
    PORT: "3000"
    BASE_PATH: "./routes"


# Navigate to project
cd my-project

# Run the project
deno run --allow-net --allow-read server.ts
```

## Uninstallation

### Linux/macOS

```bash
curl -fsSL https://github.com/AchoraSoft/tofocli_tool/releases/download/v1.3.1/uninstall.sh | sh
```

### Windows

```powershell
irm https://github.com/AchoraSoft/tofocli_tool/releases/download/v1.3.1/uninstall.ps1 | iex
```

## Requirements

- Unix-like system or Windows with WSL (recommended)
- curl (for installation)
- Internet connection
- Deno (will be installed automatically if missing)

## Project Structure

New projects will include:

```
my-project/
├── .env
├── server.ts
├── public/
├── routes/
│   ├── home
|   |----get.ts
│   └── layout.eta
```

## License

MIT License

Copyright (c) 2023 AchoraSoft

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

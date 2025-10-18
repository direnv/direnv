# Using direnv with GitHub Actions

direnv can be integrated into GitHub Actions workflows to manage environment variables from `.envrc` files. This is particularly useful for:

- Maintaining consistent environment variables between local development and CI/CD
- Managing secrets and configuration in a familiar way
- Reusing existing `.envrc` files in your CI pipeline

## Basic Usage

The `direnv export gha` command outputs environment variables in the format required by GitHub Actions. Here's a basic example:

```yaml
name: CI
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install direnv
        run: |
          curl -sfL https://direnv.net/install.sh | bash
          echo "$HOME/.local/bin" >> $GITHUB_PATH
      
      - name: Load environment variables
        run: |
          direnv allow .
          direnv export gha > "$GITHUB_ENV"
      
      - name: Run tests
        run: |
          # Your environment variables from .envrc are now available
          echo "DATABASE_URL=$DATABASE_URL"
          make test
```

## Installation in GitHub Actions

There are several ways to install direnv in your GitHub Actions workflow:

### Using the install script (recommended)

```yaml
- name: Install direnv
  run: |
    curl -sfL https://direnv.net/install.sh | bash
    echo "$HOME/.local/bin" >> $GITHUB_PATH
```

### Using package managers

On Ubuntu runners:
```yaml
- name: Install direnv
  run: sudo apt-get update && sudo apt-get install -y direnv
```

On macOS runners:
```yaml
- name: Install direnv
  run: brew install direnv
```

### Using a specific version

```yaml
- name: Install direnv
  run: |
    wget -O direnv https://github.com/direnv/direnv/releases/download/v2.34.0/direnv.linux-amd64
    chmod +x direnv
    sudo mv direnv /usr/local/bin/
```

## Security Considerations

1. **Always use `direnv allow`**: Even in CI, direnv requires explicit permission to load `.envrc` files
2. **Be careful with secrets**: Don't echo or log sensitive environment variables
3. **Use GitHub Actions secrets**: Store sensitive values in GitHub Actions secrets rather than committing them to your repository

## Troubleshooting

### Environment variables not available

Make sure you're using `> "$GITHUB_ENV"` to override to the GitHub environment file:

```yaml
# Correct
direnv export gha > "$GITHUB_ENV"

# Incorrect - this just prints to stdout
direnv export gha
```

### Permission denied errors

Ensure direnv is executable and in your PATH:

```yaml
- name: Install and setup direnv
  run: |
    curl -sfL https://direnv.net/install.sh | bash
    echo "$HOME/.local/bin" >> $GITHUB_PATH
    # Force PATH update in current step
    export PATH="$HOME/.local/bin:$PATH"
    direnv allow .
    direnv export gha >> "$GITHUB_ENV"
```

### Debugging

To debug issues, you can examine what direnv is doing:

```yaml
- name: Debug direnv
  run: |
    direnv version
    direnv status
    cat .envrc
    direnv allow .
    direnv export gha
```

## Example: Complete Node.js Workflow

Here's a complete example for a Node.js project:

```yaml
name: Node.js CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    strategy:
      matrix:
        node-version: [18.x, 20.x]
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Install direnv
      run: |
        curl -sfL https://direnv.net/install.sh | bash
        echo "$HOME/.local/bin" >> $GITHUB_PATH
    
    - name: Setup Node.js ${{ matrix.node-version }}
      uses: actions/setup-node@v4
      with:
        node-version: ${{ matrix.node-version }}
    
    - name: Load environment
      run: |
        # Create .envrc if it doesn't exist
        if [ ! -f .envrc ]; then
          cat > .envrc <<'EOF'
          export NODE_ENV=test
          export DATABASE_URL="postgresql://localhost/test_db"
          PATH_add node_modules/.bin
          EOF
        fi
        
        direnv allow .
        direnv export gha >> "$GITHUB_ENV"
    
    - name: Install dependencies
      run: npm ci
    
    - name: Run tests
      run: npm test
    
    - name: Run linter
      run: npm run lint
```

## See Also

- [GitHub Actions Environment Files](https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#environment-files)
- [direnv Documentation](https://direnv.net)
- [direnv-stdlib(1)](../man/direnv-stdlib.1.md)

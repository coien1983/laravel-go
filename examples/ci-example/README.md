# CI/CD Example for Laravel-Go Framework

This directory contains example CI/CD configurations for the Laravel-Go Framework.

## Files

- `ci.yml` - GitHub Actions CI/CD workflow example

## How to Use

### GitHub Actions

1. **Copy the CI file to your repository:**

   ```bash
   cp examples/ci-example/ci.yml .github/workflows/ci.yml
   ```

2. **Customize the configuration:**

   - Update the Docker image tags in the `docker` job
   - Modify the Go versions in the `test` job matrix
   - Adjust the build targets in the `build` job

3. **Set up secrets (optional):**

   - `DOCKER_USERNAME` - Your Docker Hub username
   - `DOCKER_PASSWORD` - Your Docker Hub password
   - `CODECOV_TOKEN` - Codecov token for coverage reporting

4. **Commit and push:**

   ```bash
   git add .github/workflows/ci.yml
   git commit -m "Add CI/CD workflow"
   git push origin main
   ```

## Workflow Jobs

### 1. Test Job

- Runs tests on multiple Go versions (1.20, 1.21, 1.22)
- Generates test coverage reports
- Uploads coverage to Codecov

### 2. Lint Job

- Runs `golint` for code style checking
- Runs `go vet` for static analysis
- Checks code formatting with `gofmt`

### 3. Build Job

- Builds the main project
- Builds example applications
- Uploads build artifacts

### 4. Security Job

- Runs `gosec` for security scanning
- Identifies potential security vulnerabilities

### 5. Docker Job (Main branch only)

- Builds and pushes Docker images
- Uses Docker Hub for image storage
- Includes caching for faster builds

### 6. Release Job (Tags only)

- Creates GitHub releases
- Builds cross-platform binaries
- Attaches binaries to releases

## Customization

### Adding New Jobs

```yaml
new-job:
  name: Custom Job
  runs-on: ubuntu-latest
  needs: [test]

  steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Custom step
      run: echo "Your custom command here"
```

### Modifying Triggers

```yaml
on:
  push:
    branches: [main, develop, feature/*]
  pull_request:
    branches: [main]
  schedule:
    - cron: "0 0 * * 0" # Weekly on Sunday
```

### Environment Variables

```yaml
env:
  GO_VERSION: "1.21"
  BUILD_FLAGS: "-ldflags=\"-s -w\""
```

## Best Practices

1. **Use specific versions** for actions and tools
2. **Cache dependencies** to speed up builds
3. **Run tests first** before other jobs
4. **Use matrix builds** for multiple Go versions
5. **Include security scanning** in your pipeline
6. **Set up proper secrets** for sensitive data

## Troubleshooting

### Common Issues

1. **Build failures:**

   - Check Go version compatibility
   - Verify all dependencies are available
   - Review test output for specific errors

2. **Docker build issues:**

   - Ensure Dockerfile exists and is valid
   - Check Docker Hub credentials
   - Verify image tags are correct

3. **Permission issues:**

   - Check repository secrets are set correctly
   - Verify GitHub token permissions
   - Ensure workflow has necessary permissions

### Debug Mode

To enable debug logging, add this to your workflow:

```yaml
env:
  ACTIONS_STEP_DEBUG: true
  ACTIONS_RUNNER_DEBUG: true
```

## Alternative CI/CD Tools

- **GitLab CI** - GitLab's built-in CI/CD
- **Jenkins** - Self-hosted CI/CD server
- **CircleCI** - Cloud-based CI/CD platform
- **Travis CI** - GitHub-focused CI/CD service

## Support

For questions about CI/CD configuration:

- Check the [GitHub Actions documentation](https://docs.github.com/en/actions)
- Review the [Go GitHub Actions examples](https://github.com/actions/setup-go)
- Open an issue in the Laravel-Go repository

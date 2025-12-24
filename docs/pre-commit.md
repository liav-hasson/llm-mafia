# Pre-Commit Hooks

Automated code quality checks that run before each commit.

---

## Setup

```bash
# Activate virtual environment
source src/.venv/bin/activate

# Install pre-commit
pip install pre-commit detect-secrets

# Generate secrets baseline (first time only)
detect-secrets scan > .secrets.baseline

# Install git hooks
pre-commit install
pre-commit install --hook-type commit-msg
```

---

## Hooks

| Hook | Language | Purpose |
|------|----------|---------|
| ruff | Python | Linting + auto-fix |
| ruff-format | Python | Code formatting |
| mypy | Python | Static type checking |
| golangci-lint | Go | Linting (multiple linters) |
| detect-secrets | All | Credential leak prevention |
| trailing-whitespace | All | Remove trailing spaces |
| end-of-file-fixer | All | Ensure newline at EOF |
| check-yaml/toml/json | All | Syntax validation |
| check-merge-conflict | All | Detect unresolved conflicts |
| conventional-pre-commit | All | Commit message format |

---

## Commit Message Format

Commits must follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <description>

[optional body]
```

**Allowed types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `refactor` - Code restructuring
- `test` - Adding/updating tests
- `chore` - Maintenance
- `style` - Formatting
- `perf` - Performance
- `ci` - CI/CD changes
- `build` - Build system

**Examples:**
```bash
git commit -m "feat: add player agent kafka consumer"
git commit -m "fix: correct vote tallying logic"
git commit -m "docs: update pre-commit setup instructions"
```

---

## Usage

```bash
# Run on staged files (automatic on commit)
git commit -m "message"

# Run on all files manually
pre-commit run --all-files

# Run specific hook
pre-commit run ruff --all-files

# Skip hooks (use sparingly)
git commit --no-verify -m "message"
```

---

## Configuration Files

- [.pre-commit-config.yaml](../.pre-commit-config.yaml) - Hook definitions
- [pyproject.toml](../pyproject.toml) - Python tool settings (ruff, mypy)
- [.golangci.yaml](../.golangci.yaml) - Go linter settings
- [.secrets.baseline](../.secrets.baseline) - Known secrets whitelist

---

## Troubleshooting

**Hook fails on commit:**
- Some hooks auto-fix issues. Re-stage files (`git add .`) and commit again.

**mypy missing type stubs:**
- Add stubs to `additional_dependencies` in `.pre-commit-config.yaml`

**Update hook versions:**
```bash
pre-commit autoupdate
```

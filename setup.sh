#!/bin/bash

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Function to show usage
usage() {
    cat << EOF
Usage: $0 <new-module-name>

This script renames the project from 'github.com/andrew-womeldorf/connect-boilerplate' to your specified module name.
It should only be run ONCE after copying or forking this project.

Arguments:
  <new-module-name>    The new fully qualified module name (e.g., github.com/user/project)

Examples:
  $0 github.com/myuser/my-api
  $0 git.company.com/team/project-api

EOF
}

# Check if new module name is provided
if [[ $# -ne 1 ]]; then
    log_error "New module name is required"
    usage
    exit 1
fi

NEW_MODULE="$1"
OLD_MODULE="github.com/andrew-womeldorf/connect-boilerplate"

# Validate new module name format
if [[ ! "$NEW_MODULE" =~ ^[a-zA-Z0-9._/-]+$ ]]; then
    log_error "Invalid module name format: $NEW_MODULE"
    log_error "Module name should contain only alphanumeric characters, dots, underscores, hyphens, and slashes"
    exit 1
fi

# Extract package name from module (last segment after last slash)
NEW_PACKAGE=$(basename "$NEW_MODULE")
OLD_PACKAGE="example"

log_info "Project setup configuration:"
echo "  Old module: $OLD_MODULE"
echo "  New module: $NEW_MODULE"
echo "  Old package: $OLD_PACKAGE"
echo "  New package: $NEW_PACKAGE"
echo

# Check if this appears to be the template project
if [[ ! -f "go.mod" ]] || ! grep -q "$OLD_MODULE" go.mod; then
    log_error "This doesn't appear to be the template project or it has already been renamed"
    log_error "Expected to find '$OLD_MODULE' in go.mod"
    exit 1
fi

# Check if setup has already been run
if grep -q "$NEW_MODULE" go.mod 2>/dev/null; then
    log_error "Setup appears to have already been run with module name: $NEW_MODULE"
    log_error "This script should only be run once"
    exit 1
fi

log_info "Starting project rename from '$OLD_MODULE' to '$NEW_MODULE'"
echo

# 1. Update go.mod
log_info "Updating Go module file: go.mod"
sed -i.bak "s|$OLD_MODULE|$NEW_MODULE|g" go.mod && rm go.mod.bak

# 2. Update all Go files
log_info "Updating Go source files..."
find . -name "*.go" -type f | while read -r file; do
    if grep -q "$OLD_MODULE" "$file"; then
        log_info "  Updating: $file"
        sed -i.bak "s|$OLD_MODULE|$NEW_MODULE|g" "$file" && rm "$file.bak"
    fi
done

# 3. Update proto files (before moving them)
log_info "Updating protobuf files..."
find proto -name "*.proto" -type f | while read -r file; do
    log_info "  Updating: $file"
    sed -i.bak "s|$OLD_MODULE|$NEW_MODULE|g" "$file"
    sed -i.bak "s|package $OLD_PACKAGE\\.v1|package $NEW_PACKAGE.v1|g" "$file"
    sed -i.bak "s|import \"$OLD_PACKAGE/v1/|import \"$NEW_PACKAGE/v1/|g" "$file"
    rm "$file.bak"
done

# 4. Move proto directory structure
if [[ -d "proto/$OLD_PACKAGE" ]]; then
    log_info "Moving proto directory: proto/$OLD_PACKAGE -> proto/$NEW_PACKAGE"
    mv "proto/$OLD_PACKAGE" "proto/$NEW_PACKAGE"
fi

# 5. Update buf.gen.yaml if it exists
if [[ -f "buf.gen.yaml" ]]; then
    log_info "Updating buf generation config: buf.gen.yaml"
    sed -i.bak "s|$OLD_MODULE|$NEW_MODULE|g" buf.gen.yaml && rm buf.gen.yaml.bak
fi

# 6. Update any other config files that might reference the old module
for file in buf.yaml mise.toml; do
    if [[ -f "$file" ]] && grep -q "$OLD_MODULE" "$file"; then
        log_info "Updating config file: $file"
        sed -i.bak "s|$OLD_MODULE|$NEW_MODULE|g" "$file" && rm "$file.bak"
    fi
done

echo
log_info "Project setup completed successfully!"
log_info "Next steps:"
echo "  1. Run: mise run proto:generate"
echo "  2. Run: mise run check"
echo "  3. Commit your changes: git add . && git commit -m 'Rename project to $NEW_MODULE'"
echo "  4. Remove this setup script: rm setup.sh"

#!/usr/bin/env sh

set -euo pipefail

TF_VERSION="${TF_VERSION:-1.13.3}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture to Terraform's naming convention
case "${ARCH}" in
    x86_64)
        TF_ARCH="amd64"
        ;;
    aarch64|arm64)
        TF_ARCH="arm64"
        ;;
    armv7l|armv6l)
        TF_ARCH="arm"
        ;;
    i386|i686)
        TF_ARCH="386"
        ;;
    *)
        echo "Unsupported architecture: ${ARCH}"
        exit 1
        ;;
esac

TF_PACKAGE="terraform_${TF_VERSION}_${OS}_${TF_ARCH}.zip"

# Install unzip if not present
apt-get update && apt-get install -y unzip

# Download terraform and checksums
curl -LO "https://releases.hashicorp.com/terraform/${TF_VERSION}/${TF_PACKAGE}"
curl -LO "https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_SHA256SUMS"

# Verify checksum
grep "${TF_PACKAGE}" "terraform_${TF_VERSION}_SHA256SUMS" | sha256sum -c -

# Install terraform
unzip "${TF_PACKAGE}"
mv terraform /usr/local/bin/
chmod +x /usr/local/bin/terraform

# Clean up
rm "${TF_PACKAGE}" "terraform_${TF_VERSION}_SHA256SUMS"

# Verify installation
terraform version

pip install terraform-local

set -x
tflocal -chdir=/terraform init -input=false
tflocal -chdir=/terraform apply -auto-approve

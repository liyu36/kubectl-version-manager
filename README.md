# kubectl-version-manager

kubectl-version-manager is a straightforward tool designed to manage various versions of the kubectl command.

It achieves this by creating symbolic link across different versions.

The tool offers convenient commands for installing, using, and removing specific versions, as well as displaying the current version and listing available versions of kubectl.

## Features

- **Install kubectl:** Easily install the kubectl command.

- **View installed version:** View the currently installed version of kubectl.

- **Use specific version:** Switch to and use a specified version of kubectl.

- **Display current version:** Check the currently active version of kubectl.

- **List available versions:** View a list of kubectl versions available for installation.

- **Remove specific version:** Uninstall a specified version of kubectl.

## Installation

``` bash
go install github.com/liyu36/kubectl-version-manager@latest
# or
version="v0.1.1"
wget -O kubectl-version-manager https://github.com/liyu36/kubectl-version-manager/releases/download/${version}/kubectl-version-manager-${version}-darwin-arm64
mv kubectl-version-manager /usr/local/bin/

# ~/.bash_profile
export PATH="~/.kube/bin:$PATH"
source <(kubectl-version-manager completion bash)
complete -F __start_kubectl-version-manager kvm
alias kvm="kubectl-version-manager"
```

## Usage

### Install kubectl

``` bash
kvm install <version>
```

### View installed version

``` bash
kvm show
```

### Use specific version

``` bash
kvm use <version>
```

### Display current version

``` bash
kvm current
```

### List available versions

``` bash
kvm list [regexp]
```

### Remove specific version

``` bash
kvm remove <version>
```

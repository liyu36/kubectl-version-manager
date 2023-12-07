#!/usr/bin/env bash
# author: liy
#
export PS4='\[\e[35m\]+ $(basename $0):${FUNCNAME}:$LINENO: \[\e[0m\]'
[[ "$debug" =~ ^(true|yes|on)$ ]] && set -x
set -e

which go &>/dev/null

arch_list=("amd64" "arm64")
os_list=("darwin" "linux")

tag="$(git describe --abbrev=0 --tags)"
name="$(basename $(pwd))"

function build() {
    output="${name}-${tag}-${os}-${arch}"
    echo "build ${output}"
    CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -o ${output}
}

function main() {
    go mod tidy
    for os in ${os_list[@]}; do
        for arch in ${arch_list[@]}; do
            build
        done
    done
}

main

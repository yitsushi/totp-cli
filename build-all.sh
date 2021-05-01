#!/bin/bash

ARCH=(amd64)
OS=(darwin linux windows)
TARGET_DIR=bin

for arch in ${ARCH[*]}; do
  for os in ${OS[*]}; do
    outname="${TARGET_DIR}/totp-cli-${os}-${arch}"

    if [[ "${os}" == "windows" ]]; then
      outname="${outname}.exe"
    fi

    GOOS="${os}" GOARCH="${arch}" CGO_ENABLED=0 go build -o "${outname}" .
  done
done

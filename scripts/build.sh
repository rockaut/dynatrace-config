#!/usr/bin/env bash
set -e

# try and use the correct MD5 lib (depending on user OS darwin/linux)
MD5=$(which md5 || which md5sum)

# for versioning
getCurrCommit() {
  echo `git rev-parse --short HEAD| tr -d "[ \r\n\']"`
}

# for versioning
getCurrTag() {
  echo `git describe --always --tags --abbrev=0 | tr -d "[v\r\n]"`
}

# remove any previous builds that may have failed
[ -e "./build/bin" ] && \
  echo "Cleaning up old builds..." && \
  rm -rf "./build/bin"

# build shaman
echo "Building dynatrace-config ..."
go build -ldflags="-s -X main.version=$(getCurrTag) -X main.commit=$(getCurrCommit)" \
  -o="./build/bin/dynatrace-config" -v -x \
  cmd/main.go

# look through each os/arch/file and generate an md5 for each
echo "Generating md5s..."
for file in $(ls ./build/bin/${os}/${arch}); do
    cat "./build/bin/${os}/${arch}/${file}" | ${MD5} | awk '{print $1}' >> "./build/bin/${os}/${arch}/${file}.md5"
done

#!/bin/sh

# Symlink to /go/src to make go environment happy enough to build
ln -s $GITHUB_WORKSPACE /go/src/app

cd /go/src/app

# GitHub actions build context do not contain the github workspace contents
# leaving the choice between:
#  - [THIS] Hackily building the binary on container run (cause $GITHUB_WORKSPACE gets volume mounted)
#  - Cloning the repository in the Dockerfile and then building (better, but I've chosen the lazy option for now)
echo "::group::go build"
go build -o bin/server ./cmd/server
go build -o bin/server-testwrapper ./.github/actions/tester/main.go
echo "::endgroup::"

cd $GITHUB_WORKSPACE

$GITHUB_WORKSPACE/bin/server-testwrapper

exit $?

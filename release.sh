#! /usr/bin/env bash

set -e -u -o pipefail

if [ "$#" -lt 1 ] ; then
	echo Usage: release.sh VERSION
	echo VERSION should be something like 1.0.3
	exit 1
fi

version="$1"
version_tag="v$version"

echo "-------------------------------------------------------------------------"
echo "Creating release for version $version
echo "-------------------------------------------------------------------------"

srcdir="$GOPATH/src/github.com/stathat/cmd"
cd $srcdir
if git tag -a $version_tag -m $version_tag ; then
	echo "Tagged stathat/cmd repo with $version_tag"
	git push --tags
else
	echo "git tag $version_tag failed on $srcdir, presumably it exists"
fi

# that's it for now...

echo "done."

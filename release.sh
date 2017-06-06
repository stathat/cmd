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
echo "Creating release for version $version"
echo "-------------------------------------------------------------------------"

srcdir="$GOPATH/src/github.com/stathat/cmd"
cd $srcdir
if git tag -a $version_tag -m $version_tag ; then
	echo "Tagged stathat/cmd repo with $version_tag"
	git push --tags
else
	echo "git tag $version_tag failed on $srcdir, presumably it exists"
fi

echo "-------------------------------------------------------------------------"
echo "Downloading release for version $version"
echo "-------------------------------------------------------------------------"

filename="v$version.tar.gz"
wget "https://github.com/stathat/cmd/archive/$filename"

echo "-------------------------------------------------------------------------"
echo "Calculating sha256"
echo "-------------------------------------------------------------------------"

shasum -a 256 $filename

rm $filename

#that's it for now...

echo "done."

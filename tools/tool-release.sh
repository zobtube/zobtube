#!/bin/bash

set -e

if [ -z "$1" ]
then
	echo "Usage $0 NEW_VERSION"
	exit 1
fi

NEW_VERSION=$1
BRANCH_NAME="chore-release-${NEW_VERSION}"

echo " - going to create release $NEW_VERSION, press enter to start"
read

echo " - checkout branch $BRANCH_NAME"
git checkout -b chore-release-$NEW_VERSION

echo " - generate release for version $NEW_VERSION"
./tools/release-generate.sh $NEW_VERSION

echo " - open editor to write release announcement"
RELEASE_MSG_FILE=$(mktemp)
trap 'rm -f "$RELEASE_MSG_FILE"' EXIT

cat > "$RELEASE_MSG_FILE" <<EOF
Hi everyone,



If you encounter any bug, feel free to drop an issue on the [Github issue page](https://github.com/zobtube/zobtube/issues).

As usual, it is available on the [Github release page](https://github.com/zobtube/zobtube/releases/tag/${NEW_VERSION}).

EOF

${EDITOR:-${VISUAL:-vi}} "$RELEASE_MSG_FILE"

VERSION_HEADER_LINE=$(grep -n "^## Version ${NEW_VERSION}$" CHANGELOG.md | head -1 | cut -d: -f1)
if [ -z "$VERSION_HEADER_LINE" ]
then
	echo "Could not find '## Version ${NEW_VERSION}' in CHANGELOG.md"
	exit 1
fi

{
	head -n "$VERSION_HEADER_LINE" CHANGELOG.md
	echo ""
	cat "$RELEASE_MSG_FILE"
	echo ""
	tail -n +$((VERSION_HEADER_LINE + 1)) CHANGELOG.md
} > CHANGELOG.md.tmp
mv CHANGELOG.md.tmp CHANGELOG.md

echo " - commit changes"
git commit -asS -m "misc: prepare change for release $NEW_VERSION"

echo " - push diff"
git push --set-upstream origin "`git branch --no-color 2>/dev/null | grep '*' | sed -e 's/\\* //'`:dev/sblablaha/`git branch --no-color 2>/dev/null | grep '*' | sed -e 's/\\* //'`"

echo " - now wait for tests to go all green and press enter"
read

echo " - checkout main"
git checkout main

echo " - pull"
git pull --stat

echo " - merge"
git merge --ff-only $BRANCH_NAME

echo " - tag"
git tag $NEW_VERSION

echo " - push"
git push

echo " - push tag"
git push --tags

echo " - cleanup release branch"
git branch -d $BRANCH_NAME

echo " - all done"

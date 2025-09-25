#!/bin/bash

set -e

finalLog=""

TAG_PREVIOUS=""
TAG_CURRENT=""

generateChangelogBetweenVersions() {
	local previousTag=$1
	local currentTag=$2
	local tag=$3

	tagLog=""
	logs=$(git log --no-decorate --no-color tags/${previousTag}..${currentTag} |grep -E '^.*?(commit|feat|chore|fix|doc)' |sed -r 's/\s\s+//g')
	commit=""
	enrichedLogs=""

	while read log
	do
		if echo "$log" |grep -q 'commit'
		then
			commit="$(echo $log |cut -d' ' -f2)"
			continue
		fi
		enrichedLogs="${enrichedLogs}* [${commit::7}](https://github.com/zobtube/zobtube/commit/$commit) $log"$'\n'
	done <<< "$logs"

	enrichedLogs=$(echo "$enrichedLogs" |sort)

	tagLog="${tagLog}"$'\n'
	tagLog="${tagLog}## Version $tag"$'\n'


	if echo "$enrichedLogs" |grep -qE '^\* \[[a-f0-9]{7}\]\([^)]+\) feat'
	then
		tagLog="${tagLog}### Features"$'\n'
		while read log
		do
			tagLog="${tagLog}${log}"$'\n'
		done <<< "$(echo "$enrichedLogs" |grep -E '^\* \[[a-f0-9]{7}\]\([^)]+\) feat')"
	fi
	if echo "$enrichedLogs" |grep -qE '^\* \[[a-f0-9]{7}\]\([^)]+\) chore'
	then
		tagLog="${tagLog}### Chores"$'\n'
		while read log
		do
			tagLog="${tagLog}${log}"$'\n'
		done <<< "$(echo "$enrichedLogs" |grep -E '^\* \[[a-f0-9]{7}\]\([^)]+\) chore')"
	fi
	if echo "$enrichedLogs" |grep -qE '^\* \[[a-f0-9]{7}\]\([^)]+\) fix'
	then
		tagLog="${tagLog}### Fixes"$'\n'
		while read log
		do
			tagLog="${tagLog}${log}"$'\n'
		done <<< "$(echo "$enrichedLogs" |grep -E '^\* \[[a-f0-9]{7}\]\([^)]+\) fix')"
	fi
	if echo "$enrichedLogs" |grep -qE '^\* \[[a-f0-9]{7}\]\([^)]+\) doc'
	then
		tagLog="${tagLog}### Documentation"$'\n'
		while read log
		do
			tagLog="${tagLog}${log}"$'\n'
		done <<< "$(echo "$enrichedLogs" |grep -E '^\* \[[a-f0-9]{7}\]\([^)]+\) doc')"
	fi

	finalLog="$tagLog$finalLog"

}

for tag in $(git tag --sort=creatordate)
do
	if [ -z "$TAG_PREVIOUS" ]
	then
		TAG_PREVIOUS=$tag
		continue
	fi

	TAG_CURRENT=$tag

	generateChangelogBetweenVersions $TAG_PREVIOUS $TAG_CURRENT $tag

	TAG_PREVIOUS=$tag
done

if [[ "$1" ]]
then
	generateChangelogBetweenVersions $TAG_CURRENT "" $1
fi

echo '# Changelog' > CHANGELOG.md
echo "$finalLog" >> CHANGELOG.md

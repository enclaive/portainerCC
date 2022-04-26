#!/bin/bash

# For reference see: https://portainer.atlassian.net/wiki/spaces/TECH/pages/570589194/Code+Freeze+Preparation

# Portainer (CE + EE)
#   Change version in package.json
#   Change APIVersion in portainer.go
#   Change @version in handler/handler.go

# This script requires jq
#   sudo apt-get install jq

CURRENT_VERSION=$(jq -r '.version' package.json)


# Parse the major, minor and patch versions
# out.
# You use it like this:
#    semver="3.4.5+xyz"
#    a=($(parse_semver "$semver"))
#    major=${a[0]}
#    minor=${a[1]}
#    patch=${a[2]}
#    printf "%-32s %4d %4d %4d\n" "$semver" $major $minor $patch
function ParseSemVer() {
    local token="$1"
    local major=0
    local minor=0
    local patch=0

    if [[ "$token" =~ ^([0-9]+)\.([0-9]+)\.([0-9]+) ]]; then
        major=${BASH_REMATCH[1]}
        minor=${BASH_REMATCH[2]}
        patch=${BASH_REMATCH[3]}
    fi

    echo "$major $minor $patch"
}

echo "Current Portainer version: ${CURRENT_VERSION}"
echo -n "New version [${CURRENT_VERSION}]: "
read NEW_VERSION

if [ -z "$NEW_VERSION" ]; then
    echo "Version left unchanged."
    exit 0
fi

a=($(ParseSemVer "$NEW_VERSION"))
major=${a[0]}
minor=${a[1]}
patch=${a[2]}

if [ "$major" == 0 ] && [ "$minor" == 0 ] && [ "$patch" = 0 ]; then
    echo "Invalid version format, must be major.minor.patch"
    exit 1
fi

echo "Version will be changed to: ${NEW_VERSION}"
# echo -n "Continue? [y/N]: "
# read CONFIRM

# if [ "$CONFIRM" != "y" ]; then
#     echo "Aborting"
#     exit 1
# fi


tmp=$(mktemp)

# Change version in package.json
filename="package.json"
jq --arg a "$NEW_VERSION" '.version = $a' package.json > "$tmp" && mv "$tmp" "$filename"
echo "updated $filename"

# Update portainer.go
filename="api/portainer.go"
sed -E "s/^([[:blank:]]*APIVersion[[:blank:]]*=[[:blank:]]*).*/\1\"$NEW_VERSION\"/" "$filename" > "$tmp" && mv "$tmp" "$filename"
echo "updated $filename"

# Change @version in handler/handler.go
filename="api/http/handler/handler.go"
sed -E "s|// @version .*|// @version $NEW_VERSION|" "$filename" > "$tmp" && mv "$tmp" "$filename"
echo "updated $filename"


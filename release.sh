#!/bin/bash

set -e

tag=$1


echo "Releasing version $tag"

# create git tag
git tag $1
git push origin master --follow-tags

# update go module version
curl https://proxy.golang.org/github.com/serverless-plus/tencent-serverless-go/@v/$tag.info
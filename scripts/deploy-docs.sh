#!/bin/bash

source ./scripts/download-mdbook.sh

bin/mdbook build docs/memberdashboard --dest-dir ../../.dist/docs/

GIT_REPO_URL=$(git config --get remote.origin.url)

cd .dist/docs
git init .
git remote add github $GIT_REPO_URL
git checkout -b gh-pages
git add .
git commit -am "Static site deploy"
git push github gh-pages --force
cd ../..
rm -rf .deploy

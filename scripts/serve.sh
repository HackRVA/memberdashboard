#!/bin/bash

npm run build --prefix web/
cp -R web/dist/* internal/fs/web/
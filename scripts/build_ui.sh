#!/bin/bash

npm run build --prefix web/
cp -R web/dist/* pkg/membermgr/ui/web/

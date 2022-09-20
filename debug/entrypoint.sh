#!/bin/sh

/dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec /usr/local/bin/api-catalog-harvester

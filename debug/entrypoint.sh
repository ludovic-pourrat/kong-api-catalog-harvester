#!/bin/sh
#
# Copyright (c) 2022 by Bank Lombard Odier & Co Ltd, Geneva, Switzerland. This software is subject
# to copyright protection under the laws of Switzerland and other countries. ALL RIGHTS RESERVED.
#
#

/dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec /usr/local/bin/api-catalog-harvester

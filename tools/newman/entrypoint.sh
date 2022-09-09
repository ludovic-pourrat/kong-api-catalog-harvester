#!/bin/sh

set -e
    
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' http://kong:8006/status)" != "200" ]]; do 
    sleep 5; 
done
  
cd /etc/newman

exec "${@:-/bin/sh}"

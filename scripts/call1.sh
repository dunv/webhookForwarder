#!/bin/bash

curl --location --request POST "$1" \
--header 'Content-Type: application/json' \
--data-raw '{
	"duration": 2000,
	"command": "testCommandFromHttp"
}'

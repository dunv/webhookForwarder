#!/bin/bash

curl --location --request POST 'http://localhost:8080/webhookCall' \
--header 'Content-Type: application/json' \
--data-raw '{
	"duration": 2000,
	"command": "testCommandFromHttp"
}'

#!/bin/bash

DOCUMENTA_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIxLTExLTEwVDEyOjQyOjM4LjgyNTAzMzgwNyswMTowMCIsInVzZXItZW1haWwiOiJ0ZXN0QGdtYWlsLmNvbSIsInVzZXItaWQiOjM4fQ.1cNCXfiZKoMUDmZfZwabqIPk18s9iufzf2KBCXRLFEM"

HEADER="Content-Type: application/json"

DATA='{"ID": 59, "content": "derp-nov-04"}'

curl -X PUT -H "$HEADER" -d "$DATA" --cookie "documentaLoginToken=$DOCUMENTA_TOKEN" http://localhost:3123/api/v1/comment | jq

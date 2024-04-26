#!/bin/bash

url="http://localhost:17000"

curl -X POST -d "white" "$url"
curl -X POST -d "bgrect 0.25 0.25 0.75 0.75" "$url"
curl -X POST -d "figure 0.5 0.5" "$url"
curl -X POST -d "green" "$url"
curl -X POST -d "figure 0.6 0.6" "$url"
curl -X POST -d "update" "$url"

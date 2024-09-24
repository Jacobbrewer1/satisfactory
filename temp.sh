#!/bin/bash

curl -s --insecure https://192.168.1.253:7777/api/v1 --header 'Authorization: Bearer ewoJInBsIjogIkFQSVRva2VuIgp9.9F4E3195CE4C3DD165790189574D3D34828A837FBCBE04485F961F9A8FED6ADF4F0386F01C5119F3FB1D3136F59D165A255955475E0138FB3B2CF52598671D3A' --header 'Content-Type: application/json' --data '{
    "function": "QueryServerState",
    "data": {
        "clientCustomData": ""
    }
}' | jq -c

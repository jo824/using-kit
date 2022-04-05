#!/bin/bash
curl -X GET http://localhost:8008/things
curl -X GET http://localhost:8008/thing/yik
curl -X GET http://localhost:8008/thing/doesNotExist

curl -X POST http://localhost:8008/thing \
     -H "Content-Type: application/json" \
     -d '{"id":"newadd" , "available": true}'

curl -X GET http://localhost:8008/thing/newadd
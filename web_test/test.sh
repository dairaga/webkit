#!/bin/bash

curl --include http://localhost:8080/test/categories
curl --include http://localhost:8080/test/categories/1
curl --include http://localhost:8080/test/categories/100
curl --include --request POST --header "Content-type: application/json" --data '{"name": "PC", "parent": 1}' http://localhost:8080/test/categories
curl --include --request PUT --header "Content-type: application/json" --data '{"name": "NB-2", "parent": 1}' http://localhost:8080/test/categories/2
curl --include --request DELETE http://localhost:8080/test/categories/2
curl --include http://localhost:8080/test/categories
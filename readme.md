Add Items
example command
curl --header "Content-Type: application/json"   --request POST   --data '{"produceCode":"upup-down-left-righ", "UnitPrice": "3.45", "Name":"Name"}' localhost:8080/addItem

Delete
curl --header "Content-Type: application/json"   --request DELETE  localhost:8080/delete/upup-down-left-righ

Get/Fetch
curl --header "Content-Type: application/json"   --request GET  localhost:8080/items

curl --header "Content-Type: application/json"   --request GET  localhost:8080/item/upup-down-left-righ

Unit testing
Run manually using 'go test'

Dockerfile
Built using alpine as base 
docker inspect `container name`
CI
test ran using github actions

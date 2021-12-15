Usage Routes 
   Server Will be listening on port 8080
   Add Produce Item - /addItem
   Delete Produce Item - /delete/{produce-code}
   Get Single Produce Item - /item/{produce-code}
   Get All Produce Items - /items

Installation and Startup
   Build dockerfile 'docker build . -t your-container-tag'
   docker run your-container-tag --name your-container-name
   It should be 172.17.0.2 because that may be the default
   If is's not Ip address can be retrieved by docker inspect `container name`, and from there it should be in network settings

Routes
   /item
Add Items
   Requires DataStruct, Return nothing, response code 201 if successful, code 400 if request was invalid
   Unit Price is optional
   Example Command
      curl -i --header "Content-Type: application/json"   --request POST   --data '{"produceCode":"upup-down-left-righ", "UnitPrice": "3.45", "Name":"Name"}' localhost:8080/addItem

Delete Item
   Requires ProduceCode, Return nothing, response code 200 if successfull, code 404 if request was on resource that did not exist
   Example Command 
      curl -i --header "Content-Type: application/json"   --request DELETE  localhost:8080/delete/upup-down-left-righ

Get/Fetch 
   Requires Nothing, Returns DataItems, response code 200 if successfull, code 404 if request was on resource that did not exist
   Example Commands
      All
         curl -i --header "Content-Type: application/json"   --request GET  localhost:8080/items
      Single Item
         curl -i --header "Content-Type: application/json"   --request GET  localhost:8080/item/upup-down-left-righ

Unit testing
   Run manually using 'go test'
   Functions Tested
      (Function Name) - param paramName paramType 
      Get All Items - param 
      Get One Item - param product code as string
      Add Item - param DataItem DataStruct
      Delete Item - param Id string

Dockerfile
   Built using go-alpine as base 
   server started by running 'docker run `containerTag`'
   Ip address can be retrieved docker inspect `container name`
CI
   Tests automatically ran using github actions
   Automatically builds dockerfile

# cookItAPI
API that return stations that corresponding to the items ID's passed in the request

If you want to run the dockerized version, you need to publish the port used, like that:
docker run -p 8080:8080 docker-go-cook_it-api


To test the API with curl:
curl --header "Content-Type: application/json" --request GET --data '{"itemIds": [127, 123]}' http://localhost:8080/stations

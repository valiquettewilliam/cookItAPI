package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type item struct {
	Id           uint16 `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	Volume       int    `json:"volume"`
	DeliveryWeek string `json:"deliveryWeek"`
	Station      string `json:"station"`
	Category     string `json:"category"`
}

type protein struct {
	Name    string `json:"name"`
	Code    rune   `json:"code"`
	Station string `json:"station"`
}

type resStationGet struct {
	Picks        []string `json:"picks"`
	Out_of_stock []string `json:"out-of-stock"`
}

type inRequestType struct {
	ItemsID []uint16 `json:"itemIds"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/stations", getStations).Methods("GET")
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func getStations(w http.ResponseWriter, r *http.Request) {

	//get needed info
	var allItems []item = getAllItems()
	var allProteins []protein = getAllproteins()

	//TODO check if valid id
	var inItemIDs inRequestType
	byteValue, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(byteValue, &inItemIDs); err != nil {
		panic(err)
	}

	//output value
	//TODO find better name for return value
	var retPicks []string
	var retOutOfStock []string

	var tmpMeatStr []string

	for _, itemID := range inItemIDs.ItemsID {
		for _, itemObj := range allItems {

			if itemID == itemObj.Id {
				retPicks = append(retPicks, itemObj.Station)
				if itemObj.Volume == 0 {
					retOutOfStock = append(retOutOfStock, itemObj.Station)
				}

				tmpMeatStr = strings.Split(itemObj.DisplayName, "-")
				//check if valid return string
				if len(tmpMeatStr) != 3 {
					// manage error better
					fmt.Println("item got invalid display name value")
					continue
				}

				if tmpMeatStr[1] != "" {
					//we got meat
					for _, meatletter := range tmpMeatStr[1] {
						for _, proteinObj := range allProteins {
							if proteinObj.Code == meatletter {
								retPicks = append(retPicks, proteinObj.Station)
							}
						}
					}

				}
			}

		}
	}

	var returnJSON resStationGet
	returnJSON.Picks = retPicks
	returnJSON.Out_of_stock = retOutOfStock

	json.NewEncoder(w).Encode(returnJSON)

}

func getAllItems() []item {
	response, _ := http.Get("https://chefcookit.proxy.beeceptor.com/items")
	byteValue, _ := ioutil.ReadAll(response.Body)
	var allItems []item
	if err := json.Unmarshal(byteValue, &allItems); err != nil {
		panic(err)
	}
	return allItems
}

func getAllproteins() []protein {
	response, _ := http.Get("https://chefcookit.proxy.beeceptor.com/proteins")
	byteValue, _ := ioutil.ReadAll(response.Body)
	var allProteins []protein
	if err := json.Unmarshal(byteValue, &allProteins); err != nil {
		panic(err)
	}
	return allProteins
}

func (p *protein) UnmarshalJSON(d []byte) error {
	tmpStruct := struct {
		Name    string      `json:"name"`
		Code    interface{} `json:"code"`
		Station string      `json:"station"`
	}{}

	if err := json.Unmarshal(d, &tmpStruct); err != nil {
		return err
	}

	//TODO verif error better

	//convert string to rune
	p.Code = []rune(tmpStruct.Code.(string))[0]

	p.Name = tmpStruct.Name

	p.Station = tmpStruct.Station

	return nil
}

func (in_req *inRequestType) UnmarshalJSON(d []byte) error {
	tmpStruct := struct {
		ItemsID []interface{} `json:"itemIds"`
	}{}

	if err := json.Unmarshal(d, &tmpStruct); err != nil {
		return err
	}

	//TODO verif error better

	//convert string to uint16

	for _, id := range tmpStruct.ItemsID {

		tmpID, ok := id.(float64)
		if !ok {
			return errors.New("Error in decoding of IDs in request: " + "not a valid type")
		}

		idUi16 := uint16(tmpID)
		in_req.ItemsID = append(in_req.ItemsID, uint16(idUi16))
	}

	return nil
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")

	handleRequests()
}

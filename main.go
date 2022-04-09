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

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/stations", getStations).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func returnError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(err.Error())
}

func getStations(w http.ResponseWriter, r *http.Request) {

	//get needed info
	var allItems []item = getAllItems()
	var allProteins []protein = getAllproteins()

	var inItemIDs inRequestType
	byteValue, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(byteValue, &inItemIDs); err != nil {
		returnError(w, err)
	}

	//output value
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(returnJSON)

}

func getAllItems() []item {
	response, _ := http.Get("https://chefcookit.proxy.beeceptor.com/items")
	byteValue, _ := ioutil.ReadAll(response.Body)
	var allItems []item
	if err := json.Unmarshal(byteValue, &allItems); err != nil {
		fmt.Println(err.Error())
	}
	return allItems
}

func getAllproteins() []protein {
	response, _ := http.Get("https://chefcookit.proxy.beeceptor.com/proteins")
	byteValue, _ := ioutil.ReadAll(response.Body)
	var allProteins []protein
	if err := json.Unmarshal(byteValue, &allProteins); err != nil {
		fmt.Println(err.Error())
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

	//convert float64 to uint16
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
	fmt.Println("REST API started")

	handleRequests()
}

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shed/bookservice/api"
)

type GetFactsRequest struct {
	UserId string `json:"user_id"`
}

func GetFacts(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body) /// Deserialize request
	var request GetFactsRequest
	json.Unmarshal(reqBody, &request)                                        ///// deserialize and map it to object
	json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 200, Error: ""}) ////write response to http writer
}

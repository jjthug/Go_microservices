package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "hit broken mf",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission (w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w,r,&requestPayload)
	if err != nil {
		app.errorJSON(w,err)
		return
	}

	switch requestPayload.Action {
	case "auth" :
		app.authenticate(w, requestPayload.Auth)

	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload){
	// create json to send to auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response,err := client.Do(request)
	if err!= nil {
		app.errorJSON(w,err)
		return
	}

	defer response.Body.Close()

	// return status code

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	} else  if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// response from auth service
	var jsonFromService jsonResponse

	// decode the json
	err = json.NewDecoder(response.Body).Decode((&jsonFromService))
	if err != nil{
		app.errorJSON(w, errors.New("err"))
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "authenticated"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted,payload)

}
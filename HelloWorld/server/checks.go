package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (server *Server) CheckLogin() map[string]string {
	data := map[string]interface{}{
		"userID": uuid.New(),
	}
	log.Println(data)
	var err error
	var bytesRepresentation []byte
	if bytesRepresentation, err = json.Marshal(data); err != nil {
		log.Fatalln(err)
	}
	var resp *http.Response
	if resp, err = http.Post("http://localhost:4000/login", "application/json", bytes.NewBuffer(bytesRepresentation)); err != nil {
		log.Fatalln(err)
	}
	var result map[string]string

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	log.Println(result["data"])
	return result
}

func (server *Server) CheckRefresh(loginResult map[string]string) {
	data := map[string]interface{}{
		"refreshToken": loginResult["refreshToken"],
		"accessToken":  loginResult["accessToken"],
	}

	var err error
	var bytesRepresentation []byte
	if bytesRepresentation, err = json.Marshal(data); err != nil {
		log.Fatalln(err)
	}
	var resp *http.Response
	if resp, err = http.Post("http://localhost:4000/refresh", "application/json", bytes.NewBuffer(bytesRepresentation)); err != nil {
		log.Fatalln(err)
	}
	var result map[string]string

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	log.Println(result["data"])
}

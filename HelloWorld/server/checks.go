package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

//Проверка логина. Формирует новый uuid пользователя и передаёт его по адресу /login
func (server *Server) CheckLogin() map[string]string {
	data := map[string]interface{}{
		"userID": uuid.New(),
	}
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

	return result
}

//Проверка обновления токенов. Передаёт полученные в результате логина данные по адресу /refresh
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
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var cars []Car

func getCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func deleteCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // Получаем параметры из URL, в данном случае - ID автомобиля
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for index := range cars {
		if index == id {
			cars = append(cars[:index], cars[index+1:]...)
			w.WriteHeader(http.StatusNoContent) // Возвращаем статус 204 No Content в случае успешного удаления
			return
		}
	}

	w.WriteHeader(http.StatusNotFound) // Возвращаем статус 404 Not Found, если автомобиль не найден
}

func updateCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	var updatedCar Car
	err = json.NewDecoder(r.Body).Decode(&updatedCar)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if id < 0 || id >= len(cars) {
		http.Error(w, "Car ID not found", http.StatusNotFound)
		return
	}

	cars[id] = updatedCar
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cars[id])
}

func addCar(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		RegNums []string `json:"regNums"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, regNum := range requestData.RegNums {
		newCar := Car{RegNum: regNum}
		cars = append(cars, newCar)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestData)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/cars", getCars).Methods("GET")
	router.HandleFunc("/cars/{id}", deleteCar).Methods("DELETE")
	router.HandleFunc("/cars/{id}", updateCar).Methods("PATCH")
	router.HandleFunc("/cars", addCar).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

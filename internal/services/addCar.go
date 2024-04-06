package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"

	"github.com/Kiruhanchik/TestovoyeGo/internal/config"
	"github.com/Kiruhanchik/TestovoyeGo/internal/models"
)

// AddCar godoc
// @Summary      AddCar
// @Description  AddCar
// @Accept       json
// @Produce      json
// @Param        request   body      models.CarRequest  true  "Body"
// @Success      200  {object} string
// @Failure      400  {object} string
// @Failure      500  {object} string
// @Router       /cars [post]
func AddCar(w http.ResponseWriter, r *http.Request) {
	slog.Info("start AddCar")

	// Получение данных об автомобиле из тела запроса
	var newCars models.CarRequest
	if err := json.NewDecoder(r.Body).Decode(&newCars); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Подключение к базе данных
	db, err := sql.Open("postgres", fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable", config.Cfg.User, config.Cfg.Password, config.Cfg.Host, config.Cfg.DbPort, config.Cfg.DbName))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	for _, regNum := range newCars.RegNums {

		requestURL := config.Cfg.CarApi + "?regNum=" + regNum
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var newCar models.Car
		err = json.Unmarshal(resBody, &newCar)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Выполнение запроса на добавление нового автомобиля в базу данных
		_, err = db.Exec("INSERT INTO cars (regNum, mark, model, year, owner_name, owner_surname, owner_patronymic) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			newCar.RegNum, newCar.Mark, newCar.Model, newCar.Year, newCar.Owner.Name, newCar.Owner.Surname, newCar.Owner.Patronymic)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	slog.Info("complete AddCar")
	w.WriteHeader(http.StatusCreated) // Отправка HTTP-статуса 201 (Created)
}

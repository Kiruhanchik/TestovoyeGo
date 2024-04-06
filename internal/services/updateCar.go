package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kiruhanchik/TestovoyeGo/internal/config"
	"github.com/Kiruhanchik/TestovoyeGo/internal/models"
	"github.com/gorilla/mux"
)

// UpdateCar godoc
// @Summary      UpdateCar
// @Description  UpdateCar
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Param        request   body      models.CarFilter  true  "Body"
// @Success      200  {object} string
// @Failure      400  {object} string
// @Failure      500  {object} string
// @Router       /cars/:id [patch]
func UpdateCar(w http.ResponseWriter, r *http.Request) {
	slog.Info("start UpdateCar")
	// Получение параметра ID из URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Получение данных обновления из тела запроса
	var updatedCar models.CarFilter
	if err := json.NewDecoder(r.Body).Decode(&updatedCar); err != nil {
		slog.Error(err.Error())
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

	// Выполнение запроса на обновление автомобиля по идентификатору
	request := "UPDATE cars SET "
	args := []any{id}
	i := 2
	if updatedCar.RegNum != nil {
		request += fmt.Sprintf(" regNum=$%v ", i)
		i++
		args = append(args, updatedCar.RegNum)
	}
	if updatedCar.Mark != nil {
		request += fmt.Sprintf(" mark=$%v ", i)
		i++
		args = append(args, updatedCar.Mark)
	}
	if updatedCar.Model != nil {
		request += fmt.Sprintf(" model=$%v ", i)
		i++
		args = append(args, updatedCar.Model)
	}
	if updatedCar.Year != nil {
		request += fmt.Sprintf(" year=$%v ", i)
		i++
		args = append(args, updatedCar.Year)
	}
	if updatedCar.Owner != nil && updatedCar.Owner.Name != nil {
		request += fmt.Sprintf(" owner_name=$%v ", i)
		i++
		args = append(args, updatedCar.Owner.Name)
	}
	if updatedCar.Owner != nil && updatedCar.Owner.Surname != nil {
		request += fmt.Sprintf(" owner_surname=$%v ", i)
		i++
		args = append(args, updatedCar.Owner.Surname)
	}
	if updatedCar.Owner != nil && updatedCar.Owner.Patronymic != nil {
		request += fmt.Sprintf(" owner_patronymic=$%v ", i)
		i++
		args = append(args, updatedCar.Owner.Patronymic)
	}
	request += " WHERE id = $1"
	_, err = db.Exec(request, args...)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("complete UpdateCar")
	w.WriteHeader(http.StatusOK)
}

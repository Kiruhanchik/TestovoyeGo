package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

	"github.com/Kiruhanchik/TestovoyeGo/internal/config"
	"github.com/Kiruhanchik/TestovoyeGo/internal/models"
)

// GetCars godoc
// @Summary      GetCars
// @Description  GetCars
// @Accept       json
// @Produce      json
// @Param        limit   query      int  true  "limit"
// @Param        offset   query      int  true  "offset"
// @Param        regNum   query      string  false  "regNum"
// @Param        mark   query      string  false  "mark"
// @Param        model   query      string  false  "model"
// @Param        year   query      int  false  "year"
// @Param        owner_name   query      string  false  "owner_name"
// @Param        owner_surname   query      string  false  "owner_surname"
// @Param        owner_patronymic   query      string  false  "owner_patronymic"
// @Success      200  {object} string
// @Failure      400  {object} string
// @Failure      500  {object} string
// @Router       /cars [get]
func GetCars(w http.ResponseWriter, r *http.Request) {
	slog.Info("start GetCars")

	// Получение параметров фильтрации и пагинации из запроса
	params := r.URL.Query()
	limitStr := params.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	offsetStr := params.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var searchCar models.CarFilter

	regNum := params.Get("regNum")
	if len(regNum) != 0 {
		searchCar.RegNum = &regNum
	}
	mark := params.Get("mark")
	if len(mark) != 0 {
		searchCar.Mark = &mark
	}
	model := params.Get("model")
	if len(model) != 0 {
		searchCar.Model = &model
	}
	yearStr := params.Get("year")
	if len(yearStr) != 0 {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		searchCar.Year = &year
	}
	owner_name := params.Get("owner_name")
	owner_surname := params.Get("owner_surname")
	owner_patronymic := params.Get("owner_patronymic")
	if len(owner_name) != 0 || len(owner_surname) != 0 || len(owner_patronymic) != 0 {
		searchCar.Owner = &models.PeopleFilter{}

		if len(owner_name) != 0 {
			searchCar.Owner.Name = &owner_name
		}
		if len(owner_surname) != 0 {
			searchCar.Owner.Surname = &owner_surname
		}
		if len(owner_patronymic) != 0 {
			searchCar.Owner.Patronymic = &owner_patronymic
		}
	}

	// Подключение к базе данных
	db, err := sql.Open("postgres", fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable", config.Cfg.User, config.Cfg.Password, config.Cfg.Host, config.Cfg.DbPort, config.Cfg.DbName))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Выполнение запроса к базе данных
	query := "SELECT * FROM cars "
	where := make([]string, 0)
	i := 3
	args := []any{
		limit,
		offset,
	}

	if searchCar.RegNum != nil {
		where = append(where, fmt.Sprintf(" regNum=$%v ", i))
		i++
		args = append(args, searchCar.RegNum)
	}
	if searchCar.Mark != nil {
		where = append(where, fmt.Sprintf(" mark=$%v ", i))
		i++
		args = append(args, searchCar.Mark)
	}
	if searchCar.Model != nil {
		where = append(where, fmt.Sprintf(" model=$%v ", i))
		i++
		args = append(args, searchCar.Model)
	}
	if searchCar.Year != nil {
		where = append(where, fmt.Sprintf(" year=$%v ", i))
		i++
		args = append(args, searchCar.Year)
	}
	if searchCar.Owner != nil && searchCar.Owner.Name != nil {
		where = append(where, fmt.Sprintf(" owner_name=$%v ", i))
		i++
		args = append(args, searchCar.Owner.Name)
	}
	if searchCar.Owner != nil && searchCar.Owner.Surname != nil {
		where = append(where, fmt.Sprintf(" owner_surname=$%v ", i))
		i++
		args = append(args, searchCar.Owner.Surname)
	}
	if searchCar.Owner != nil && searchCar.Owner.Patronymic != nil {
		where = append(where, fmt.Sprintf(" owner_patronymic=$%v ", i))
		i++
		args = append(args, searchCar.Owner.Patronymic)
	}

	if len(where) > 0 {
		query += " where " + strings.Join(where, " and ")
	}
	query += " LIMIT $1 OFFSET $2"
	rows, err := db.Query(query, args...)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Создание списка автомобилей на основе данных из базы
	var carList []models.Car
	for rows.Next() {
		var car models.Car
		// Сканирование данных из строки результата в структуру Car
		err := rows.Scan(&car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Добавление машины в список
		carList = append(carList, car)
	}

	// Конвертация списка в JSON и отправка ответа
	jsonResult, err := json.Marshal(carList)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	slog.Info("complete GetCars")
	w.Write(jsonResult)
}

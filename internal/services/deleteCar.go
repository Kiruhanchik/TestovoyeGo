package services

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kiruhanchik/TestovoyeGo/internal/config"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

// DeleteCar godoc
// @Summary      DeleteCar
// @Description  DeleteCar
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Success      200  {object} string
// @Failure      400  {object} string
// @Failure      500  {object} string
// @Router       /cars/:id [delete]
func DeleteCar(w http.ResponseWriter, r *http.Request) {
	slog.Info("start DeleteCar")
	// Получение параметра ID из URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Подключение к базе данных
	db, err := sql.Open("postgres", fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable", config.Cfg.User, config.Cfg.Password, config.Cfg.Host, config.Cfg.DbPort, config.Cfg.DbName))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Выполнение запроса на удаление автомобиля по идентификатору
	_, err = db.Exec("DELETE FROM cars WHERE id = $1", id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("complete DeleteCar")
	w.WriteHeader(http.StatusNoContent)
}

package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"movie-festival-app/database"
)

type AdminController struct {
	db *gorm.DB
}

func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{db}
}

func (c *AdminController) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie database.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		SendResponse(w, http.StatusBadRequest, "Failed to decode request body", nil)
		return
	}

	if err := c.db.Create(&movie).Error; err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to create movie", nil)
		return
	}
	SendResponse(w, http.StatusCreated, "Success", movie)
}

func (c *AdminController) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID := vars["id"]

	var movie database.Movie
	if err := c.db.First(&movie, movieID).Error; err != nil {
		SendResponse(w, http.StatusNotFound, "Movie not found", nil)
		return
	}

	var updatedMovie database.Movie
	if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
		SendResponse(w, http.StatusBadRequest, "Failed to decode request body", nil)
		return
	}

	movie.Title = updatedMovie.Title
	movie.Description = updatedMovie.Description
	movie.Duration = updatedMovie.Duration
	movie.Artists = updatedMovie.Artists
	movie.Genres = updatedMovie.Genres
	movie.WatchURL = updatedMovie.WatchURL

	if err := c.db.Save(&movie).Error; err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to update movie", nil)
		return
	}
	SendResponse(w, http.StatusOK, "Success", movie)
}

func (c *AdminController) GetMostViewedMovie(w http.ResponseWriter, r *http.Request) {
	var movie database.Movie
	err := c.db.Table("movies").
		Select("movies.*, SUM(viewerships.view_counts) as total_view_counts").
		Joins("JOIN viewerships ON viewerships.movie_id = movies.id").
		Group("movies.id").
		Order("total_view_counts DESC").
		Limit(1).
		Find(&movie).Error

	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to get most viewed movie", nil)
		return
	}

	SendResponse(w, http.StatusOK, "Success", movie)
}

func (c *AdminController) GetMostViewedGenre(w http.ResponseWriter, r *http.Request) {
	var genres []database.Genre
	err := c.db.Model(&database.Viewership{}).
		Select("genres.*, COUNT(*) as viewership").
		Joins("JOIN genres ON viewerships.genre_id = genres.id").
		Group("genres.id").
		Order("viewership DESC").
		Limit(1).
		Find(&genres)

	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to get most viewed genres", nil)
		return
	}

	SendResponse(w, http.StatusOK, "Success", genres)
}

func SendResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	respPayload := map[string]interface{}{
		"stat_code": statusCode,
		"stat_msg":  message,
		"data":      data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(respPayload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

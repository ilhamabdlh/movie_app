package controllers

import (
	"encoding/json"
	"errors"
	"movie-festival-app/database"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db}
}

func (c *UserController) ListMovies(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	offset := (page - 1) * limit
	var movies []database.Movie
	if err := c.db.Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to get movies", nil)
		return
	}

	SendResponse(w, http.StatusOK, "Success", movies)
}

func (c *UserController) SearchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("title")
	movies := []database.Movie{}

	if err := c.db.Where("title LIKE ?", "%"+query+"%").Find(&movies).Error; err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to search movies", nil)
		return
	}

	SendResponse(w, http.StatusOK, "Success", movies)
}

func (c *UserController) TrackViewership(w http.ResponseWriter, r *http.Request) {
	viewershipID := mux.Vars(r)["movie_id"]
	viewership := database.Viewership{}
	query := "SELECT * FROM viewerships WHERE id = ?"
	if err := c.db.Raw(query, viewershipID).Scan(&viewership).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			SendResponse(w, http.StatusNotFound, "Viewership not found", nil)
		} else {
			SendResponse(w, http.StatusInternalServerError, "Failed to track viewership", nil)
		}
		return
	}

	SendResponse(w, http.StatusOK, "Success", viewership)
}

func (c *UserController) CreateViewership(w http.ResponseWriter, r *http.Request) {
	var viewership database.Viewerships
	var viewershipx database.Viewershipx
	var viewerCount int

	if err := json.NewDecoder(r.Body).Decode(&viewershipx); err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to input viewership", nil)

		return
	}

	userID, err := GetUserIDFromToken(r)
	if err != nil {
		SendResponse(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	if viewershipx.MovieID == 0 || userID == 0 {
		SendResponse(w, http.StatusBadRequest, "Invalid viewership data", nil)
		return
	}
	existingViewership := database.Viewerships{}
	if err := c.db.Where("user_id = ? AND movie_id = ?", userID, viewershipx.MovieID).First(&existingViewership).Error; err != nil {
		viewerCount = 1
	} else {
		viewerCount = existingViewership.ViewCounts + 1
	}

	location, _ := time.LoadLocation("Asia/Jakarta")
	currentTime := time.Now().In(location)
	offset := 7 * time.Hour
	current := currentTime.Add(offset)

	parsedTime, _ := time.Parse("2006-01-02 15:04:05", viewershipx.EndTime)
	durations := int(parsedTime.Sub(current).Round(time.Minute).Minutes())

	if existingViewership.MovieID != 0 {
		if err := c.db.Model(&existingViewership).
			Updates(database.Viewerships{
				EndTime:    parsedTime,
				Duration:   durations,
				ViewCounts: viewerCount,
			}).Error; err != nil {
			SendResponse(w, http.StatusInternalServerError, "Failed to update viewership", nil)
			return
		}
	} else {
		newViewership := database.Viewerships{
			MovieID:    uint(viewershipx.MovieID),
			UserID:     userID,
			StartTime:  current,
			EndTime:    parsedTime,
			Duration:   durations,
			ViewCounts: viewerCount,
		}
		if err := c.db.Create(&newViewership).Error; err != nil {
			SendResponse(w, http.StatusInternalServerError, "Failed to create viewership", nil)
			return
		}
	}

	viewership.UserID = userID
	viewership.MovieID = uint(viewershipx.MovieID)
	viewership.StartTime = current
	viewership.EndTime = parsedTime
	viewership.Duration = durations
	viewership.ViewCounts = viewerCount

	SendResponse(w, http.StatusCreated, "Success", viewership)
}

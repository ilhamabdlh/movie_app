package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Movie struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Artists     string `json:"artists"`
	Genres      string `json:"genres"`
	WatchURL    string `json:"watch_url"`
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Vote struct {
	gorm.Model
	UserID  uint `json:"user_id"`
	MovieID uint `json:"movie_id"`
}

type Genre struct {
	gorm.Model
	Name string `json:"name"`
}

type Viewership struct {
	gorm.Model
	MovieID    uint      `json:"movie_id"`
	UserID     uint      `json:"user_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Duration   int       `json:"duration"`
	ViewCounts int       `json:"view_counts"`
}
type Viewerships struct {
	MovieID    uint      `json:"movie_id"`
	UserID     uint      `json:"user_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Duration   int       `json:"duration"`
	ViewCounts int       `json:"view_counts"`
}

type Viewershipx struct {
	MovieID int    `json:"movie_id"`
	EndTime string `json:"end_time"`
}

type MovieResult struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Genres      string `json:"genres"`
	VoteCount   int    `json:"vote_count"`
}

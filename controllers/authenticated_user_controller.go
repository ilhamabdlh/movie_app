package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"movie-festival-app/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticatedUserController struct {
	db *gorm.DB
}

func NewAuthenticatedUserController(db *gorm.DB) *AuthenticatedUserController {
	return &AuthenticatedUserController{db}
}

func (c *AuthenticatedUserController) Register(w http.ResponseWriter, r *http.Request) {
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		SendResponse(w, http.StatusBadRequest, "Failed to decode request body", nil)
		return
	}
	// Hash the user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to hash password", nil)
		return
	}
	user.Password = string(hashedPassword)

	if err := c.db.Create(&user).Error; err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to create user", nil)
		return
	}

	SendResponse(w, http.StatusCreated, "Success", user)
}

func (c *AuthenticatedUserController) Login(w http.ResponseWriter, r *http.Request) {
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to decode request body", nil)
		return
	}
	var foundUser database.User
	if err := c.db.Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
		SendResponse(w, http.StatusInternalServerError, "Invalid email or password", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		SendResponse(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": foundUser.ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	// Generate the JWT token string
	secretKey := []byte("the secret of kalimdor") // Replace with your secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		SendResponse(w, http.StatusUnauthorized, "Failed to generate token", nil)
		return
	}

	response := map[string]string{
		"token": signedToken,
	}
	SendResponse(w, http.StatusOK, "Success", response)
}
func (c *AuthenticatedUserController) Logout(w http.ResponseWriter, r *http.Request) {
	SendResponse(w, http.StatusOK, "Logout successful", nil)
}

func (c *AuthenticatedUserController) ListVotes(w http.ResponseWriter, r *http.Request) {
	var votes []database.Vote
	if err := c.db.Find(&votes).Error; err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to get user votes", nil)
		return
	}

	SendResponse(w, http.StatusOK, "Success", votes)
}

func (c *AuthenticatedUserController) VoteMovie(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromToken(r)
	if err != nil {
		SendResponse(w, http.StatusBadRequest, "cannot vote because youre not login", nil)
		return
	}

	vars := mux.Vars(r)
	movId := vars["movie_id"]
	movieID, _ := strconv.ParseUint(movId, 10, 64)

	if c.hasUserVotedMovie(userID, uint(movieID)) {
		SendResponse(w, http.StatusBadRequest, "User has already voted for this movie", nil)
		return
	}

	vote := database.Vote{
		UserID:  userID,
		MovieID: uint(movieID),
	}

	if err := c.db.Create(&vote).Error; err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to vote for the movie", nil)
		return
	}

	SendResponse(w, http.StatusOK, "Vote successful", nil)
}

func (c *AuthenticatedUserController) UnvoteMovie(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	movieID := mux.Vars(r)["movie_id"]

	query := `delete from votes where user_id = ? and movie_id= ?`

	er := c.db.Exec(query, userID, movieID)
	if er != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to unvote movi", nil)
		return
	}

	if err := c.db.Where("user_id = ? AND movie_id = ?", userID, movieID).Delete(&database.Vote{}).Error; err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to unvote movie", nil)
		return
	}
	SendResponse(w, http.StatusOK, "Movie unvoted successfully", nil)
}

func (c *AuthenticatedUserController) GetMostVotedMovie(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT movies.id, movies.title, movies.description, movies.genres, COUNT(votes.movie_id) AS vote_count
		FROM movies
		LEFT JOIN votes ON movies.id = votes.movie_id
		WHERE votes.deleted_at IS NULL
		GROUP BY movies.id
		ORDER BY vote_count DESC
		LIMIT 1
	`

	rows, err := c.db.Raw(query).Rows()
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to get most voted movie", nil)
		return
	}
	defer rows.Close()

	var movieResult database.MovieResult

	if rows.Next() {
		if err := rows.Scan(&movieResult.ID, &movieResult.Title, &movieResult.Description, &movieResult.Genres, &movieResult.VoteCount); err != nil {
			SendResponse(w, http.StatusInternalServerError, "Failed to scan rows", nil)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	SendResponse(w, http.StatusOK, "Success", movieResult)

}

func (c *AuthenticatedUserController) hasUserVotedMovie(userID, movieID uint) bool {
	var vote database.Vote
	err := c.db.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&vote).Error
	if err != nil {
		return false
	}

	return true
}
func GetUserIDFromToken(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("Missing authorization token")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("the secret of kalimdor"), nil
	})
	if err != nil {
		return 0, fmt.Errorf("Failed to parse token: %v", err)
	}

	if !token.Valid {
		return 0, errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid token claims")
	}

	userID, ok := claims["userID"].(float64)
	if !ok {
		return 0, errors.New("Invalid user ID in token")
	}

	userIDUint := uint(userID)

	return userIDUint, nil
}

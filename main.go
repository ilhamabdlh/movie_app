package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"movie-festival-app/controllers"
	"movie-festival-app/database"
)

func main() {
	// connection to database MySQL
	db, err := gorm.Open("mysql", "root:satu2tiga45@tcp(localhost:3306)/movie_festival_app?charset=utf8&parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Migrate database
	db.AutoMigrate(&database.Movie{}, &database.User{}, &database.Vote{}, &database.Viewership{})

	// Inisialisasi controller
	adminController := controllers.NewAdminController(db)
	userController := controllers.NewUserController(db)
	authenticatedUserController := controllers.NewAuthenticatedUserController(db)

	// Inisialisasi router
	router := mux.NewRouter()

	// Admin APIs
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/movies", adminController.CreateMovie).Methods("POST")
	adminRouter.HandleFunc("/movies/{id}", adminController.UpdateMovie).Methods("PUT")
	adminRouter.HandleFunc("/movies/most-viewed", adminController.GetMostViewedMovie).Methods("GET")
	adminRouter.HandleFunc("/movies/most-viewed-genre", adminController.GetMostViewedGenre).Methods("GET")

	// All Users APIs
	userRouter := router.PathPrefix("/movies").Subrouter()
	userRouter.HandleFunc("", userController.ListMovies).Methods("GET")
	userRouter.HandleFunc("/search", userController.SearchMovies).Methods("GET")
	userRouter.HandleFunc("/viewership/{movie_id}", userController.TrackViewership).Methods("GET")
	userRouter.HandleFunc("/viewership", userController.CreateViewership).Methods("POST")

	// Authenticated User APIs
	authUserRouter := router.PathPrefix("/users").Subrouter()
	authUserRouter.HandleFunc("/register", authenticatedUserController.Register).Methods("POST")
	authUserRouter.HandleFunc("/login", authenticatedUserController.Login).Methods("POST")
	authUserRouter.HandleFunc("/logout", authenticatedUserController.Logout).Methods("POST")
	authUserRouter.HandleFunc("/votes", authenticatedUserController.ListVotes).Methods("GET")
	authUserRouter.HandleFunc("/vote/{movie_id}", authenticatedUserController.VoteMovie).Methods("POST")
	authUserRouter.HandleFunc("/unvote/{movie_id}", authenticatedUserController.UnvoteMovie).Methods("POST")
	authUserRouter.HandleFunc("/most-voted-movie", authenticatedUserController.GetMostVotedMovie).Methods("GET")

	// run server HTTP
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

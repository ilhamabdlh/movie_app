# Movie Festival App

This is a backend application for a movie festival app. It provides APIs for managing movies, user authentication, voting, and viewership tracking.

## Requirements

- Go 1.15 or higher
- MySQL database

## Installation

1. Clone the repository:

git clone https://github.com/your-username/movie-festival-app.git

2. Navigate to the project directory:

cd movie-festival-app

3. Install the dependencies:

go mod download

4. Set up the database:

Create a MySQL database with the name movie_festival_app.
Update the database connection details in the main.go file.

5. Run the application:

go run main.go

The application will be accessible at http://localhost:8080.

API Endpoints
Admin APIs
Create Movie

Endpoint: POST /admin/movies
Description: Create a new movie.
Request Body: JSON representation of the movie.
Response: Created movie details.
Update Movie

Endpoint: PUT /admin/movies/{id}
Description: Update an existing movie.
Request Body: JSON representation of the updated movie.
Response: Updated movie details.
Get Most Viewed Movie

Endpoint: GET /admin/movies/most-viewed
Description: Get the most viewed movie.
Response: Details of the most viewed movie.
Get Most Viewed Genre

Endpoint: GET /admin/movies/most-viewed-genre
Description: Get the most viewed genre.
Response: Details of the most viewed genre.
User APIs
List Movies

Endpoint: GET /movies
Description: Get a list of all movies.
Response: List of movies.
Search Movies

Endpoint: GET /movies/search
Description: Search for movies based on title, description, artist, or genre.
Query Parameters: Search criteria.
Response: List of matching movies.
Track Viewership

Endpoint: GET /movies/viewership/{id}
Description: Get viewership details of a movie by ID.
Path Parameter: Movie ID.
Response: Viewership details.
Create Viewership

Endpoint: POST /movies/viewership
Description: Create viewership record for a movie.
Request Body: JSON representation of viewership details.
Response: Created viewership details.
Authenticated User APIs
Register

Endpoint: POST /users/register
Description: Register a new user.
Request Body: User registration details.
Response: User registration confirmation.
Login

Endpoint: POST /users/login
Description: Log in as an existing user.
Request Body: User login credentials.
Response: Login token.
Logout

Endpoint: POST /users/logout
Description: Log out the currently authenticated user.
Response: Logout confirmation.
List Votes

Endpoint: GET /users/votes
Description: Get a list of votes by the authenticated user.
Response: List of votes.
Vote Movie

Endpoint: POST /users/vote/{id}
Description: Vote for a movie by ID.
Path Parameter: Movie ID.
Response: Vote confirmation.
Unvote Movie

Endpoint: POST /users/unvote/{id}
Description: Remove the vote for a movie by ID.
Path Parameter: Movie ID.
Response: Unvote confirmation.
Get Most Voted Movie

Endpoint: GET /users/most-voted-movie
Description: Get the most voted movie.
Response: Details of the most voted movie.




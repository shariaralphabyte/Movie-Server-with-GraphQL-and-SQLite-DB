package database

import (
	"database/sql"
	"fmt"
	"log"
	"movie-app/internal/models"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDatabase() error {
	var err error
	DB, err = sql.Open("sqlite3", "./movies.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	// Seed sample data
	seedData()

	return nil
}

func EnsureDirector(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("director name is required")
	}

	var id string
	err := DB.QueryRow("SELECT id FROM directors WHERE name = ?", name).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return "", err
	}

	id = uuid.New().String()
	_, err = DB.Exec("INSERT INTO directors (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		err2 := DB.QueryRow("SELECT id FROM directors WHERE name = ?", name).Scan(&id)
		if err2 == nil {
			return id, nil
		}
		return "", err
	}

	return id, nil
}

func createTables() error {
	moviesTable := `
	CREATE TABLE IF NOT EXISTS movies (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		year INTEGER,
		rating REAL,
		duration INTEGER,
		genre TEXT,
		director TEXT,
		poster_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	directorsTable := `
	CREATE TABLE IF NOT EXISTS directors (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE
	);`

	actorsTable := `
	CREATE TABLE IF NOT EXISTS actors (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		birth_date TEXT,
		nationality TEXT,
		biography TEXT,
		profile_url TEXT
	);`

	movieActorsTable := `
	CREATE TABLE IF NOT EXISTS movie_actors (
		movie_id TEXT,
		actor_id TEXT,
		character_name TEXT,
		FOREIGN KEY (movie_id) REFERENCES movies(id),
		FOREIGN KEY (actor_id) REFERENCES actors(id),
		PRIMARY KEY (movie_id, actor_id)
	);`

	reviewsTable := `
	CREATE TABLE IF NOT EXISTS reviews (
		id TEXT PRIMARY KEY,
		movie_id TEXT,
		user_name TEXT NOT NULL,
		rating INTEGER CHECK(rating >= 1 AND rating <= 5),
		comment TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (movie_id) REFERENCES movies(id)
	);`

	tables := []string{moviesTable, directorsTable, actorsTable, movieActorsTable, reviewsTable}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	return nil
}

func seedData() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM movies").Scan(&count)
	if err != nil {
		return
	}
	if count >= 10 {
		return
	}

	movies := []models.Movie{
		{ID: "1", Title: "Inception", Description: "A thief who steals corporate secrets through dream-sharing technology.", Year: 2010, Rating: 8.8, Duration: 148, Genre: "Sci-Fi, Action", Director: "Christopher Nolan", PosterURL: "https://example.com/inception.jpg"},
		{ID: "2", Title: "The Shawshank Redemption", Description: "Two imprisoned men bond over a number of years.", Year: 1994, Rating: 9.3, Duration: 142, Genre: "Drama", Director: "Frank Darabont", PosterURL: "https://example.com/shawshank.jpg"},
		{ID: "3", Title: "The Dark Knight", Description: "Batman faces the Joker, a criminal mastermind wreaking havoc on Gotham.", Year: 2008, Rating: 9.0, Duration: 152, Genre: "Action, Crime", Director: "Christopher Nolan", PosterURL: "https://example.com/dark-knight.jpg"},
		{ID: "4", Title: "Interstellar", Description: "A team travels through a wormhole in search of a new home for humanity.", Year: 2014, Rating: 8.6, Duration: 169, Genre: "Sci-Fi, Adventure", Director: "Christopher Nolan", PosterURL: "https://example.com/interstellar.jpg"},
		{ID: "5", Title: "Fight Club", Description: "An insomniac forms an underground fight club that evolves into something much more.", Year: 1999, Rating: 8.8, Duration: 139, Genre: "Drama", Director: "David Fincher", PosterURL: "https://example.com/fight-club.jpg"},
		{ID: "6", Title: "Pulp Fiction", Description: "The lives of two mob hitmen, a boxer, and others intertwine in tales of violence and redemption.", Year: 1994, Rating: 8.9, Duration: 154, Genre: "Crime, Drama", Director: "Quentin Tarantino", PosterURL: "https://example.com/pulp-fiction.jpg"},
		{ID: "7", Title: "Forrest Gump", Description: "The life journey of Forrest Gump, a man with a low IQ but a big heart.", Year: 1994, Rating: 8.8, Duration: 142, Genre: "Drama, Romance", Director: "Robert Zemeckis", PosterURL: "https://example.com/forrest-gump.jpg"},
		{ID: "8", Title: "The Matrix", Description: "A hacker learns the shocking truth about reality and his role in the war against its controllers.", Year: 1999, Rating: 8.7, Duration: 136, Genre: "Sci-Fi, Action", Director: "Lana Wachowski, Lilly Wachowski", PosterURL: "https://example.com/matrix.jpg"},
		{ID: "9", Title: "Gladiator", Description: "A former Roman General seeks revenge after being betrayed.", Year: 2000, Rating: 8.5, Duration: 155, Genre: "Action, Drama", Director: "Ridley Scott", PosterURL: "https://example.com/gladiator.jpg"},
		{ID: "10", Title: "Parasite", Description: "A poor family schemes to become employed by a wealthy household.", Year: 2019, Rating: 8.5, Duration: 132, Genre: "Thriller, Drama", Director: "Bong Joon-ho", PosterURL: "https://example.com/parasite.jpg"},
	}

	for _, movie := range movies {
		if movie.Director != "" {
			_, _ = EnsureDirector(movie.Director)
		}

		_, err := DB.Exec(
			`INSERT OR IGNORE INTO movies (id, title, description, year, rating, duration, genre, director, poster_url)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			movie.ID, movie.Title, movie.Description, movie.Year, movie.Rating,
			movie.Duration, movie.Genre, movie.Director, movie.PosterURL,
		)
		if err != nil {
			log.Printf("Failed to insert movie: %v", err)
		}
	}

	actors := []models.Actor{
		{ID: "a1", Name: "Leonardo DiCaprio", BirthDate: "1974-11-11", Nationality: "American", Biography: "Actor and producer.", ProfileURL: "https://example.com/actors/leo"},
		{ID: "a2", Name: "Morgan Freeman", BirthDate: "1937-06-01", Nationality: "American", Biography: "Actor, director, narrator.", ProfileURL: "https://example.com/actors/morgan"},
		{ID: "a3", Name: "Christian Bale", BirthDate: "1974-01-30", Nationality: "British", Biography: "Actor.", ProfileURL: "https://example.com/actors/bale"},
		{ID: "a4", Name: "Keanu Reeves", BirthDate: "1964-09-02", Nationality: "Canadian", Biography: "Actor.", ProfileURL: "https://example.com/actors/keanu"},
		{ID: "a5", Name: "Brad Pitt", BirthDate: "1963-12-18", Nationality: "American", Biography: "Actor and producer.", ProfileURL: "https://example.com/actors/brad"},
		{ID: "a6", Name: "Tom Hanks", BirthDate: "1956-07-09", Nationality: "American", Biography: "Actor and filmmaker.", ProfileURL: "https://example.com/actors/hanks"},
	}

	for _, actor := range actors {
		_, err := DB.Exec(
			`INSERT OR IGNORE INTO actors (id, name, birth_date, nationality, biography, profile_url)
			VALUES (?, ?, ?, ?, ?, ?)`,
			actor.ID, actor.Name, actor.BirthDate, actor.Nationality, actor.Biography, actor.ProfileURL,
		)
		if err != nil {
			log.Printf("Failed to insert actor: %v", err)
		}
	}

	links := []models.MovieActor{
		{MovieID: "1", ActorID: "a1", CharacterName: "Cobb"},
		{MovieID: "2", ActorID: "a2", CharacterName: "Red"},
		{MovieID: "3", ActorID: "a3", CharacterName: "Bruce Wayne"},
		{MovieID: "8", ActorID: "a4", CharacterName: "Neo"},
		{MovieID: "5", ActorID: "a5", CharacterName: "Tyler Durden"},
		{MovieID: "7", ActorID: "a6", CharacterName: "Forrest Gump"},
	}

	for _, link := range links {
		_, err := DB.Exec(
			`INSERT OR IGNORE INTO movie_actors (movie_id, actor_id, character_name)
			VALUES (?, ?, ?)`,
			link.MovieID, link.ActorID, link.CharacterName,
		)
		if err != nil {
			log.Printf("Failed to insert movie_actors: %v", err)
		}
	}

	reviews := []models.Review{
		{ID: "r1", MovieID: "1", UserName: "alice", Rating: 5, Comment: "Mind-bending and brilliant."},
		{ID: "r2", MovieID: "2", UserName: "bob", Rating: 5, Comment: "One of the best movies ever made."},
		{ID: "r3", MovieID: "3", UserName: "charlie", Rating: 5, Comment: "Legendary superhero film."},
		{ID: "r4", MovieID: "8", UserName: "diana", Rating: 4, Comment: "A sci-fi classic."},
		{ID: "r5", MovieID: "10", UserName: "eve", Rating: 5, Comment: "Masterpiece."},
	}

	for _, review := range reviews {
		_, err := DB.Exec(
			`INSERT OR IGNORE INTO reviews (id, movie_id, user_name, rating, comment)
			VALUES (?, ?, ?, ?, ?)`,
			review.ID, review.MovieID, review.UserName, review.Rating, review.Comment,
		)
		if err != nil {
			log.Printf("Failed to insert review: %v", err)
		}
	}
}

func CloseDatabase() {
	if DB != nil {
		DB.Close()
	}
}

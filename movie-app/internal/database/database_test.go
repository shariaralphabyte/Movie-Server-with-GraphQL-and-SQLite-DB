package database

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func TestInsertTenMoviesWithDetails(t *testing.T) {
	var err error
	DB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer DB.Close()

	if err := createTables(); err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	type movieSeed struct {
		id          string
		title       string
		description string
		year        int
		rating      float64
		duration    int
		genre       string
		director    string
		posterURL   string
	}

	movies := []movieSeed{
		{"m1", "Movie 1", "Desc 1", 2001, 7.1, 101, "Drama", "Director A", "https://example.com/m1.jpg"},
		{"m2", "Movie 2", "Desc 2", 2002, 7.2, 102, "Action", "Director B", "https://example.com/m2.jpg"},
		{"m3", "Movie 3", "Desc 3", 2003, 7.3, 103, "Comedy", "Director C", "https://example.com/m3.jpg"},
		{"m4", "Movie 4", "Desc 4", 2004, 7.4, 104, "Thriller", "Director D", "https://example.com/m4.jpg"},
		{"m5", "Movie 5", "Desc 5", 2005, 7.5, 105, "Sci-Fi", "Director E", "https://example.com/m5.jpg"},
		{"m6", "Movie 6", "Desc 6", 2006, 7.6, 106, "Drama", "Director A", "https://example.com/m6.jpg"},
		{"m7", "Movie 7", "Desc 7", 2007, 7.7, 107, "Action", "Director B", "https://example.com/m7.jpg"},
		{"m8", "Movie 8", "Desc 8", 2008, 7.8, 108, "Comedy", "Director C", "https://example.com/m8.jpg"},
		{"m9", "Movie 9", "Desc 9", 2009, 7.9, 109, "Thriller", "Director D", "https://example.com/m9.jpg"},
		{"m10", "Movie 10", "Desc 10", 2010, 8.0, 110, "Sci-Fi", "Director E", "https://example.com/m10.jpg"},
	}

	for _, m := range movies {
		if _, err := EnsureDirector(m.director); err != nil {
			t.Fatalf("failed to ensure director %q: %v", m.director, err)
		}

		_, err := DB.Exec(
			`INSERT INTO movies (id, title, description, year, rating, duration, genre, director, poster_url)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			m.id, m.title, m.description, m.year, m.rating, m.duration, m.genre, m.director, m.posterURL,
		)
		if err != nil {
			t.Fatalf("failed to insert movie %s: %v", m.id, err)
		}

		// 2 actors per movie
		for i := 0; i < 2; i++ {
			actorID := uuid.New().String()
			actorName := m.title + " Actor " + string(rune('A'+i))
			_, err := DB.Exec(
				`INSERT INTO actors (id, name, birth_date, nationality, biography, profile_url)
				VALUES (?, ?, ?, ?, ?, ?)`,
				actorID, actorName, "1990-01-01", "Test", "Bio", "https://example.com/actor",
			)
			if err != nil {
				t.Fatalf("failed to insert actor for movie %s: %v", m.id, err)
			}

			_, err = DB.Exec(
				`INSERT INTO movie_actors (movie_id, actor_id, character_name)
				VALUES (?, ?, ?)`,
				m.id, actorID, "Character",
			)
			if err != nil {
				t.Fatalf("failed to link actor for movie %s: %v", m.id, err)
			}
		}

		// 2 reviews per movie
		for i := 0; i < 2; i++ {
			reviewID := uuid.New().String()
			_, err := DB.Exec(
				`INSERT INTO reviews (id, movie_id, user_name, rating, comment)
				VALUES (?, ?, ?, ?, ?)`,
				reviewID, m.id, "user", 5, "great",
			)
			if err != nil {
				t.Fatalf("failed to insert review for movie %s: %v", m.id, err)
			}
		}
	}

	assertCount := func(table string, want int) {
		var c int
		if err := DB.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&c); err != nil {
			t.Fatalf("failed to count %s: %v", table, err)
		}
		if c != want {
			t.Fatalf("expected %d rows in %s, got %d", want, table, c)
		}
	}

	assertCount("movies", 10)
	assertCount("reviews", 20)
	assertCount("movie_actors", 20)
	assertCount("actors", 20)

	// We used 5 unique directors.
	assertCount("directors", 5)
}

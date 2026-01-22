package resolvers

import (
	"database/sql"
	"fmt"
	"log"
	"movie-app/internal/database"
	"movie-app/internal/models"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

var actorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Actor",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.ID},
		"name":        &graphql.Field{Type: graphql.String},
		"birth_date":  &graphql.Field{Type: graphql.String},
		"nationality": &graphql.Field{Type: graphql.String},
		"biography":   &graphql.Field{Type: graphql.String},
		"profile_url": &graphql.Field{Type: graphql.String},
	},
})

var reviewType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Review",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.ID},
		"movie_id":   &graphql.Field{Type: graphql.ID},
		"user_name":  &graphql.Field{Type: graphql.String},
		"rating":     &graphql.Field{Type: graphql.Int},
		"comment":    &graphql.Field{Type: graphql.String},
		"created_at": &graphql.Field{Type: graphql.String},
	},
})

var movieType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Movie",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.ID},
		"title":       &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"year":        &graphql.Field{Type: graphql.Int},
		"rating":      &graphql.Field{Type: graphql.Float},
		"duration":    &graphql.Field{Type: graphql.Int},
		"genre":       &graphql.Field{Type: graphql.String},
		"director":    &graphql.Field{Type: graphql.String},
		"poster_url":  &graphql.Field{Type: graphql.String},
		"created_at":  &graphql.Field{Type: graphql.String},
		"updated_at":  &graphql.Field{Type: graphql.String},
		"actors": &graphql.Field{
			Type: graphql.NewList(actorType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				movie, ok := p.Source.(*models.Movie)
				if !ok {
					return []models.Actor{}, nil
				}
				return getActorsForMovie(movie.ID)
			},
		},
		"reviews": &graphql.Field{
			Type: graphql.NewList(reviewType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				movie, ok := p.Source.(*models.Movie)
				if !ok {
					return []models.Review{}, nil
				}
				return getReviewsForMovie(movie.ID)
			},
		},
	},
})

var paginationInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PaginationInfo",
	Fields: graphql.Fields{
		"page":        &graphql.Field{Type: graphql.Int},
		"limit":       &graphql.Field{Type: graphql.Int},
		"total":       &graphql.Field{Type: graphql.Int},
		"total_pages": &graphql.Field{Type: graphql.Int},
	},
})

var moviesResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MoviesResult",
	Fields: graphql.Fields{
		"movies":     &graphql.Field{Type: graphql.NewList(movieType)},
		"pagination": &graphql.Field{Type: paginationInfoType},
	},
})

var movieInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "MovieInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"title":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"description": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"year":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
		"rating":      &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
		"duration":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
		"genre":       &graphql.InputObjectFieldConfig{Type: graphql.String},
		"director":    &graphql.InputObjectFieldConfig{Type: graphql.String},
		"poster_url":  &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

var movieFilterType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "MovieFilter",
	Fields: graphql.InputObjectConfigFieldMap{
		"genre":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		"min_year":   &graphql.InputObjectFieldConfig{Type: graphql.Int},
		"max_year":   &graphql.InputObjectFieldConfig{Type: graphql.Int},
		"min_rating": &graphql.InputObjectFieldConfig{Type: graphql.Float},
		"search":     &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

var reviewInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ReviewInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"movie_id":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.ID)},
		"user_name": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"rating":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
		"comment":   &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

var actorInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ActorInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"birth_date":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"nationality": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"biography":   &graphql.InputObjectFieldConfig{Type: graphql.String},
		"profile_url": &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

var reviewCreateInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ReviewCreateInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"user_name": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"rating":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
		"comment":   &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

var movieWithDetailsInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "MovieWithDetailsInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"movie":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(movieInputType)},
		"actors":  &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.NewNonNull(actorInputType))},
		"reviews": &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.NewNonNull(reviewCreateInputType))},
	},
})

func GetMovie(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required")
	}

	movie, err := getMovieByID(id)
	if err != nil {
		return nil, err
	}

	// Load actors for the movie
	actors, err := getActorsForMovie(id)
	if err != nil {
		log.Printf("Failed to load actors: %v", err)
	}
	movie.Actors = actors

	// Load reviews for the movie
	reviews, err := getReviewsForMovie(id)
	if err != nil {
		log.Printf("Failed to load reviews: %v", err)
	}
	movie.Reviews = reviews

	return movie, nil
}

func GetMovies(p graphql.ResolveParams) (interface{}, error) {
	page := 1
	limit := 10

	if p.Args["page"] != nil {
		page = p.Args["page"].(int)
	}
	if p.Args["limit"] != nil {
		limit = p.Args["limit"].(int)
	}

	offset := (page - 1) * limit

	// Build query with filters
	query := "SELECT * FROM movies WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM movies WHERE 1=1"
	args := []interface{}{}
	countArgs := []interface{}{}

	if p.Args["filter"] != nil {
		filter := p.Args["filter"].(map[string]interface{})

		if genre, ok := filter["genre"].(string); ok && genre != "" {
			query += " AND genre LIKE ?"
			countQuery += " AND genre LIKE ?"
			args = append(args, "%"+genre+"%")
			countArgs = append(countArgs, "%"+genre+"%")
		}

		if minYear, ok := filter["min_year"].(int); ok && minYear > 0 {
			query += " AND year >= ?"
			countQuery += " AND year >= ?"
			args = append(args, minYear)
			countArgs = append(countArgs, minYear)
		}

		if maxYear, ok := filter["max_year"].(int); ok && maxYear > 0 {
			query += " AND year <= ?"
			countQuery += " AND year <= ?"
			args = append(args, maxYear)
			countArgs = append(countArgs, maxYear)
		}

		if minRating, ok := filter["min_rating"].(float64); ok && minRating > 0 {
			query += " AND rating >= ?"
			countQuery += " AND rating >= ?"
			args = append(args, minRating)
			countArgs = append(countArgs, minRating)
		}

		if search, ok := filter["search"].(string); ok && search != "" {
			query += " AND (title LIKE ? OR description LIKE ? OR director LIKE ?)"
			countQuery += " AND (title LIKE ? OR description LIKE ? OR director LIKE ?)"
			searchTerm := "%" + search + "%"
			args = append(args, searchTerm, searchTerm, searchTerm)
			countArgs = append(countArgs, searchTerm, searchTerm, searchTerm)
		}
	}

	// Get total count
	var total int
	err := database.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count movies: %v", err)
	}

	// Add pagination
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query movies: %v", err)
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID, &movie.Title, &movie.Description, &movie.Year, &movie.Rating,
			&movie.Duration, &movie.Genre, &movie.Director, &movie.PosterURL,
			&movie.CreatedAt, &movie.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %v", err)
		}
		movies = append(movies, movie)
	}

	totalPages := (total + limit - 1) / limit

	result := map[string]interface{}{
		"movies": movies,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}

	return result, nil
}

func SearchMovies(p graphql.ResolveParams) (interface{}, error) {
	query, ok := p.Args["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query is required")
	}

	page := 1
	limit := 10
	if p.Args["page"] != nil {
		page = p.Args["page"].(int)
	}
	if p.Args["limit"] != nil {
		limit = p.Args["limit"].(int)
	}

	offset := (page - 1) * limit

	searchTerm := "%" + query + "%"
	sqlQuery := `
		SELECT * FROM movies 
		WHERE title LIKE ? OR description LIKE ? OR director LIKE ? OR genre LIKE ?
		ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	rows, err := database.DB.Query(sqlQuery, searchTerm, searchTerm, searchTerm, searchTerm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %v", err)
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID, &movie.Title, &movie.Description, &movie.Year, &movie.Rating,
			&movie.Duration, &movie.Genre, &movie.Director, &movie.PosterURL,
			&movie.CreatedAt, &movie.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %v", err)
		}
		movies = append(movies, movie)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*) FROM movies 
		WHERE title LIKE ? OR description LIKE ? OR director LIKE ? OR genre LIKE ?
	`
	var total int
	err = database.DB.QueryRow(countQuery, searchTerm, searchTerm, searchTerm, searchTerm).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count movies: %v", err)
	}

	totalPages := (total + limit - 1) / limit

	result := map[string]interface{}{
		"movies": movies,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}

	return result, nil
}

func CreateSchema() (graphql.Schema, error) {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"movie": &graphql.Field{
				Type: movieType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: GetMovie,
			},
			"movies": &graphql.Field{
				Type: moviesResultType,
				Args: graphql.FieldConfigArgument{
					"page":   &graphql.ArgumentConfig{Type: graphql.Int},
					"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
					"filter": &graphql.ArgumentConfig{Type: movieFilterType},
				},
				Resolve: GetMovies,
			},
			"searchMovies": &graphql.Field{
				Type: moviesResultType,
				Args: graphql.FieldConfigArgument{
					"query": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"page":  &graphql.ArgumentConfig{Type: graphql.Int},
					"limit": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: SearchMovies,
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createMovie": &graphql.Field{
				Type: movieType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(movieInputType)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input, ok := p.Args["input"].(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("input is required")
					}

					id := uuid.New().String()

					if director, ok := input["director"].(string); ok && director != "" {
						if _, err := database.EnsureDirector(director); err != nil {
							return nil, fmt.Errorf("failed to ensure director: %v", err)
						}
					}

					_, err := database.DB.Exec(`
                        INSERT INTO movies (id, title, description, year, rating, duration, genre, director, poster_url, created_at, updated_at)
                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
						id,
						input["title"].(string),
						input["description"],
						input["year"].(int),
						input["rating"].(float64),
						input["duration"].(int),
						input["genre"],
						input["director"],
						input["poster_url"],
					)
					if err != nil {
						return nil, fmt.Errorf("failed to create movie: %v", err)
					}

					return getMovieByID(id)
				},
			},
			"updateMovie": &graphql.Field{
				Type: movieType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(movieInputType)},
				},
				Resolve: UpdateMovie,
			},
			"deleteMovie": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: DeleteMovie,
			},
			"createReview": &graphql.Field{
				Type: reviewType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(reviewInputType)},
				},
				Resolve: CreateReview,
			},
			"createMovieWithDetails": &graphql.Field{
				Type: movieType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(movieWithDetailsInputType)},
				},
				Resolve: CreateMovieWithDetails,
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
}

func CreateMovieWithDetails(p graphql.ResolveParams) (interface{}, error) {
	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("input is required")
	}

	movieInput := input["movie"].(map[string]interface{})
	var actorsInput []interface{}
	if input["actors"] != nil {
		if v, ok := input["actors"].([]interface{}); ok {
			actorsInput = v
		}
	}
	var reviewsInput []interface{}
	if input["reviews"] != nil {
		if v, ok := input["reviews"].([]interface{}); ok {
			reviewsInput = v
		}
	}

	movieID := uuid.New().String()

	if director, ok := movieInput["director"].(string); ok && director != "" {
		if _, err := database.EnsureDirector(director); err != nil {
			return nil, fmt.Errorf("failed to ensure director: %v", err)
		}
	}

	// Insert movie
	_, err := database.DB.Exec(`
        INSERT INTO movies (id, title, description, year, rating, duration, genre, director, poster_url) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		movieID, movieInput["title"], movieInput["description"], movieInput["year"], movieInput["rating"],
		movieInput["duration"], movieInput["genre"], movieInput["director"], movieInput["poster_url"],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert movie: %v", err)
	}

	// Insert actors
	for _, actor := range actorsInput {
		actorMap := actor.(map[string]interface{})
		actorID := uuid.New().String()
		_, err := database.DB.Exec(`
            INSERT INTO actors (id, name, birth_date, nationality, biography, profile_url) 
            VALUES (?, ?, ?, ?, ?, ?)`,
			actorID, actorMap["name"], actorMap["birth_date"], actorMap["nationality"], actorMap["biography"], actorMap["profile_url"],
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert actor: %v", err)
		}

		// Link actor to movie
		_, err = database.DB.Exec(`
            INSERT INTO movie_actors (movie_id, actor_id) 
            VALUES (?, ?)`,
			movieID, actorID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to link actor to movie: %v", err)
		}
	}

	// Insert reviews
	for _, review := range reviewsInput {
		reviewMap := review.(map[string]interface{})
		reviewID := uuid.New().String()
		_, err := database.DB.Exec(`
            INSERT INTO reviews (id, movie_id, user_name, rating, comment) 
            VALUES (?, ?, ?, ?, ?)`,
			reviewID, movieID, reviewMap["user_name"], reviewMap["rating"], reviewMap["comment"],
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert review: %v", err)
		}
	}

	return getMovieByID(movieID)
}

func UpdateMovie(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required")
	}

	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("input is required")
	}

	if director, ok := input["director"].(string); ok && director != "" {
		if _, err := database.EnsureDirector(director); err != nil {
			return nil, fmt.Errorf("failed to ensure director: %v", err)
		}
	}

	_, err := database.DB.Exec(`
		UPDATE movies 
		SET title = ?, description = ?, year = ?, rating = ?, duration = ?, 
		    genre = ?, director = ?, poster_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		input["title"].(string),
		input["description"],
		input["year"].(int),
		input["rating"].(float64),
		input["duration"].(int),
		input["genre"],
		input["director"],
		input["poster_url"],
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update movie: %v", err)
	}

	return getMovieByID(id)
}

func DeleteMovie(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required")
	}

	// Delete related records first
	_, err := database.DB.Exec("DELETE FROM movie_actors WHERE movie_id = ?", id)
	if err != nil {
		return false, fmt.Errorf("failed to delete movie actors: %v", err)
	}

	_, err = database.DB.Exec("DELETE FROM reviews WHERE movie_id = ?", id)
	if err != nil {
		return false, fmt.Errorf("failed to delete reviews: %v", err)
	}

	// Delete the movie
	result, err := database.DB.Exec("DELETE FROM movies WHERE id = ?", id)
	if err != nil {
		return false, fmt.Errorf("failed to delete movie: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}

func CreateReview(p graphql.ResolveParams) (interface{}, error) {
	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("input is required")
	}

	// Validate rating
	rating := input["rating"].(int)
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	id := uuid.New().String()

	_, err := database.DB.Exec(`
		INSERT INTO reviews (id, movie_id, user_name, rating, comment)
		VALUES (?, ?, ?, ?, ?)`,
		id,
		input["movie_id"].(string),
		input["user_name"].(string),
		rating,
		input["comment"],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %v", err)
	}

	return getReviewByID(id)
}

func getMovieByID(id string) (*models.Movie, error) {
	var movie models.Movie
	err := database.DB.QueryRow(`
		SELECT id, title, description, year, rating, duration, genre, director, poster_url, created_at, updated_at
		FROM movies WHERE id = ?`, id).Scan(
		&movie.ID, &movie.Title, &movie.Description, &movie.Year, &movie.Rating,
		&movie.Duration, &movie.Genre, &movie.Director, &movie.PosterURL,
		&movie.CreatedAt, &movie.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("movie not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query movie: %v", err)
	}
	return &movie, nil
}

func getActorsForMovie(movieID string) ([]models.Actor, error) {
	rows, err := database.DB.Query(`
		SELECT a.id, a.name, a.birth_date, a.nationality, a.biography, a.profile_url
		FROM actors a
		INNER JOIN movie_actors ma ON a.id = ma.actor_id
		WHERE ma.movie_id = ?`, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []models.Actor
	for rows.Next() {
		var actor models.Actor
		err := rows.Scan(&actor.ID, &actor.Name, &actor.BirthDate, &actor.Nationality, &actor.Biography, &actor.ProfileURL)
		if err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}
	return actors, nil
}

func getReviewsForMovie(movieID string) ([]models.Review, error) {
	rows, err := database.DB.Query(`
		SELECT id, movie_id, user_name, rating, comment, created_at
		FROM reviews WHERE movie_id = ? ORDER BY created_at DESC`, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		err := rows.Scan(&review.ID, &review.MovieID, &review.UserName, &review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func getReviewByID(id string) (*models.Review, error) {
	var review models.Review
	err := database.DB.QueryRow(`
		SELECT id, movie_id, user_name, rating, comment, created_at
		FROM reviews WHERE id = ?`, id).Scan(
		&review.ID, &review.MovieID, &review.UserName, &review.Rating, &review.Comment, &review.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review with id %s not found", id)
		}
		return nil, fmt.Errorf("error retrieving review: %v", err)
	}
	return &review, nil
}

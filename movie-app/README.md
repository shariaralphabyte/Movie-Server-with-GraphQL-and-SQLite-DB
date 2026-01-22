# Movie App (GraphQL + Go + SQLite)

This project is a Go GraphQL API backed by SQLite.

## Run

```bash
go test ./...
PORT=8081 go run ./cmd/server
```

- GraphQL endpoint: `http://localhost:8081/graphql`
- GraphiQL UI: `http://localhost:8081/graphql`

> Note: `PORT` defaults to `8080`. If `8080` is already in use, run with another port (example above).

## Database

- SQLite file: `movies.db`
- On startup the app runs `database.InitDatabase()` which:
  - Creates tables (movies, actors, movie_actors, reviews, directors)
  - Seeds data (movies + some actors/reviews/links)

## Core Types

### `Movie`
Includes related data:
- `actors`: loaded via `movie_actors` + `actors`
- `reviews`: loaded from `reviews`

## Queries

### Get all movies

```graphql
query GetAllMovies {
  movies {
    movies {
      id
      title
      year
      rating
      genre
      director
      poster_url
    }
  }
}
```

### Get movie by id (with actors + reviews)

```graphql
query GetMovie($id: ID!) {
  movie(id: $id) {
    id
    title
    description
    year
    rating
    duration
    genre
    director
    poster_url
    created_at
    updated_at
    actors {
      id
      name
      birth_date
      nationality
      biography
      profile_url
    }
    reviews {
      id
      movie_id
      user_name
      rating
      comment
      created_at
    }
  }
}
```

Variables:

```json
{ "id": "1" }
```

### List movies (pagination + filter)

```graphql
query ListMovies($page: Int, $limit: Int, $filter: MovieFilter) {
  movies(page: $page, limit: $limit, filter: $filter) {
    pagination {
      page
      limit
      total
      total_pages
    }
    movies {
      id
      title
      year
      rating
      genre
      director
      poster_url
    }
  }
}
```

Variables example:

```json
{
  "page": 1,
  "limit": 10,
  "filter": {
    "genre": "Drama",
    "min_year": 1990,
    "min_rating": 7.5,
    "search": "dark"
  }
}
```

### Search movies

```graphql
query Search($query: String!, $page: Int, $limit: Int) {
  searchMovies(query: $query, page: $page, limit: $limit) {
    pagination { page limit total total_pages }
    movies { id title year rating }
  }
}
```

## Mutations

### Create movie

```graphql
mutation CreateMovie($input: MovieInput!) {
  createMovie(input: $input) {
    id
    title
    director
  }
}
```

Variables:

```json
{
  "input": {
    "title": "My Movie",
    "description": "A test movie",
    "year": 2024,
    "rating": 8.2,
    "duration": 120,
    "genre": "Action",
    "director": "Some Director",
    "poster_url": "https://example.com/poster.jpg"
  }
}
```

### Update movie

```graphql
mutation UpdateMovie($id: ID!, $input: MovieInput!) {
  updateMovie(id: $id, input: $input) {
    id
    title
    rating
    updated_at
  }
}
```

### Delete movie

```graphql
mutation DeleteMovie($id: ID!) {
  deleteMovie(id: $id)
}
```

### Create review (standalone)

```graphql
mutation CreateReview($input: ReviewInput!) {
  createReview(input: $input) {
    id
    movie_id
    user_name
    rating
    comment
    created_at
  }
}
```

Variables:

```json
{
  "input": {
    "movie_id": "1",
    "user_name": "alice",
    "rating": 5,
    "comment": "Great movie"
  }
}
```

### Create movie WITH details (movie + actors + reviews)

Use this when you want to insert into **multiple tables** in one request.

```graphql
mutation CreateMovieWithDetails($input: MovieWithDetailsInput!) {
  createMovieWithDetails(input: $input) {
    id
    title
    director
    actors { id name }
    reviews { id user_name rating comment }
  }
}
```

Variables:

```json
{
  "input": {
    "movie": {
      "title": "Movie With Details",
      "description": "Inserted with actors + reviews",
      "year": 2025,
      "rating": 8.4,
      "duration": 140,
      "genre": "Drama",
      "director": "Director X",
      "poster_url": "https://example.com/movie-with-details.jpg"
    },
    "actors": [
      {
        "name": "Actor One",
        "birth_date": "1980-01-01",
        "nationality": "US",
        "biography": "Bio...",
        "profile_url": "https://example.com/actor-one"
      },
      { "name": "Actor Two" }
    ],
    "reviews": [
      { "user_name": "john", "rating": 5, "comment": "Amazing" },
      { "user_name": "sara", "rating": 4 }
    ]
  }
}
```

## Notes

- The authoritative GraphQL schema is defined in Go (`internal/resolvers/resolvers.go`) via `github.com/graphql-go/graphql`.
- `internal/schema/schema.graphql` can be used as a reference, but the server does not load it at runtime.

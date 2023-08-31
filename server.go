package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	omdb "./internal/types" // Update this to your actual package path
)

type omdbServer struct {
	omdbAPIKey string
}

func (s *omdbServer) GetMovieByID(ctx context.Context, req *omdb.GetMovieByIDRequest) (*omdb.GetMovieByIDResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "ID cannot be empty")
	}

	// Make an HTTP request to the OMDb API
	client := resty.New()
	resp, err := client.R().SetQueryParams(map[string]string{
		"apikey": s.omdbAPIKey,
		"i":      req.Id,
	}).Get("http://www.omdbapi.com/")
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to fetch movie data from OMDb API")
	}

	var omdbResponse map[string]interface{}
	err = resp.UnmarshalJSON(&omdbResponse)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to parse OMDb API response")
	}

	if omdbResponse["Response"] == "False" {
		return nil, status.Error(codes.NotFound, "Movie not found")
	}

	movie := &omdb.GetMovieByIDResponse{
		id:         req.Id,
		title:      omdbResponse["Title"].(string),
		year:       omdbResponse["Year"].(string),
		rated:      omdbResponse["Rated"].(string),
		genre:      omdbResponse["Genre"].(string),
		plot:       omdbResponse["Plot"].(string),
		director:   omdbResponse["Director"].(string),
		actors:     strings.Split(omdbResponse["Actors"].(string), ", "),
		language:   omdbResponse["Language"].(string),
		country:    omdbResponse["Country"].(string),
		type:       omdbResponse["Type"].(string),
		poster_url:  omdbResponse["Poster"].(string),
	}

	return movie, nil
}

func (s *omdbServer) SearchMovies(ctx context.Context, req *omdb.SearchMoviesRequest) (*omdb.SearchMoviesResponse, error) {
	if len(req.Query) < 3 {
		return nil, status.Error(codes.InvalidArgument, "Query must have at least 3 characters")
	}

	// Make an HTTP request to the OMDb API
	client := resty.New()
	resp, err := client.R().SetQueryParams(map[string]string{
		"apikey": s.omdbAPIKey,
		"s":      req.Query,
		"type":   req.Type,
		"page":   fmt.Sprintf("%d", req.Page),
	}).Get("http://www.omdbapi.com/")
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to fetch movie data from OMDb API")
	}

	var omdbResponse map[string]interface{}
	err = resp.UnmarshalJSON(&omdbResponse)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to parse OMDb API response")
	}

	if omdbResponse["Response"] == "False" {
		return nil, status.Error(codes.NotFound, "No movies found")
	}

	movies := make([]*omdb.MovieResult, 0)
	for _, item := range omdbResponse["Search"].([]interface{}) {
		movieItem := item.(map[string]interface{})
		movies = append(movies, &omdb.MovieResult{
			id:        movieItem["imdbID"].(string),
			title:     movieItem["Title"].(string),
			year:      movieItem["Year"].(string),
			type:      movieItem["Type"].(string),
			poster_url: movieItem["Poster"].(string),
		})
	}

	totalResults := uint64(omdbResponse["totalResults"].(float64))

	return &omdb.SearchMoviesResponse{
		movies:       movies,
		total_results: totalResults,
	}, nil
}

func main() {
	// Replace with your actual OMDb API key
	omdbAPIKey := "faf7e5bb"

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	omdb.RegisterOMDBServiceServer(server, &omdbServer{
		omdbAPIKey: omdbAPIKey,
	})

	log.Println("Starting gRPC server on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
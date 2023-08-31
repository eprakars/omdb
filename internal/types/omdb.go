syntax = "proto3";
package omdb;
option go_package = "internal/types"; // customize based on your structure
service OMDBService {
rpc GetMovieByID(GetMovieByIDRequest) returns (GetMovieByIDResponse) {}

Building a Go server with gRPC for querying the OMDb API - Estimation 3 hours. 2

rpc SearchMovies(SearchMoviesRequest) returns (SearchMoviesResponse) {}
}
message GetMovieByIDRequest {
string id = 1;
}
message GetMovieByIDResponse {
string id = 1;
string title = 2;
string year = 3;
string rated = 4;
string genre = 5;
string plot = 6;
string director = 7;
repeated string actors = 8;
string language = 9;
string country = 10;
string type = 11;
string poster_url = 12;
}
message SearchMoviesRequest {
string query = 1;
string type = 2;
uint64 page = 3;
}
message SearchMoviesResponse {
repeated MovieResult movies = 1;
uint64 total_results = 2;
}
message MovieResult {
string id = 1;
string title = 2;
string year = 3;
string type = 4;
string poster_url = 5;
}
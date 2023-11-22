package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/* Khai báo cấu trúc dữ liệu Movie == tưng tư như định nghĩa các trường trong db */
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

/* Khai báo biến movies có kiểu slice chứa các phần tử có trong Movie */
var movies []Movie

/* Get Movies
 * http.ResponseWriter: đại diện cho phản hồi từ http mà server gửi về client
 * *http.Request : yêu cầu từ client đến server
 * NewEncoder : tạo một đối tượng json.Encoder mới
 * json.Encoder: mã hóa (encode) dữ liệu thành định dạng json
 */
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

/*
 * mux.Vars(r): trích xuất các biến đường dẫn từ routes
 *
 */
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			/* giải thích
			append : thêm các phần tử vào 1 slice và trả về slice mới chứa các phần tử đã thêm
			movies[:index] : là slice con của movies bắt đầu từ 0 đến index-1
			movies[index + 1:]: ... bắt đầu từ index + 1 đến hết
			... :  toán tử ellipsis mở rộng slice thàng các phần tử riêng lẻ
			*/
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

/*
 * Decode : giải mã decode từ định dạng json thành các đối tưởng Go tương ứng.
 */
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(101))
	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	// Set json content type
	w.Header().Set("Content-Type", "application/json")
	// params
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{
		ID:       "1",
		Isbn:     "438227",
		Title:    "Movie One",
		Director: &Director{Firstname: "Jone", Lastname: "Doe"},
	})
	movies = append(movies, Movie{
		ID:       "2",
		Isbn:     "438228",
		Title:    "Movie Two",
		Director: &Director{Firstname: "Phan", Lastname: "Quoc"},
	})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

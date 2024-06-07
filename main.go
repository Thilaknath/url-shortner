package main

import (
    "context"
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/gorilla/mux"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
    rdb = redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Use the appropriate address if running Redis elsewhere
    })
}

type URLRequest struct {
    LongURL string `json:"long_url"`
}

type URLResponse struct {
    ShortURL string `json:"short_url"`
}

func generateShortURL() string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    const length = 8
    rand.Seed(time.Now().UnixNano())
    shortURL := make([]byte, length)
    for i := range shortURL {
        shortURL[i] = charset[rand.Intn(len(charset))]
    }
    return string(shortURL)
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
    var request URLRequest
    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    shortURL := generateShortURL()
    err = rdb.Set(ctx, shortURL, request.LongURL, 0).Err()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    response := URLResponse{ShortURL: shortURL}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    shortURL := vars["shortURL"]
    longURL, err := rdb.Get(ctx, shortURL).Result()
    if err == redis.Nil {
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/shorten", shortenURL).Methods("POST")
    r.HandleFunc("/{shortURL}", redirectURL).Methods("GET")
    http.Handle("/", r)
    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", nil)
}

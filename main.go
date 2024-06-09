package main

import (
    "context"
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
)

var ctx = context.Background()
var rdb *redis.Client
var log = logrus.New()

func init() {
    redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_ENDPOINT"), os.Getenv("REDIS_PORT"))
    rdb = redis.NewClient(&redis.Options{
        Addr: redisAddr,
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
        log.WithError(err).Error("Failed to decode request body")
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    shortURL := generateShortURL()
    log.WithFields(logrus.Fields{
        "shortURL": shortURL,
        "longURL":  request.LongURL,
    }).Info("Generated short URL")
    err = rdb.Set(ctx, shortURL, request.LongURL, 0).Err()
    if err != nil {
        log.WithError(err).Error("Failed to save URL to Redis")
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
        log.WithField("shortURL", shortURL).Warn("Short URL not found")
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    } else if err != nil {
        log.WithError(err).Error("Failed to get URL from Redis")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.WithFields(logrus.Fields{
        "shortURL": shortURL,
        "longURL":  longURL,
    }).Info("Redirecting to long URL")
    http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
    log.SetFormatter(&logrus.JSONFormatter{})
    log.SetLevel(logrus.InfoLevel)

    r := mux.NewRouter()
    r.HandleFunc("/shorten", shortenURL).Methods("POST")
    r.HandleFunc("/{shortURL}", redirectURL).Methods("GET")
    http.Handle("/", r)
    log.Info("Server is running on port 8080")
    http.ListenAndServe(":8080", nil)
}

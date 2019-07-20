package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
)

// Reply structure
type Reply struct {
	OriginalLink string `json:"original_link,omitempty"`
	ShortLink    string `json:"short_link,omitempty"`
}

// Request structure
type Request struct {
	OriginalLink string `json:"url,omitempty"`
}

// RedirectEndpoint Redirect to original url link
func RedirectEndpoint(w http.ResponseWriter, req *http.Request) {

	shortURI := "http://" + req.Host + req.RequestURI
	originalLink, err := redisdb.Get(shortURI).Result()
	if err == redis.Nil {
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		panic(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		http.Redirect(w, req, originalLink, 301)
	}

}

// CreateShortLinkEndpoint create short url with hash, or return existing one from Redis
func CreateShortLinkEndpoint(w http.ResponseWriter, req *http.Request) {
	var url Request
	_ = json.NewDecoder(req.Body).Decode(&url)

	shortLink, err := redisdb.Get(url.OriginalLink).Result()
	if err == redis.Nil {
		// create hash
		hd := hashids.NewData()
		h, _ := hashids.NewWithData(hd)
		now := time.Now()

		nowTime := int(now.Unix())
		ID, _ := h.Encode([]int{nowTime})

		shortLink = "http://localhost/" + ID
		err = redisdb.Set(shortLink, url.OriginalLink, 0).Err()
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	reply := Reply{url.OriginalLink, shortLink}

	js, err := json.Marshal(reply)
        
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

var redisdb *redis.Client

func initRedis() {
	redisdb = redis.NewClient(&redis.Options{
		Addr:         "redis:6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
}

func main() {
	router := mux.NewRouter()
	initRedis()

	router.HandleFunc("/shorten", CreateShortLinkEndpoint).Methods("POST")
	router.HandleFunc("/{id}", RedirectEndpoint).Methods("GET")

	log.Fatal(http.ListenAndServe(":80", router))
}

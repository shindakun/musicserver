package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func getRandFile() string {
	d, err := os.ReadDir("./music")
	if err != nil {
		panic(err)
	}

	n := rand.Intn(len(d) + 1)

	s := "http://localhost:8080/music/" + d[n].Name()

	return s
}

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		fs.ServeHTTP(w, r)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)

	f := getRandFile()

	w.Write([]byte(f))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("musicserver")

	fs := http.FileServer(http.Dir("./music"))
	t := cors(fs)
	http.Handle("/music/", http.StripPrefix("/music", t))

	http.HandleFunc("/", handle)
	http.ListenAndServe(":8080", nil)
}

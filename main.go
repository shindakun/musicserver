package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/shindakun/envy"
)

var client *twitch.Client

func getRandFile() (string, error) {
	d, err := os.ReadDir("./music")
	if err != nil {
		panic(err)
	}

	n := rand.Intn(len(d)+1) - 1

	if !d[n].IsDir() {
		s := "http://localhost:8080/music/" + d[n].Name()

		name := strings.ReplaceAll(d[n].Name(), ".", " ")

		log.Println(name)

		client.Say("shindakun", "Now playing: "+name)

		return s, nil
	}
	return "", fmt.Errorf("is dir")
}

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		fs.ServeHTTP(w, r)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.RequestURI == "/" {
		w.WriteHeader(http.StatusAccepted)
		f, err := getRandFile()
		if err != nil {
			panic(err)
		}
		w.Write([]byte(f))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}

}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("musicserver")

	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	go func() {
		oauth, err := envy.Get("TWITCHAUTH")
		if err != nil {
			panic(err)
		}

		c := twitch.NewClient("shinbot", oauth)

		client = c

		err = c.Connect()
		if err != nil {
			panic(err)
		}

		c.Join("shindakun")

		c.OnUserJoinMessage(func(message twitch.UserJoinMessage) {
			fmt.Println(message)
			c.Say("shindakun", "Welcome!")
		})
	}()

	fs := http.FileServer(http.Dir("./music"))
	t := cors(fs)
	http.Handle("/music/", http.StripPrefix("/music", t))

	http.HandleFunc("/", handle)
	http.ListenAndServe(":8080", nil)
}

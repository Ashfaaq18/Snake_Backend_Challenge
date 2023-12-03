package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type state struct {
	GameID string `json:"gameId"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Score  int    `json:"score"`
	Fruit  fruit  `json:"fruit"`
	Snake  snake  `json:"snake"`
}

type fruit struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type snake struct {
	X    int `json:"x"`    // horizontal pos
	Y    int `json:"y"`    // vertical pos
	VelX int `json:"velX"` // X velocity of the snake (-1, 0, 1) where -1 is left, 1 is right
	VelY int `json:"velY"` // Y velocity of the snake (-1, 0, 1) where -1 is up, 1 is down
}

type velocity struct {
	VelX int `json:"velX"`
	VelY int `json:"velY"`
}

type gameStates struct {
	state
	Ticks []velocity `json:"ticks"`
}

// random position for fruit
func randFruitPosition(width, height int) fruit {
	return fruit{
		X: rand.Intn(width-1) + 1,
		Y: rand.Intn(height-1) + 1,
	}
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else {
		fmt.Printf("GET working %d\n", rand.Intn(20-1)+1)
		w.WriteHeader(http.StatusOK) //200
		w.Header().Set("Content-Type", "application/text")
		fmt.Printf("w value: %s\n", r.URL.Query().Get("w"))
		w.Write([]byte(r.URL.Query().Get("w")))
	}
}

func validateGameHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/new", newGameHandler)
	http.HandleFunc("/validate", validateGameHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

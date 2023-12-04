package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
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
	RecvState state      `json:"recvState"`
	Ticks     []velocity `json:"ticks"`
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
		w.WriteHeader(http.StatusMethodNotAllowed) //Response 405
		return
	} else {
		w.WriteHeader(http.StatusOK) //Response 200
		w.Header().Set("Content-Type", "application/text")

		//send random fruit position as JSON marshalled data back to frontend
		// string to int
		reqWidth, err := strconv.Atoi(r.URL.Query().Get("w"))
		if err != nil {
			// ... handle error
			http.Error(w, err.Error(), http.StatusBadRequest) //Response 405, (4xx: Client Error - The request contains bad syntax or cannot be fulfilled)
			return
		}
		reqHeight, err := strconv.Atoi(r.URL.Query().Get("h"))
		if err != nil {
			// ... handle error
			http.Error(w, err.Error(), http.StatusBadRequest) //Response 405
			return
		}

		//initialized state for new game
		state_marshalled, err := json.Marshal(
			state{
				GameID: "001",
				Width:  reqWidth,
				Height: reqHeight,
				Score:  0,
				Fruit:  randFruitPosition(reqWidth, reqHeight), //randomized fruit position
				Snake: snake{
					X:    0,
					Y:    0,
					VelX: 1,
					VelY: 0,
				},
			})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) //Response 500, (5xx: Server Error - The server failed to fulfill an apparently valid request)
			return
		}
		w.Write(state_marshalled)
		return
	}
}

func validateGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else {
		w.WriteHeader(http.StatusOK) //Response 200
		fmt.Printf("validate function called\n")

		var gs gameStates

		decoder := json.NewDecoder(r.Body)
		fmt.Println("response Body:", r.Body)
		err := decoder.Decode(&gs)

		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", gs.Ticks[0].VelX)

	}
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/new", newGameHandler)
	http.HandleFunc("/validate", validateGameHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

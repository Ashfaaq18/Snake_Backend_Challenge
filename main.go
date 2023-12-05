package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
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
		X: (rand.Intn(width-1) + 1) / 16, //1-96
		Y: (rand.Intn(height-1) + 1) / 16,
	}
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed) //Response 405
		return
	} else {
		w.WriteHeader(http.StatusOK) //Response 200
		w.Header().Set("Content-Type", "application/text")

		//send random fruit position as JSON marshalled data back as a response
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

// validateState validates state for incorrect / missed data
func validateState(gs *gameStates) (validationErrors []string) {
	//validationErrors = append(validationErrors, "validation state errors: ")
	if gs.RecvState.GameID == "" {
		validationErrors = append(validationErrors, "GameID not specified.")
	}

	if gs.RecvState.Width <= 0 || gs.RecvState.Height <= 0 {
		validationErrors = append(validationErrors, "Game board has incorrect size.")
	}

	if gs.RecvState.Fruit.X < 0 || gs.RecvState.Fruit.X >= gs.RecvState.Width ||
		gs.RecvState.Fruit.Y < 0 || gs.RecvState.Fruit.Y >= gs.RecvState.Height {
		validationErrors = append(validationErrors, "Fruit has incorrect position.")
	}

	if gs.RecvState.Snake.X < 0 || gs.RecvState.Snake.X >= gs.RecvState.Width ||
		gs.RecvState.Snake.Y < 0 || gs.RecvState.Snake.Y >= gs.RecvState.Width {
		validationErrors = append(validationErrors, "Snake has incorrect initial position ")
	}
	if gs.RecvState.Snake.VelX < -1 || gs.RecvState.Snake.VelX > 1 ||
		gs.RecvState.Snake.VelY < -1 || gs.RecvState.Snake.VelY > 1 ||
		gs.RecvState.Snake.VelX == gs.RecvState.Snake.VelY {
		validationErrors = append(validationErrors, "Snake has incorrect initial velocity")
	}
	if gs.RecvState.Score < 0 {
		validationErrors = append(validationErrors, "Score cannot be negative number.")
	}
	if len(gs.Ticks) == 0 {
		validationErrors = append(validationErrors, "Ticks are not specified.")
	}
	return
}

func validateMoveSet(gs *gameStates) (validationErrors []string) {
	prevX, prevY := gs.RecvState.Snake.X, gs.RecvState.Snake.Y
	prevVelX, prevVelY := -2, -2 // init with non-possible values do indicate we have no prev velocity before 1st move

	grid := 16
	for i := len(gs.Ticks) - 1; i >= 0; i-- {
		currX, currY := prevX-(gs.Ticks[i].VelX*grid), prevY-(gs.Ticks[i].VelY*grid) // current position
		fmt.Printf("currX: %d , currY: %d | prevX: %d , prevY: %d | tick.VelX: %d , tick.VelY: %d \n", currX, currY, prevX, prevY, gs.Ticks[i].VelX, gs.Ticks[i].VelY)

		// check if snake out of game board borders
		if currX < 0 || currX >= gs.RecvState.Width || currY < 0 || currY >= gs.RecvState.Height {
			tempString := "Snake went out of bounds. currX: " + strconv.Itoa(currX) + ", currY: " + strconv.Itoa(currY)
			validationErrors = append(validationErrors, tempString)
		}

		// check if snake made an invalid move (e.g., immediate 180-degree turn not allowed)
		if (-prevVelX == gs.Ticks[i].VelX && gs.Ticks[i].VelX != 0) ||
			(-prevVelY == gs.Ticks[i].VelY && gs.Ticks[i].VelY != 0) ||
			(gs.Ticks[i].VelX == gs.Ticks[i].VelY) {
			validationErrors = append(validationErrors, "Snake made an invalid move.")
		}

		// update prev before the next iteration
		prevX, prevY = currX, currY
		prevVelX, prevVelY = gs.Ticks[i].VelX, gs.Ticks[i].VelY
	}

	return
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

		valid := false
		validationErrors := validateState(&gs) //validate current state
		if len(validationErrors) > 0 {
			http.Error(w, strings.Join(validationErrors, "\n"), http.StatusBadRequest) //error 400
			valid = false
		} else {
			valid = true
		}

		validationErrors = validateMoveSet(&gs) //validate the snake's moveset
		if len(validationErrors) > 0 {
			http.Error(w, strings.Join(validationErrors, "\n"), http.StatusTeapot) //error 400
			valid = false
		} else {
			valid = true
		}

		if valid {
			//increment game score, generate new position for the fruit, send new game state

			state_marshalled, err := json.Marshal(
				gameStates{
					RecvState: state{
						GameID: gs.RecvState.GameID,
						Width:  gs.RecvState.Width,
						Height: gs.RecvState.Height,
						Score:  gs.RecvState.Score + 1,
						Fruit:  randFruitPosition(gs.RecvState.Width, gs.RecvState.Height), //randomized fruit position
						Snake: snake{
							X:    gs.RecvState.Snake.X,
							Y:    gs.RecvState.Snake.X,
							VelX: gs.RecvState.Snake.VelX,
							VelY: gs.RecvState.Snake.VelY,
						},
					},
					Ticks: []velocity{},
				})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError) //Response 500, (5xx: Server Error - The server failed to fulfill an apparently valid request)
				return
			}
			w.Write(state_marshalled)
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

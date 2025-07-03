package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Game represents the state of a tic-tac-toe game
type Game struct {
	Board       [9]string `json:"board"`
	CurrentTurn string    `json:"currentTurn"`
	Winner      string    `json:"winner"`
	GameOver    bool      `json:"gameOver"`
}

// Global game state
var currentGame *Game

// Initialize a new game
func newGame() *Game {
	return &Game{
		Board:       [9]string{"", "", "", "", "", "", "", "", ""},
		CurrentTurn: "○",
		Winner:      "",
		GameOver:    false,
	}
}

// Check if there's a winner
func (g *Game) checkWinner() string {
	// Winning combinations
	winPatterns := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // columns
		{0, 4, 8}, {2, 4, 6}, // diagonals
	}

	for _, pattern := range winPatterns {
		if g.Board[pattern[0]] != "" &&
			g.Board[pattern[0]] == g.Board[pattern[1]] &&
			g.Board[pattern[1]] == g.Board[pattern[2]] {
			return g.Board[pattern[0]]
		}
	}

	// Check for draw
	draw := true
	for _, cell := range g.Board {
		if cell == "" {
			draw = false
			break
		}
	}

	if draw {
		return "draw"
	}

	return ""
}

// Make a move
func (g *Game) makeMove(position int) bool {
	if g.GameOver || position < 0 || position > 8 || g.Board[position] != "" {
		return false
	}

	g.Board[position] = g.CurrentTurn
	winner := g.checkWinner()

	if winner != "" {
		g.Winner = winner
		g.GameOver = true
	} else {
		// Switch turns
		if g.CurrentTurn == "○" {
			g.CurrentTurn = "×"
		} else {
			g.CurrentTurn = "○"
		}
	}

	return true
}

// HTTP handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		currentGame = newGame()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentGame)
}

func makeMoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if currentGame == nil {
		currentGame = newGame()
	}

	var request struct {
		Position int `json:"position"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	success := currentGame.makeMove(request.Position)
	if !success {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentGame)
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentGame = newGame()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentGame)
}

func main() {
	// Initialize game
	currentGame = newGame()

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/game", gameStateHandler)
	http.HandleFunc("/api/move", makeMoveHandler)
	http.HandleFunc("/api/new-game", newGameHandler)

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Starting tic-tac-toe server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/notnil/chess"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type ChessMoveObj struct {
	Start string
	End   string
}

type ClientChessIncomingMsg struct {
	MoveType     string `json:"type"`
	StartSquare  string `json:"start"`
	EndSquare    string `json:"end"`
	PieceInitial string `json:"piece"`
}

type ClientNewBoardState struct {
	MsgType    string
	Board      map[chess.Square]chess.Piece
	ValidMoves []string
}

func getValidMoves(game *chess.Game) []string {
	moveList := []string{}
	for _, v := range game.ValidMoves() {
		move := v.S1().String() + v.S2().String()
		moveList = append(moveList, move)
	}

	return moveList
}

func reader(conn *websocket.Conn) {

	// create new game with UCI notation
	game := chess.NewGame(chess.UseNotation(chess.UCINotation{}))

	// get valid moves and write to it with ClientNewBoardState
	moves := getValidMoves(game)
	conn.WriteJSON(ClientNewBoardState{
		MsgType:    "newBoard",
		Board:      game.Position().Board().SquareMap(),
		ValidMoves: moves,
	})

	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))
		m := ClientChessIncomingMsg{}
		err = json.Unmarshal([]byte(p), &m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(m)
		if m.MoveType == "move" {
			err = game.MoveStr(strings.Join([]string{m.StartSquare, m.EndSquare}, ""))
			if err != nil {
				log.Println(err)
			}
		} else {
			err = game.MoveStr(strings.Join([]string{m.StartSquare, m.EndSquare, m.PieceInitial}, ""))
			if err != nil {
				log.Println(err)
			}
		}

		newBoardState := ClientNewBoardState{
			MsgType:    "newBoard",
			Board:      game.Position().Board().SquareMap(),
			ValidMoves: getValidMoves(game),
		}
		err = conn.WriteJSON(newBoardState)
		if err != nil {
			log.Println(err)
		}

		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	log.Println(err)
		// 	return
		// }

	}
}

func upgradeWs(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	return ws, nil
}

func serveChessWebsocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket endpoint hit")
	conn, err := upgradeWs(w, r)

	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	reader(conn)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is the main endpoint...?")
	})

	http.HandleFunc("/ws", serveChessWebsocket)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

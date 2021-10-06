package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Slayzur02/GoChess/pkg/chess"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func upgradeWs(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	return ws, nil
}

type ClientNewBoardState struct {
	MsgType       string
	BoardAndMoves chess.BoardAndMoves
}

func createNewBoardClientMsg(bm *chess.BoardAndMoves) *ClientNewBoardState {
	return &ClientNewBoardState{
		MsgType:       "newBoard",
		BoardAndMoves: *bm,
	}
}

type ClientChessIncomingMsg struct {
	MoveType     string `json:"type"`
	StartSquare  string `json:"start"`
	EndSquare    string `json:"end"`
	PieceInitial string `json:"piece"`
}

func reader(conn *websocket.Conn) {

	// create new game with UCI notation
	chessManager := chess.NewService()

	bAndM := chessManager.GetMovesAndBoard()

	// get valid moves and write to it with ClientNewBoardState
	conn.WriteJSON(*createNewBoardClientMsg(bAndM))

	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))
		m := chess.BoardMoveMsg{}
		err = json.Unmarshal([]byte(p), &m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(m)
		chessManager.MovePieceOrPromote(&m)

		bAndM = chessManager.GetMovesAndBoard()
		err = conn.WriteJSON(*createNewBoardClientMsg(bAndM))
		if err != nil {
			log.Println(err)
		}
	}
}

func ServeChessWebsocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket endpoint hit")
	conn, err := upgradeWs(w, r)

	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	reader(conn)
}

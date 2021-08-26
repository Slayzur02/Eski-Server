package chess

import (
	"log"
	"strings"

	"github.com/notnil/chess"
)

type ClientNewBoardState struct {
	MsgType    string
	Board      map[chess.Square]chess.Piece
	ValidMoves []string
}

type ClientIncomingMsg struct {
	MoveType     string `json:"type"`
	StartSquare  string `json:"start"`
	EndSquare    string `json:"end"`
	PieceInitial string `json:"piece"`
}

type GenericBoard struct {
	game *chess.Game
}

func (b *GenericBoard) NewGame() *chess.Game {
	g := chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	return g
}

func (b *GenericBoard) GetMovesAndBoard() *ClientNewBoardState {
	moveList := []string{}
	for _, v := range b.game.ValidMoves() {
		move := v.S1().String() + v.S2().String()
		moveList = append(moveList, move)
	}

	return &ClientNewBoardState{
		MsgType:    "newBoard",
		Board:      b.game.Position().Board().SquareMap(),
		ValidMoves: moveList,
	}
}

func (b *GenericBoard) movePieceOrPromote(m *ClientIncomingMsg) {
	if m.MoveType == "move" {
		err := b.game.MoveStr(strings.Join([]string{m.StartSquare, m.EndSquare}, ""))
		if err != nil {
			log.Println(err)
		}
	} else {
		err := b.game.MoveStr(strings.Join([]string{m.StartSquare, m.EndSquare, m.PieceInitial}, ""))
		if err != nil {
			log.Println(err)
		}
	}
}

package chess

import "github.com/notnil/chess"

type BoardMoveMsg struct {
	MoveType     string `json:"type"`
	StartSquare  string `json:"start"`
	EndSquare    string `json:"end"`
	PieceInitial string `json:"piece"`
}

type BoardAndMoves struct {
	Board      map[chess.Square]chess.Piece
	ValidMoves []string
}

type ClientNewBoardState struct {
	MsgType       string
	BoardAndMoves BoardAndMoves
}

type ClientChessIncomingMsg struct {
	MoveType     string `json:"type"`
	StartSquare  string `json:"start"`
	EndSquare    string `json:"end"`
	PieceInitial string `json:"piece"`
}

package chess

import (
	"log"
	"strings"

	"github.com/notnil/chess"
)

type Repository interface {
}

type Service interface {
	GetMovesAndBoard() *BoardAndMoves
	MovePieceOrPromote(*BoardMoveMsg)
}

type service struct {
	game *chess.Game
	// r *Repository
}

func NewService() Service {
	return &service{
		game: chess.NewGame(chess.UseNotation(chess.UCINotation{})),
	}
}

func (s *service) GetMovesAndBoard() *BoardAndMoves {
	moveList := []string{}

	for _, v := range s.game.ValidMoves() {
		move := v.S1().String() + v.S2().String()
		moveList = append(moveList, move)
	}

	return &BoardAndMoves{
		Board:      s.game.Position().Board().SquareMap(),
		ValidMoves: moveList,
	}
}

func (s *service) MovePieceOrPromote(m *BoardMoveMsg) {
	if m.MoveType == "move" {
		err := s.game.MoveStr(strings.Join([]string{m.StartSquare, m.EndSquare}, ""))
		if err != nil {
			log.Println(err)
		}
	} else {
		err := s.game.MoveStr(strings.Join([]string{m.StartSquare, m.EndSquare, m.PieceInitial}, ""))
		if err != nil {
			log.Println(err)
		}
	}
}

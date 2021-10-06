-- +goose Up
-- +goose StatementBegin
CREATE TABLE app_user (
  id    BIGSERIAL PRIMARY KEY,
  email text      NOT NULL,
  username  text NOT NULL,
  hashedPwd  text NOT NULL
);

CREATE TABLE game (
	id BIGSERIAL PRIMARY KEY,
	whiteID bigint references app_user,
	blackID bigint references app_user
);

CREATE TABLE game_moves (
	turn_number int,
	color int,
	startSquare varchar(10),
	endSquare varchar(10),
	promotionPiece varchar(10),
	game_id bigint references game
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP Table app_user;
DROP TABLE game;
DROP TABLE game_moves;
-- +goose StatementEnd

package datasource

import (
	"context"
	"encoding/json"
	"krestikinoliki/internal/app"
	"log"
	"time"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresGameRepo struct {
	conn *pgx.Conn
}

func NewPostgresGameRepo(conn *pgx.Conn) *PostgresGameRepo {
	return &PostgresGameRepo{conn: conn}
}

func (s *PostgresGameRepo)SaveGame(currentgame *app.CurrentGame) error {

	context_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	entity := ToEntity(currentgame)
	

	fieldJSON, err := json.Marshal(entity.Field)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO games (id, field, status, vs_computer)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET field = EXCLUDED.field
	`
	_, err = s.conn.Exec(context_, query, entity.ID, string(fieldJSON), entity.Status, entity.Computer)
	return err
}

func (s *PostgresGameRepo) LoadGame(ID uuid.UUID) (*app.CurrentGame, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, field, status, vs_computer FROM games WHERE id = $1`
	row := s.conn.QueryRow(ctx, query, ID)

	var entity GameEntity
	var fieldJSON string

	err := row.Scan(&entity.ID, &fieldJSON, &entity.Status, &entity.Computer)
	if err != nil {
		return nil, err // можно уточнить: pgx.ErrNoRows → "игра не найдена"
	}

	err = json.Unmarshal([]byte(fieldJSON), &entity.Field)
	if err != nil {
		return nil, err
	}
	return FromEntity(&entity), nil
}


func ConnectDB() *pgx.Conn{
	connStr := "postgres://postgres:161902@localhost:5432/TicTacToe" // спрятать бы куда?
	context_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := pgx.Connect(context_, connStr)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return conn
}

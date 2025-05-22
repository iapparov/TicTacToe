package datasource

import (
	"context"
	"encoding/json"
	"krestikinoliki/internal/app"
	"log"
	"os"
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

func (s *PostgresGameRepo) CurrentGame(Userid string) []string{
	context_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	
	query := `
	SELECT id 
	From Games
	Where playerx = $1 OR playero = $1
	`

	rows, err := s.conn.Query(context_, query, Userid)
    if err != nil {
		log.Print("failed to execute query: %w", err)
        return nil
    }
    defer rows.Close()
	games := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err!=nil{
			return nil
		}
		games = append(games, id)
	}
	if rows.Err() != nil {
		return nil
	}

	return games
}

func (s *PostgresGameRepo) GetGames() []string{
	context_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	
	query := `
	SELECT id 
	From Games
	Where Status = 0
	`

	rows, err := s.conn.Query(context_, query)
	if err != nil {
		return nil
	}
	games := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err!=nil{
			return nil
		}
		games = append(games, id)
	}
	if rows.Err() != nil {
		return nil
	}

	return games
}

func (s *PostgresGameRepo) SaveGame(currentgame *app.CurrentGame) error {

	context_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	entity := ToEntity(currentgame)
	

	fieldJSON, err := json.Marshal(entity.Field)
	if err != nil {
		return err
	}
	
	query := `
	INSERT INTO games (id, field, status, vs_computer, playerx, playero, createdat)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (id) DO UPDATE SET 
		field = EXCLUDED.field,
		status = EXCLUDED.status,
		vs_computer = EXCLUDED.vs_computer,
		playerx = CASE 
					WHEN games.playerx IS NULL OR games.playerx = '00000000-0000-0000-0000-000000000000' 
					THEN EXCLUDED.playerx 
					ELSE games.playerx 
				  END,
		playero = CASE 
					WHEN games.playero IS NULL OR games.playero = '00000000-0000-0000-0000-000000000000' 
					THEN EXCLUDED.playero 
					ELSE games.playero 
				  END
	`
	_, err = s.conn.Exec(context_, query, entity.ID, string(fieldJSON), entity.Status, entity.Computer, entity.PlayerX, entity.PlayerO, entity.CreatedAt)
	return err
}

func (s *PostgresGameRepo) LoadGame(ID uuid.UUID) (*app.CurrentGame, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, field, status, vs_computer, PlayerX, PlayerO, createdat FROM games WHERE id = $1`
	row := s.conn.QueryRow(ctx, query, ID)

	var entity GameEntity
	var fieldJSON string

	err := row.Scan(&entity.ID, &fieldJSON, &entity.Status, &entity.Computer, &entity.PlayerX, &entity.PlayerO, &entity.CreatedAt)
	if err != nil {
		return nil, err // можно уточнить: pgx.ErrNoRows → "игра не найдена"
	}

	err = json.Unmarshal([]byte(fieldJSON), &entity.Field)
	if err != nil {
		return nil, err
	}
	return FromEntity(&entity), nil
}

func (s *PostgresGameRepo) GetEndedGames(uuid string) ([]string){

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `Select distinct id from games where status IN (3,4,5) AND (playerO = $1 OR PlayerX = $1)`
	rows, err := s.conn.Query(ctx, query, uuid)
    if err != nil {
		log.Print("failed to execute query: %w", err)
        return nil
    }
    defer rows.Close()
	games := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err!=nil{
			return nil
		}
		games = append(games, id)
	}
	if rows.Err() != nil {
		log.Print("failed to execute query: %w", err)
		return nil
	}
	
	return games
}

func (s *PostgresGameRepo) GetLeaderBoard(count int) ([]app.LeaderBoard, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	query := `with game_stats as (
				Select playerx as player_uuid,
						sum(CASE WHEN status = '4' THEN 1 ELSE 0 END) as wins,
						sum(CASE WHEN status = '5' THEN 1 ELSE 0 END) as losses,
						sum(CASE WHEN status = '3' THEN 1 ELSE 0 END) as draws
				From games
				GROUP BY playerx
				UNION ALL
					Select playero as player_uuid,
						sum(CASE WHEN status = '5' THEN 1 ELSE 0 END) as wins,
						sum(CASE WHEN status = '4' THEN 1 ELSE 0 END) as losses,
						sum(CASE WHEN status = '3' THEN 1 ELSE 0 END) as draws
				From games
				GROUP BY playero
			),
			win_ratios as (
				Select player_uuid,
					CAST (wins as Float) / coalesce(NULLIF((losses+draws), 0), 1) as win_ratio
				From game_stats
			)
			Select player_uuid, win_ratio, users.login
			From win_ratios
            Join users ON player_uuid::UUID = users.id
			ORDER By win_ratio desc
			LIMIT $1;`

	rows, err := s.conn.Query(ctx, query, count)
	if err != nil {
		log.Println("Bad query request %w", err)
		return nil, err
	}
	defer rows.Close()
	Leaders := make([]app.LeaderBoard, 0)
	for rows.Next() {
		var id string
		var win_ratio float64
		var login string
		if err := rows.Scan(&id, &win_ratio, &login); err != nil{
			log.Println(err)
			return nil, err
		}
		parsed_id, err := uuid.Parse(id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		Leaders = append(Leaders, app.LeaderBoard{
			UUID: parsed_id,
			Win: win_ratio,
			Login: login,
		})
	}

	return Leaders, nil
}

func ConnectDB() *pgx.Conn{
	

	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB_URL is not set")
	}
	context_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := pgx.Connect(context_, connStr)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return conn
}

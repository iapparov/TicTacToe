package datasource

import (
	"context"
	"krestikinoliki/internal/app"
	"time"
	"errors"
	"github.com/jackc/pgx/v5"
)


type PostgresUserRepo struct {
	conn *pgx.Conn
}

func NewPostgresUserRepo(conn *pgx.Conn) *PostgresUserRepo {
	return &PostgresUserRepo{conn: conn}
}


func (s *PostgresUserRepo) Save(user app.User) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (id, login, password)
		VALUES ($1, $2, $3)
		ON CONFLICT (login) DO UPDATE SET password = EXCLUDED.password
	`
	_, err := s.conn.Exec(ctx, query, user.UUID, user.Login, user.Password)
	return err
}

func (s *PostgresUserRepo) FindByLogin(login string) (app.User, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, login, password FROM users WHERE login = $1
	`

	var user app.User
	err := s.conn.QueryRow(ctx, query, login).Scan(&user.UUID, &user.Login, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return app.User{}, errors.New("user not found")
		}
		return app.User{}, err
	}

	return user, nil
}

func (s *PostgresUserRepo) FindByUUID(uuid string) (bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT * FROM users WHERE id = $1
	`

	var user app.User
	err := s.conn.QueryRow(ctx, query, uuid).Scan(&user.UUID, &user.Login, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false
		}
		return false
	}


	return true
}

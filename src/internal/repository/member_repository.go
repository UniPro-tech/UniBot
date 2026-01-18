package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MemberRepository struct {
	db *pgxpool.Pool
}

func NewMemberRepository(db *pgxpool.Pool) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) Create(ctx context.Context, discordUserID string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO members (discord_user_id)
		 VALUES ($1)
		 ON CONFLICT DO NOTHING`,
		discordUserID,
	)
	return err
}

func (r *MemberRepository) Exists(ctx context.Context, discordUserID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM members WHERE discord_user_id = $1
		)`,
		discordUserID,
	).Scan(&exists)

	return exists, err
}

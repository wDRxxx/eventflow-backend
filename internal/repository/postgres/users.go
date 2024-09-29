package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

func (r *repo) InsertUser(ctx context.Context, user *models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	m, err := utils.MapByStructTag(structTag, *user)
	if err != nil {
		return 0, err
	}

	builder := sq.Insert(usersTable).
		SetMap(m).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) User(ctx context.Context, userEmail string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	sql := `
	SELECT u.id, u.email, u.password, coalesce(tg_username, ''), coalesce(s.shop_id, ''), coalesce(s.shop_key, '')
	FROM users u
	LEFT JOIN users_yookassa_settings s ON u.id = s.user_id
	WHERE U.email = $1
`

	var user models.User
	err := r.db.QueryRow(ctx, sql, userEmail).
		Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.TGUsername,
			&user.YookassaSettings.ShopID,
			&user.YookassaSettings.ShopKey,
		)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) UpdateUserTGUsername(ctx context.Context, userID int64, username string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Update(usersTable).
		Set("tg_username", username).
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) UpdateYookassaSettings(ctx context.Context, settings *models.YookassaSettings) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	m, err := utils.MapByStructTag(structTag, *settings)

	builder := sq.Insert(yookassaSettingsTable).
		SetMap(m).
		Suffix(`ON CONFLICT (user_id)
        DO UPDATE SET shop_id = EXCLUDED.shop_id, shop_key = EXCLUDED.shop_key, updated_at = now()`).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

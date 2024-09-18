package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"time"
)

type repo struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

const (
	usersTable            = "users"
	eventsTable           = "events"
	pricesTable           = "prices"
	TicketsTable          = "tickets"
	yookassaSettingsTable = "users_yookassa_settings"
)

func NewPostgresRepo(db *pgxpool.Pool, timeout time.Duration) repository.Repository {
	return &repo{
		db:      db,
		timeout: timeout,
	}
}

func (r *repo) EventByURLTitle(ctx context.Context, urlTitle string) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Select(
		"title",
		"description",
		"capacity",
		"beginning_time",
		"end_time",
		"creator_id",
		"is_public",
		"location",
		"is_free",
		"preview_image",
		"utc_offset",
		"created_at",
		"updated_at",
	).From(eventsTable).
		Where(sq.Eq{"url_title": urlTitle}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sql, args...)
	var event models.Event

	err = row.Scan(
		&event.Title,
		&event.Description,
		&event.Capacity,
		&event.BeginningTime,
		&event.EndTime,
		&event.CreatorID,
		&event.IsPublic,
		&event.Location,
		&event.IsFree,
		&event.PreviewImage,
		&event.UTCOffset,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	builder = sq.Select(
		"id",
		"price",
		"currency",
		"created_at",
		"updated_at",
	).From(pricesTable).
		Where(sq.Eq{"event_id": event.ID})

	sql, args, err = builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var price models.Price
		err = rows.Scan(
			&price.Price,
			&price.Currency,
			&price.CreatedAt,
			&price.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		event.Prices = append(event.Prices, &price)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *repo) InsertEvent(ctx context.Context, event *models.Event) (id int64, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}

		err = tx.Commit(ctx)
	}()

	m := make(map[string]interface{})
	m["title"] = event.Title
	m["url_title"] = event.URLTitle
	m["description"] = event.Description
	m["beginning_time"] = event.BeginningTime
	m["end_time"] = event.EndTime
	m["creator_id"] = event.CreatorID
	m["is_public"] = event.IsPublic
	m["location"] = event.Location
	m["is_free"] = event.IsFree
	m["utc_offset"] = event.UTCOffset
	m["capacity"] = event.Capacity
	m["minimal_age"] = event.MinimalAge

	if event.PreviewImage != "" {
		m["preview_image"] = event.PreviewImage
	}

	builder := sq.Insert(eventsTable).
		SetMap(m).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	builder = sq.Insert(pricesTable).
		Columns("event_id", "price", "currency")

	for _, p := range event.Prices {
		builder = builder.Values(id, p.Price, p.Currency)
	}
	builder.PlaceholderFormat(sq.Dollar)

	sql, args, err = builder.ToSql()
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) InsertUser(ctx context.Context, user *models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	m := make(map[string]interface{})
	m["email"] = user.Email
	m["password"] = user.Password
	m["created_at"] = time.Now()
	m["updated_at"] = time.Now()

	if user.TGUsername != "" {
		m["tg_username"] = user.TGUsername
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

func (r *repo) InsertYookassaSettings(ctx context.Context, settings *models.YookassaSettings) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Insert(yookassaSettingsTable).
		Columns("user_id", "shop_id", "shop_key").
		Values(settings.UserID, settings.ShopID, settings.ShopKey).
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

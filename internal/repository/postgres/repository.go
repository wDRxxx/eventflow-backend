package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

type repo struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

const (
	usersTable            = "users"
	eventsTable           = "events"
	pricesTable           = "prices"
	ticketsTable          = "tickets"
	yookassaSettingsTable = "users_yookassa_settings"

	structTag = "db"
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
		"id",
		"title",
		"description",
		"capacity",
		"beginning_time",
		"end_time",
		"creator_id",
		"is_public",
		"location",
		"is_free",
		"coalesce(preview_image, '') as preview_image",
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
		&event.ID,
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

	if !event.IsFree {
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

	m, err := utils.MapByStructTag(structTag, *event)

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

	if !event.IsFree {
		builder = sq.Insert(pricesTable).
			Columns("event_id", "price", "currency").
			PlaceholderFormat(sq.Dollar)

		for _, p := range event.Prices {
			builder = builder.Values(id, p.Price, p.Currency)
		}

		sql, args, err = builder.ToSql()
		if err != nil {
			return 0, err
		}

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (r *repo) UpdateEvent(ctx context.Context, event *models.Event) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}

		err = tx.Commit(ctx)
	}()

	err = tx.QueryRow(
		ctx,
		"select id from events where url_title = $1",
		event.URLTitle).
		Scan(&event.ID)
	if err != nil {
		return err
	}

	m, err := utils.MapByStructTag(structTag, *event)
	if err != nil {
		return err
	}

	builder := sq.Update(eventsTable).
		SetMap(m).
		Where(sq.Eq{"url_title": event.URLTitle}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if !event.IsFree {
		insertBuilder := sq.Insert(pricesTable).
			Columns("event_id", "price", "currency").
			PlaceholderFormat(sq.Dollar)

		for _, p := range event.Prices {
			if p.Price > 0 {
				insertBuilder = insertBuilder.Values(event.ID, p.Price, p.Currency)

				continue
			}

			_, err = tx.Exec(
				ctx,
				`delete from prices where event_id = $1 and currency = $2`,
				event.ID,
				p.Currency,
			)
			if err != nil {
				return err
			}
		}

		insertBuilder = insertBuilder.Suffix(`ON CONFLICT (event_id, currency)
        DO UPDATE SET price = EXCLUDED.price, updated_at = now()`)

		sql, args, err = insertBuilder.ToSql()
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repo) DeleteEvent(ctx context.Context, urlTitle string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Delete(eventsTable).
		Where(sq.Eq{"url_title": urlTitle}).
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

//

func (r *repo) InsertTicket(ctx context.Context, ticket *models.Ticket) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	m, err := utils.MapByStructTag(structTag, *ticket)

	builder := sq.Insert(ticketsTable).
		SetMap(m).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return "", err
	}

	return ticket.ID, nil
}

func (r *repo) Ticket(ctx context.Context, ticketID string) (*models.Ticket, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Select(
		"id",
		"user_id",
		"event_id",
		"is_used",
		"first_name",
		"last_name",
	).
		From(ticketsTable).
		Where(sq.Eq{"id": ticketID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var ticket models.Ticket
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&ticket.ID,
		&ticket.UserID,
		&ticket.EventID,
		&ticket.IsUsed,
		&ticket.FirstName,
		&ticket.LastName,
	)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

//

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

	builder := sq.Select("id", "email", "password").
		From(usersTable).
		Where(sq.Eq{"email": userEmail}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.db.QueryRow(ctx, sql, args...).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) InsertYookassaSettings(ctx context.Context, settings *models.YookassaSettings) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	m, err := utils.MapByStructTag(structTag, *settings)

	builder := sq.Insert(yookassaSettingsTable).
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

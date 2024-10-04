package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

func (r *repo) Events(ctx context.Context, page int) ([]*models.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Select(
		"title",
		"coalesce(preview_image, '')",
		"url_title",
		"location",
	).
		From(eventsTable).
		Where(sq.Eq{"is_public": true}).
		OrderBy("created_at DESC").
		Limit(10).
		Offset(uint64((page - 1) * 10)).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event

		err = rows.Scan(
			&event.Title,
			&event.PreviewImage,
			&event.URLTitle,
			&event.Location,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
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
		"minimal_age",
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
		&event.MinimalAge,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if !event.IsFree {
		builder := sq.Select(
			"id",
			"price",
			"currency",
			"created_at",
			"updated_at",
		).From(pricesTable).
			Where(sq.Eq{"event_id": event.ID}).
			PlaceholderFormat(sq.Dollar)

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
				&price.ID,
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
	if err != nil {
		return 0, err
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

	if len(event.Prices) == 0 {
		sql = `delete from prices where event_id = $1`
		_, err = tx.Exec(ctx, sql, event.ID)
		if err != nil {
			return err
		}
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

func (r *repo) UserEvents(ctx context.Context, userID int64) ([]*models.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Select(
		"title",
		"coalesce(preview_image, '')",
		"url_title",
		"location",
	).
		From(eventsTable).
		Where(sq.Eq{"creator_id": userID}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event

		err = rows.Scan(
			&event.Title,
			&event.PreviewImage,
			&event.URLTitle,
			&event.Location,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

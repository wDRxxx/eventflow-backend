package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

func (r *repo) InsertTicket(ctx context.Context, ticket *models.Ticket) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}

		err = tx.Commit(ctx)
	}()

	m, err := utils.MapByStructTag(structTag, *ticket)

	builder := sq.Insert(ticketsTable).
		SetMap(m).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return "", err
	}

	sql = `update events set capacity = capacity - 1 where id = $1`
	_, err = tx.Exec(ctx, sql, ticket.EventID)
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

func (r *repo) UserTickets(ctx context.Context, userID int64) ([]*models.Ticket, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	sql := `SELECT t.id, t.event_id, t.is_used, t.first_name, t.last_name,
	e.title, e.beginning_time, e.preview_image FROM tickets t
	LEFT JOIN events e ON t.event_id = e.id
	WHERE user_id = $1`

	rows, err := r.db.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*models.Ticket
	for rows.Next() {
		ticket := models.Ticket{
			Event: &models.Event{},
		}
		err = rows.Scan(
			&ticket.ID,
			&ticket.EventID,
			&ticket.IsUsed,
			&ticket.FirstName,
			&ticket.LastName,
			&ticket.Event.Title,
			&ticket.Event.BeginningTime,
			&ticket.Event.PreviewImage,
		)
		if err != nil {
			return nil, err
		}

		tickets = append(tickets, &ticket)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

package goflextime

import "time"

type Notification struct {
	NotificationID string    `db:"notification_id" csv:"notification_id"`
	UserID         string    `db:"user_id" csv:"user_id"`
	RegisteredAt   time.Time `db:"registered_at" csv:"registered_at"`
	Title          string    `db:"title" csv:"title"`
	Content        string    `db:"content" csv:"content"`
	ReadStatusTyp  string    `db:"read_status_typ" csv:"read_status_typ"`
	OpenedAt       *NullTime `db:"opened_at" csv:"opened_at"`

	CreatedAt time.Time `db:"created_at" csv:"created_at"`
	UpdatedAt time.Time `db:"updated_at" csv:"updated_at"`
	Revision  int64     `db:"revision" csv:"revision"`
}

type NullTime time.Time

func (date *NullTime) UnmarshalCSV(csv string) (err error) {
	if csv == "" {
		return
	}
	t, err := time.Parse(time.RFC3339, csv)

	dt := NullTime(t)
	*date = dt
	return err
}

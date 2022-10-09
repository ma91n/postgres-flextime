package example

import (
	"os"
	"testing"
	"time"

	"github.com/Songmu/flextime"
	"github.com/gocarina/gocsv"
	"github.com/google/go-cmp/cmp"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func TestUpdateAlreadyRead(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name         string
		conditionSQL string
		args         args
		wantErr      bool
		wantCnt      int
		wantPath     string
	}{
		{
			name: "既読に全更新",
			args: args{
				userID: "00001",
			},
			conditionSQL: `TRUNCATE TABLE notification;
INSERT INTO notification (notification_id, user_id, registered_at, title, content, read_status_typ, opened_at,
                          created_at, updated_at, revision)
VALUES
('1', '00001', '2022-04-01 15:30:00+09', 'テスト1', 'テスト本文1', '0', null, '2022-04-01 15:30:00+09', '2022-04-01 15:30:00+09', 1),
('2', '00001', '2022-04-02 16:30:00+09', 'テスト2', 'テスト本文2', '0', null, '2022-04-02 15:30:00+09', '2022-04-02 15:30:00+09', 1),
('3', '22222', '2022-04-03 16:30:00+09', 'テスト3', 'テスト本文3', '0', null, '2022-04-03 15:30:00+09', '2022-04-03 15:30:00+09', 1);
`,
			wantCnt:  2,
			wantPath: "testdata/want1.csv",
		},
		{
			name: "ユーザーが存在しない",
			args: args{
				userID: "99999",
			},
			conditionSQL: `TRUNCATE TABLE notification;
		INSERT INTO notification (notification_id, user_id, registered_at, title, content, read_status_typ, opened_at,
		                         created_at, updated_at, revision)
		VALUES
		('1', '00001', '2022-04-01 15:30:00', 'テスト1', 'テスト本文1', '0', null, '2022-04-01 15:30:00', '2022-04-01 15:30:00', 1),
		('2', '00001', '2022-04-02 16:30:00', 'テスト2', 'テスト本文2', '0', null, '2022-04-02 15:30:00', '2022-04-02 15:30:00', 1),
		('3', '22222', '2022-04-03 16:30:00', 'テスト3', 'テスト本文3', '0', null, '2022-04-03 15:30:00', '2022-04-03 15:30:00', 1);
		`,
			wantCnt:  0,
			wantPath: "testdata/want2.csv",
		},
		{
			name: "すでに全てが既読",
			args: args{
				userID: "00002",
			},
			conditionSQL: `TRUNCATE TABLE notification;
		INSERT INTO notification (notification_id, user_id, registered_at, title, content, read_status_typ, opened_at,
		                         created_at, updated_at, revision)
		VALUES
		('1', '00001', '2022-04-01 15:30:00', 'テスト1', 'テスト本文1', '1', null, '2022-04-01 15:30:00', '2022-04-01 15:30:00', 1),
		('2', '00001', '2022-04-02 16:30:00', 'テスト2', 'テスト本文2', '1', null, '2022-04-02 15:30:00', '2022-04-02 15:30:00', 1),
		('3', '22222', '2022-04-03 16:30:00', 'テスト3', 'テスト本文3', '0', null, '2022-04-03 15:30:00', '2022-04-03 15:30:00', 1);
		`,
			wantCnt:  0,
			wantPath: "testdata/want3.csv",
		},
	}

	conn, err := sqlx.Open("pgx", "postgres://sample:password@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	tx, err := conn.Beginx()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	restore := flextime.Fix(time.Date(2022, time.October, 11, 10, 10, 10, 0, jst))
	defer restore()

	_, err = tx.Exec("INSERT INTO flex_time (fix_time) VALUES (TO_TIMESTAMP('2022-10-11 10:10:10', 'YYYY-MM-DD HH24:MI:SS'));")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, err := tx.Exec("TRUNCATE TABLE flex_time")
		if err != nil {
			t.Error(err)
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// setup
			if _, err := tx.Exec(tt.conditionSQL); err != nil {
				t.Fatal(err)
			}

			// run
			gotCnt, err := UpdateAlreadyRead(tx, tt.args.userID)

			// verify
			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateAlreadyRead() error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotCnt != tt.wantCnt {
				t.Errorf("UpdateAlreadyRead() gotCnt = %v, wantErr %v", gotCnt, tt.wantCnt)
			}
			gotRecords := gotRecords(t, tx)
			wantRecords := wantRecords(t, tt.wantPath)

			// ★updated_atをテスト対象外にしたくない
			//opts := cmpopts.IgnoreFields("updated_at")
			//if diff := cmp.Diff(wantRecords, gotRecords, opts); diff != "" {
			//	t.Errorf("records mismatch (-want +got):\n%s", diff)
			//}

			if diff := cmp.Diff(wantRecords, gotRecords); diff != "" {
				t.Errorf("records mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func gotRecords(t *testing.T, tx *sqlx.Tx) []Notification {
	rows, err := tx.Queryx("SELECT * FROM notification ORDER BY notification_id")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var resp []Notification
	for rows.Next() {
		var row Notification
		if err := rows.StructScan(&row); err != nil {
			t.Fatal(err)
		}
		resp = append(resp, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return resp
}

func wantRecords(t *testing.T, path string) []Notification {
	wantFile, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	var resp []Notification
	if err := gocsv.UnmarshalFile(wantFile, &resp); err != nil {
		t.Fatal(err)
	}

	for i, row := range resp {
		if time.Time(*row.OpenedAt).IsZero() {
			resp[i].OpenedAt = nil
		}
	}

	return resp
}

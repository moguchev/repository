package schema

import "database/sql"

type Item struct {
	ID     int64         `db:"id"`
	Name   string        `db:"name"`
	Price  int64         `db:"price"`
	Tags   []int64       `db:"tags"`
	Length sql.NullInt64 `db:"length"`
	Width  sql.NullInt64 `db:"width"`
	Height sql.NullInt64 `db:"height"`
}

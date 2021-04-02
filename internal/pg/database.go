package pg

import (
	"fmt"

	"github.com/jackc/pgx"
)

type Database struct {
	conn   *pgx.Conn
	tables map[string]Table
}

func New() *Database {
	return &Database{}
}

func Connect(s string) (*Database, error) {
	conf, err := pgx.ParseConnectionString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %v", s)
	}
	conn, err := pgx.Connect(conf)
	// TODO...
	return &Database{
		conn: conn,
	}, nil
}

func (d *Database) Table(name string) (Table, error) {
	if 

}

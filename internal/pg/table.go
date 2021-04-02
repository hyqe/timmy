package pg

import (
	"fmt"
	"net/http"
)

type Table struct {
	Name  string
	db    *Database
	colls map[string]Column
}

func (t *Table) NewRow(columns []string) ([]interface{}, error) {
	row := make([]interface{}, 0, len(columns))
	for _, name := range columns {
		column, ok := t.colls[name]
		if !ok {
			return nil, fmt.Errorf("table '%v' does not contain column '%v'", t.Name, name)
		}
		row = append(row, column.NewValue())
	}
	return row, nil
}

func (t *Table) Refresh() {
	t.colls = make(map[string]Columner)
	for _, c := range t.Sqler.colls() {
		t.colls[c.Name()] = c
	}
}

func (t *Table) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

// Put values in a table. v must be a reference type.
//func (t *Table) Put(v interface{}) error {}

// Get values from a table. v must be a reference type.
//func (t *Table) Get(Where map[string]interface{}, v interface{}) error {}

// delete rows from a table.
//func (t *Table) Delete(Where map[string]interface{}) error {}

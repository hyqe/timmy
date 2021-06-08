package paper_test

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/hyqe/timmy/internal/paper"
)

type Row struct {
	ID uuid.UUID
	FN string
	LN string
}

func TestPage(t *testing.T) {

	page, err := paper.NewPage("page.txt")
	if err != nil {
		t.Fatal(err)
	}

	index := []int64{}

	for _, row := range []Row{
		{
			ID: uuid.New(),
			FN: "abby",
			LN: "walice",
		},
		{
			ID: uuid.New(),
			FN: "Bob",
			LN: "walice",
		},
		{
			ID: uuid.New(),
			FN: "Craig",
			LN: "Gordon",
		},
		{
			ID: uuid.New(),
			FN: "Dan",
			LN: "Shure",
		},
		{
			ID: uuid.New(),
			FN: "Eddy",
			LN: "Martin",
		},
	} {
		key, err := page.Put(row)
		if err != nil {
			t.Fatal(err)
		}
		index = append(index, key)
	}

	rows := make([]Row, 0)
	for _, key := range index {
		var row Row
		err = page.Get(key, &row)
		if err != nil {
			t.Fatal(err)
		}
		rows = append(rows, row)
	}

	log.Printf("%+v", rows)
}

func TestPageWalk(t *testing.T) {
	page, err := paper.NewPage("page.txt")
	if err != nil {
		t.Fatal(err)
	}

	row := Row{}
	err = page.Walk(func(key int64, decoder paper.Decoder) error {
		err := decoder(&row)
		if err != nil {
			return err
		}

		log.Println(key, row)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPageIterate(t *testing.T) {
	page, err := paper.NewPage("page.txt")
	if err != nil {
		t.Fatal(err)
	}

	row := Row{}
	for line := range page.Iterate(context.Background()) {
		err := line.Decoder(&row)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(line.Key, row)
	}
}

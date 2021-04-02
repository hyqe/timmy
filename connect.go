package timmy

import (
	"fmt"
	"net/url"

	"github.com/jakobii/timmy/internal/pg"
)

func Connect(conn string) (Databaser, error) {
	u, err := url.Parse(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conntection string: %v", err)
	}

	switch u.Scheme {
	case "postgresql":
		return pg.Connect(conn)
	default:
		return nil, fmt.Errorf("url scheme did not match any driver")
	}
}

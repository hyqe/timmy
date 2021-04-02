package timmy

type Databaser interface {
	Connect(string) error
	Disconnect() error
	Table(string) (Tabler, error)
}

type Tabler interface {
	Name() string
	Columns() []string
	Insert(values []map[string]interface{}) error
	Select(where map[string]interface{}) ([]map[string]interface{}, error)
	Update(set []map[string]interface{}) error
	Delete(where map[string]interface{}) error
}

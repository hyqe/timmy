package tables

type Tabler interface {
	Insert(Row)
	Select(matcher Matcher, Sort Sorter) Interator
	Delete(matcher Matcher) error
}

type Matcher interface {
	Match() bool
}

type Sorter interface {
	Sort()
}

type Interator interface {
	Next()
	Scan(pointer interface{})
}

type Row interface {
}

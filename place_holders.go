package engine

import "fmt"

type PlaceHolder interface {
	NextToken() string
	Reset() bool
}

type PostgresPlaceholder struct {
	index int
}

func NewPostgresPlaceHolder() *PostgresPlaceholder {
	return &PostgresPlaceholder{
		index: 0,
	}
}

func (p *PostgresPlaceholder) NextToken() string {
	p.index += 1

	return fmt.Sprintf("$%d", p.index)
}

func (p *PostgresPlaceholder) Reset() bool {
	p.index = 0
	return true
}

type SqlitePlaceholder struct{}

func NewSqlitePlaceholder() *SqlitePlaceholder {
	return &SqlitePlaceholder{}
}

func (s *SqlitePlaceholder) NextToken() string {
	return "?"
}

func (s *SqlitePlaceholder) Reset() bool {
	return true
}

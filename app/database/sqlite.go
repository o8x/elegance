package database

import (
	_ "embed"

	"github.com/o8x/jk/v2/puresqlite"

	"github.com/o8x/elegance/app/database/queries"
)

var (
	//go:embed ddl.sql
	ddl   string
	query *queries.Queries
)

var DataFile = "elegance.sqlite"

func init() {
	if err := Init(DataFile); err != nil {
		panic(err)
	}
}

func GetQueries() *queries.Queries {
	return query
}

func Init(filename string) error {
	if err := puresqlite.Init(filename); err != nil {
		return err
	}

	query = queries.New(puresqlite.Get())
	_, err := puresqlite.Get().Exec(ddl)
	return err
}

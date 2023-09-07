package postgres

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	_ "github.com/lib/pq"
)

type Message struct {
	TableName         struct{} `pg:"message"`
	Session_id        string   `pg:"session_id"`
	Frequency         string   `pg:"frequency"`
	Current_timestamp string   `pg:"timestamp"`
}

func NewDBConn() (con *pg.DB) {
	address := fmt.Sprintf("%s:%s", "localhost", "5432")
	options := &pg.Options{
		User:     "postgres",
		Password: "1234",
		Addr:     address,
		Database: "goteam",
		PoolSize: 50,
	}
	con = pg.Connect(options)
	if con == nil {
		log.Fatal("cannot connect to postgres")
	}
	err := createSchema(con)
	if err == nil {
		fmt.Println("Table messages created")
	} else if strings.HasSuffix(err.Error(), "already exists") {
		fmt.Println("DB already exists and wasn't been created")
	} else {
		log.Fatal(err)
	}
	return con
}

func InsertDB(pg *pg.DB, post *Message) error {
	_, err := pg.Model(post).Insert()
	return err
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Message)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

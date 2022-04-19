package main

import (
	"log"

	"codeberg.org/rchan/hmn/business"
	"codeberg.org/rchan/hmn/config"
	"codeberg.org/rchan/hmn/data"
	"codeberg.org/rchan/hmn/model"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(log.Llongfile | log.LUTC)
}

func main() {

	conf := config.LoadConfig("./config.json")

	db := data.OpenDB(conf)

	business := business.NewBusunessLayer(db)

	c, tx, err := business.GetContextFor(1)
	n := &model.Note{}
	n.SetTitle("first Note")
	n.SetContent("this is the first note")
	n.SetParentID(1)
	n.SetIndex(0)

	err = business.Note().AddNote(c, n)
	if err != nil {
		log.Println(err)
	}

	notes, err := business.Note().GetAllNote(c)
	if err != nil {
		log.Println(err)
	}

	for _, x := range notes {
		log.Printf("%v", x)
	}

	tx.Commit()
	log.Println("end")
}

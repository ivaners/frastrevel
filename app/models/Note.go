package models

import (
	"gopkg.in/mgo.v2/bson"
	"myapp/app/db"
	"myapp/app/table"
)

type Note struct{}

func (n *Note) GetNoteAll() (notes []table.Notes) {
	notes = []table.Notes{}
	db.ListByQ(db.Notes, bson.M{}, &notes)
	return notes
}

package controllers

import (
	"myapp/app/models"
)

var note *models.Note

func Init() {
	note = models.NoteS
}

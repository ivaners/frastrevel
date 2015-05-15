package controllers

import (
	"myapp/app/models"
)

var note *models.Note
var session *models.Session
var auth *models.Auth

func Init() {
	note = models.NoteS
	session = models.Session
	auth = models.Auth
}

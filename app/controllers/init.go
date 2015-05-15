package controllers

import (
	"myapp/app/models"
)

var note *models.Note
var session *models.Session
var auth *models.Auth
var user *models.User

func Init() {
	note = models.NoteS
	session = models.SessionS
	auth = models.AuthS
	user = models.Users
}

package models

var NoteS *Note
var SessionS, session *Session
var AuthS *Auth
var Users *User

func Init() {
	NoteS = &Note{}
	SessionS = &Session{}
	session = SessionS
	AuthS = &Auth{}
	Users = &User{}
}

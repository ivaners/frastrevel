package models

var NoteS *Note
var SessionS, session *Session
var AuthS *Auth

func Init() {
	NoteS = &Note{}
	SessionS = &Session{}
	session = SessionS
	AuthS = &Auth{}
}

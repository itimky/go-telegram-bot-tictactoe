package session

type ISessionStorage interface {
	Load(id int) (*session, error)
	Save(session *session) error
}

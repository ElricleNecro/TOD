package commands

type User interface {
	GetPrivateKey() string
	GetUsername() string
}

type Command struct {
	User    User
	Command string
}

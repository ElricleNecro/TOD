package user

type User struct {
	Username   string
	PrivateKey string
}

func (user *User) GetUsername() string {
	return user.Username
}

func (user *User) GetPrivateKey() string {
	return user.PrivateKey
}

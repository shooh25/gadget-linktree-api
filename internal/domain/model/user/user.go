package user

type User struct {
	id      UserId
	profile Profile
}

func NewUser(id UserId, profile Profile) *User {
	return &User{
		id:      id,
		profile: profile,
	}
}

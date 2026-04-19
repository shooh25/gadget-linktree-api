package user

type User struct {
	id UserId
	profile string
}

func NewUser(userId UserId, profile string) *User {
	return &User{
		id: userId,
		profile: profile,
	}
}

package user

type User struct {
	userId   UserId
	googleId string
	profile  Profile
}

func NewUser(userid UserId, googleId string, profile Profile) *User {
	return &User{
		userId: userid,
		googleId: googleId,
		profile: profile,
	}
}

func (u *User) UserId() UserId {
	return u.userId
}

func (u *User) GoogleId() string {
	return u.googleId
}

func (u *User) Profile() Profile {
	return u.profile
}

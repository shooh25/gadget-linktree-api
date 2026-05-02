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

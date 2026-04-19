package user

type UserId struct {
	value string
}

func NewUserId(value string) UserId {
	return UserId{
		value: value,
	}
}

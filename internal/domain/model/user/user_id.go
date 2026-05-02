package user

import "github.com/google/uuid"

type UserId struct {
	value string
}

func NewUserId(value string) UserId {
	return UserId{
		value: value,
	}
}

func GenerateUserId() UserId {
	return UserId{
		value: uuid.New().String(),
	}
}

func (u UserId) Value() string {
	return u.value
}

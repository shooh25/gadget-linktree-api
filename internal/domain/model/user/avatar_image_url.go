package user

type AvatarImageURL struct {
	value string
}

func NewAvatarImageURL(value string) AvatarImageURL {
	return AvatarImageURL{
		value: value,
	}
}

func (a AvatarImageURL) Value() string {
	return a.value
}

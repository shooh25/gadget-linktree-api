package user

type SocialURL struct {
	value string
}

func NewSocialURL(value string) SocialURL {
	return SocialURL{
		value: value,
	}
}

func (s SocialURL) Value() string {
	return s.value
}

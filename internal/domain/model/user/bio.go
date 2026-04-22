package user

type Bio struct {
	value string
}

func NewBio(value string) Bio {
	return Bio{
		value: value,
	}
}

func (b Bio) Value() string {
	return b.value
}

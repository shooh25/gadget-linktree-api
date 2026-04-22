package user

type Profile struct {
	bio         Bio
	displayName string
	avatarImage AvatarImageURL
	socialLinks []SocialLink
}

func NewProfile(displayName string, avatarImage AvatarImageURL, bio Bio, socialLinks []SocialLink) Profile {
	return Profile{
		bio:         bio,
		displayName: displayName,
		avatarImage: avatarImage,
		socialLinks: socialLinks,
	}
}

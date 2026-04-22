package user

type SocialLink struct {
	platform string
	url      SocialURL
}

func NewSocialLink(platform string, url SocialURL) SocialLink {
	return SocialLink{
		platform: platform,
		url:      url,
	}
}

func (s SocialLink) Platform() string {
	return s.platform
}

func (s SocialLink) URL() SocialURL {
	return s.url
}

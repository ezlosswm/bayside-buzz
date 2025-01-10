package domain

type Settings struct {
	PageData                                   PageData
	IsDisabled, IsLoggedIn, IsCreated, IsError bool
}

func NewSettings() Settings {
	return Settings{
		PageData: PageData{
			Title:       "Bayside Buzz",
			SiteName:    "Bayside Buzz",
			Description: "Discover Events in Corozal, Belize - Bayside Buzz",
			Type:        "website",
			Image:       "/assets/images/corozal-sign.jpg",
			URL:         "https://baysidebuzz.com", // r.Host
		},
		IsDisabled: false,
		IsLoggedIn: false,
		IsCreated:  false,
		IsError:    false,
	}
}

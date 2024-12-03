package domain

type PageData struct {
	SiteName, Title, Description, Type, Image, URL string
}

func NewPageData(siteName, title, description, pageType, image, url string) *PageData {
    return &PageData{
        SiteName: siteName,
        Title: title,
        Description: description,
        Type: pageType,
        Image: image,
        URL: url,
    }
}

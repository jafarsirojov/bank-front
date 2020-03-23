package history



type Url string

type History struct {
	url Url
}

func NewHistory(url Url) *History {
	return &History{url: url}
}

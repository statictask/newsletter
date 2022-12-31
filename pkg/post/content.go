package post

import "fmt"

type ContentItem interface {
	GetTitle() string
	GetContent() string
	GetLink() string
}

type Content struct {
	title string
	items []ContentItem
}

func NewContent(title string, items []ContentItem) *Content {
	return &Content{title, items}
}

func (c *Content) Build() (string, error) {
	if len(c.items) == 0 {
		return "", fmt.Errorf("empty list of content items to build")
	}

	body := "<html>\n" +
		"  <head>\n" +
		"    <title>" + c.title + "</title>\n" +
		"  </head>\n" +
		"  <body>\n" +
		"    <h1>" + c.title + "</h1>\n" +
		"    <br>\n"

	for _, i := range c.items {
		body = body +
			"    <hr>\n" +
			"    <a href=\"" + i.GetLink() + "\">\n" +
			"      <h3>" + i.GetTitle() + "</h3>\n" +
			"    </a>\n" +
			"    <br>\n" +
			"    <p>" + i.GetContent() + "</p>\n"
	}

	body = body +
		"    <br><br>\n" +
		"  </body>\n" +
		"</html>\n"

	return body, nil
}

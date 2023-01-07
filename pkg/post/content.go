package post

import (

	"fmt"

	"github.com/statictask/newsletter/pkg/postitem"
)

type ContentBuilder struct {
	title string
	items []*postitem.PostItem
}

func NewContentBuilder(title string, items []*postitem.PostItem) *ContentBuilder {
	return &ContentBuilder{title, items}
}

func (c *ContentBuilder) BuildHTML() (string, error) {
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
			"    <a href=\"" + i.Link + "\">\n" +
			"      <h3>" + i.Title + "</h3>\n" +
			"    </a>\n" +
			"    <br>\n" +
			"    <p>" + i.Content + "</p>\n"
	}

	body = body +
		"    <br><br>\n" +
		"  </body>\n" +
		"</html>\n"

	return body, nil
}

func (c *ContentBuilder) BuildPlainText() (string, error) {
	if len(c.items) == 0 {
		return "", fmt.Errorf("empty list of content items to build")
	}

	body := c.title

	for _, i := range c.items {
		body = body +
			"\n\n" +
			i.Title +
			"\n" +
			i.Content +
			"\n" +
			i.Link +
			"\n"
	}

	return body, nil
}

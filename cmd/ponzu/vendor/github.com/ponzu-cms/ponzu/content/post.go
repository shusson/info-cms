package content

import (
	"fmt"

	"github.com/ponzu-cms/ponzu/management/editor"
	"github.com/ponzu-cms/ponzu/system/item"
	"encoding/json"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"
	"regexp"
)

type Post struct {
	item.Item

	Title string `json:"title"`
	Image string `json:"image"`
	Body  string `json:"body"`
}

// MarshalEditor writes a buffer of html to edit a Post within the CMS
// and implements editor.Editable
func (p *Post) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(p,
		// Take note that the first argument to these Input-like functions
		// is the string version of each Post field, and must follow
		// this pattern for auto-decoding and auto-encoding reasons:
		editor.Field{
			View: editor.Input("Title", p, map[string]string{
				"label":       "Title",
				"type":        "text",
				"placeholder": "Enter the Title here",
			}),
		},
		editor.Field{
			View: editor.Input("Image", p, map[string]string{
				"label":       "Image",
				"type":        "text",
				"placeholder": "Enter the image path here",
			}),
		},
		editor.Field{
			View: editor.Textarea("Body", p, map[string]string{
				"label":       "Body",
				"placeholder": "Enter the Body here",
			}),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to render Post editor view: %s", err.Error())
	}

	return view, nil
}

func (p *Post) String() string {
	return p.Title
}

func (p *Post) generateMarkdown() string {
	unsafe := blackfriday.MarkdownCommon([]byte(p.Body))
	x := bluemonday.UGCPolicy()
	x.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	return string(x.SanitizeBytes(unsafe))
}

func (p *Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	return json.Marshal(&struct {
		Body string `json:"body"`
		*Alias
	}{
		Body: p.generateMarkdown(),
		Alias: (*Alias)(p),
	})
}

func (p *Post) UnmarshalJSON(data []byte) error {
	type Alias Post
	aux := &struct {
		Body string `json:"body"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.Body = aux.Body
	return nil
}

func init() {
	item.Types["Post"] = func() interface{} { return new(Post) }
}

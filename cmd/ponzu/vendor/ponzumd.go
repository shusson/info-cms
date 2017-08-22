package ponzumd

import (
	"bytes"
	"github.com/ponzu-cms/ponzu/management/editor"
)

func Input(fieldName string, p interface{}, attrs map[string]string) []byte {
	html := &bytes.Buffer{}
	e := editor.NewElement("textarea", attrs["label"], fieldName, p, attrs)
	html.Write(editor.DOMElement(e))
	html.WriteString(`WTF`)
	return html.Bytes()
}
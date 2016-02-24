package beego_assets

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/css"
)

func MinifyStylesheet(body string) (string, error) {
	b_result := bytes.NewBuffer([]byte{})
	body_reader := strings.NewReader(body)
	err := css.Minify(minifier, b_result, body_reader, map[string]string{})
	if err != nil {
		return "", errors.New(fmt.Sprintf("Minification error: %v", err))
	}
	return b_result.String(), nil
}

func MinifyJavascript(body string) (string, error) {
	b_result := bytes.NewBuffer([]byte{})
	body_reader := strings.NewReader(body)

	err := js.Minify(minifier, b_result, body_reader, map[string]string{})
	if err != nil {
		return "", errors.New(fmt.Sprintf("Minification error: %v", err))
	}
	return b_result.String(), nil
}
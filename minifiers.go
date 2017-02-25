package beegoAssets

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

var minifier = minify.New()

//MinifyStylesheet - Minify css file
func MinifyStylesheet(body string) (string, error) {
	bResult := bytes.NewBuffer([]byte{})
	bodyReader := strings.NewReader(body)
	err := css.Minify(minifier, bResult, bodyReader, map[string]string{})
	if err != nil {
		return "", fmt.Errorf("Minification error: %v", err)
	}
	return bResult.String(), nil
}

//MinifyJavascript - Minify js file
func MinifyJavascript(body string) (string, error) {
	bResult := bytes.NewBuffer([]byte{})
	bodyReader := strings.NewReader(body)
	err := js.Minify(minifier, bResult, bodyReader, map[string]string{})
	if err != nil {
		return "", fmt.Errorf("Minification error: %v", err)
	}
	return bResult.String(), nil
}

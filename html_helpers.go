package beegoAssets
import (
"fmt"
"html/template"
)

func jsTag(location string) template.HTML {
	return template.HTML(fmt.Sprintf("<script type=text/javascript src=\"%s\"></script>", location))
}
func cssTag(location string) template.HTML {
	return template.HTML(fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\"> ", location))
}

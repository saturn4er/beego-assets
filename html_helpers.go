package beego_assets
import (
"fmt"
"html/template"
)

func js_tag(location string) template.HTML {
	return template.HTML(fmt.Sprintf("<script type=text/javascript src=\"%s\"></script>", location))
}
func css_tag(location string) template.HTML {
	return template.HTML(fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\"> ", location))
}

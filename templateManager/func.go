package templateManager

import (
	"html/template"
)

func templateFunc_html(x string) interface{} {
	return template.HTML(x)
}

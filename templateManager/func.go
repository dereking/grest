package templateManager

import (
	"html/template"
)

func html(x string) interface{} {
	return template.HTML(x)
}

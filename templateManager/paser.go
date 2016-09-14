package templateManager

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/dereking/grest/debug"
)

//var allTemplates map[string]*template.Template
var allTemplates *template.Template

const suffix = ".HTML" //大写
var PthSep string

func init() {
	//allTemplates = make(map[string]*template.Template)
	PthSep = string(os.PathSeparator)
	allTemplates = template.New("")

	allTemplates = allTemplates.Funcs(template.FuncMap{"html": html})

	loadTemplateDir("/", "views")
}

func loadTemplateDir(templateRoot, dir string) {

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("Loading template err:", err)
		return
	}

	for _, f := range fs {
		//文件名
		fn := dir + PthSep + f.Name()

		//模板名字
		tn := templateRoot + strings.ToLower(f.Name())

		if f.IsDir() {
			loadTemplateDir(templateRoot+strings.ToLower(f.Name())+"/", fn)
		} else {

			if strings.HasSuffix(strings.ToUpper(f.Name()), suffix) { //匹配文件

				debug.Debug("Parsing template Name:", tn)
				t := allTemplates.New(tn)

				content, err := ioutil.ReadFile(fn)
				if err != nil {
					debug.Debug("parse", fn, err)
				}
				_, err = t.Parse(string(content))
				if err != nil {
					debug.Debug(" * parse ERROR:", fn, err)
				}
			}
		}
	}
}

func getTemplate(name string) *template.Template {
	//return allTemplates[name]
	return allTemplates.Lookup(name)
}

//Render execute the template.
// if template not found , renturn nil,error
func Render(ctlName, actName string, viewBag map[string]interface{}) ([]byte, error) {

	//tname := fmt.Sprintf("views/%s/%s.html", ctlName, actName)
	fn := strings.ToLower("/" + ctlName + "/" + actName + ".html")

	t := getTemplate(fn)
	if t == nil {
		return nil, errors.New(fmt.Sprintf("template %s not found", fn))
		/*
			debug.Debug("Render HTML:", fn)

			content, err := ioutil.ReadFile(fn)
			if err != nil {
				if os.IsNotExist(err) {
					return nil, errors.New(fmt.Sprintf("template %s-%s not found", ctlName, actName))
				}
			}

			debug.Debug("Render content:", string(content))
			//t = allTemplates.New(tname)

			t, _ = t.Parse(string(content))
			//		allTemplates[tname] = t
		*/
	}
	b := bytes.NewBuffer(make([]byte, 0))
	err := t.Execute(b, viewBag)
	if err != nil {
		return b.Bytes(), err

	} else {
		return b.Bytes(), nil
	}

}

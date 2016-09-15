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

	"github.com/dereking/grest/config"
	"github.com/dereking/grest/debug"
	"github.com/fsnotify/fsnotify"
)

//var allTemplates map[string]*template.Template
var allTemplates *template.Template

var suffix string //大写
var TMPLATE_DIR string

var PthSep string
var watcher *fsnotify.Watcher

//Initialize the templates.
// args
func Initialize() {

	//TMPLATE_DIR : the directory of templates exists. eg: "views"
	//suffix : the ext name of the template file. UPCASE. eg: ".HTML"
	//bMoniteTemplate: bool, Need monite the template file modify, and auto reload it?
	TMPLATE_DIR = config.AppConfig.StringDefault("TemplateDir", "views")
	suffix = config.AppConfig.StringDefault("TemplateExt", ".HTML")
	bMoniteTemplate := config.AppConfig.BoolDefault("AutoReloadTemplate", false)

	PthSep = string(os.PathSeparator)
	allTemplates = template.New("")

	allTemplates = allTemplates.Funcs(template.FuncMap{"html": html})

	if bMoniteTemplate {
		var err error
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(" fsnotify.NewWatcher err", err)
		}

		loadTemplateDir("/", TMPLATE_DIR, true)

		// Process events
		go func() {
			for {
				select {
				case ev := <-watcher.Events:
					debug.Debug("template modified, reload all template.", ev)

					parseTemplate(ev.Name)

				case err := <-watcher.Errors:
					log.Println("error:", err)
				}
			}
		}()
	} else {

		loadTemplateDir("/", TMPLATE_DIR, false)
	}
}

func loadTemplateDir(templateRoot, dir string, bMoniteTemplate bool) {

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("Loading template err:", err)
		return
	}

	if bMoniteTemplate {
		//moniter this dir
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal("watcher.Watch err", TMPLATE_DIR, err)
		}
	}

	for _, f := range fs {
		//文件名
		fn := dir + PthSep + f.Name()

		if f.IsDir() {
			loadTemplateDir(templateRoot+strings.ToLower(f.Name())+"/", fn, bMoniteTemplate)
		} else {
			parseTemplate(fn)
		}
	}
}

func parseTemplate(fn string) {

	if strings.HasSuffix(strings.ToUpper(fn), strings.ToUpper(suffix)) { //匹配文件
		//模板名字
		tn := strings.Replace(fn, TMPLATE_DIR, "", 1)
		tn = strings.ToLower(strings.Replace(tn, "\\", "/", -1))

		debug.Debug("Parsing template Name:", tn, fn)

		t := allTemplates.Lookup(tn)
		if t == nil {
			t = allTemplates.New(tn)
		}

		content, err := ioutil.ReadFile(fn)
		if err != nil {
			debug.Debug("parse", fn, tn, err)
		}

		_, err = t.Parse(string(content))
		if err != nil {
			debug.Debug(" * parse ERROR:", fn, tn, err)
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

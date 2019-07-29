package templateManager

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dereking/grest/config"
	//"github.com/dereking/grest/debug"

	"github.com/dereking/grest/log"
	"go.uber.org/zap"

	"github.com/fsnotify/fsnotify"
)

//var allTemplates map[string]*template.Template
var allTemplates *template.Template

var suffix string //大写
var TMPLATE_DIR string

var LAYOUT_TEMPLATE string //layout模板路径
var LAYOUT_DATA string

var PthSep string
var watcher *fsnotify.Watcher

const (
	LAYOUT_BODY_TAG = "{{ @RenderBody() }}"
)

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

	allTemplates = allTemplates.Funcs(template.FuncMap{"html": templateFunc_html})
	allTemplates = allTemplates.Funcs(template.FuncMap{"fileSize": templateFunc_FileSize})
	allTemplates = allTemplates.Funcs(template.FuncMap{"datetime": templateFunc_DateTime})
	allTemplates = allTemplates.Funcs(template.FuncMap{"add": templateFunc_add})

	//读取layout数据
	layoutFN := fmt.Sprintf("%s%sShared%s_Layout.html", TMPLATE_DIR, PthSep, PthSep)
	content, err := ioutil.ReadFile(layoutFN)
	if err != nil {
		log.Logger().Fatal("load layout file",
			zap.String("layoutFN", layoutFN),
			zap.Error(err))
	} else {
		LAYOUT_DATA = string(content)
	}

	if bMoniteTemplate {
		var err error
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Logger().Error(" fsnotify.NewWatcher err", zap.Error(err))
		}

		loadTemplateDir("/", TMPLATE_DIR, true)

		// Process events
		go func() {
			for {
				select {
				case ev := <-watcher.Events:
					log.Logger().Debug("template modified, reload template:",
						zap.String("templatename", ev.Name))

					parseTemplate(ev.Name)

				case err := <-watcher.Errors:
					log.Logger().Error("error:", zap.Error(err))
				}
			}
		}()
	} else {

		loadTemplateDir("/", TMPLATE_DIR, false)
	}
}

//遍历模板目录，并进行模板编译。如果需要动态监视模板变动，进行目录监控。
func loadTemplateDir(templateRoot, dir string, bMoniteTemplate bool) {

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Logger().Error("Loading template err:", zap.Error(err))
		return
	}

	if bMoniteTemplate {
		//moniter this dir
		err = watcher.Add(dir)
		if err != nil {
			log.Logger().Error("watcher.Watch err",
				zap.String("TMPLATE_DIR", TMPLATE_DIR),
				zap.Error(err))
		}
	}

	for _, f := range fs {
		//文件名
		fn := dir + PthSep + f.Name()

		if f.IsDir() {
			//递归子目录遍历
			loadTemplateDir(templateRoot+strings.ToLower(f.Name())+"/", fn, bMoniteTemplate)
		} else {
			//处理当前模板，编译并记录。
			parseTemplate(fn)
		}
	}

	for i, t := range allTemplates.Templates() {
		log.Logger().Debug("allTemplates", zap.Int("index", i), zap.Any("t", t.Name()))
	}

}

//编译指定模板。 fn为模板文件全路径
func parseTemplate(fn string) {

	if strings.HasSuffix(strings.ToUpper(fn), strings.ToUpper(suffix)) { //匹配文件

		//模板名字 /controller/action.html
		tn := strings.Replace(fn, TMPLATE_DIR, "", 1) //取相对路径
		tn = strings.ToLower(strings.Replace(tn, "\\", "/", -1))

		log.Logger().Debug("Parsing template",
			zap.String("templateName", tn),
			zap.String("templateFileName", fn))

		// 忽略 LAYOUT_TEMPLATE， 这个文件不独立添加模板。
		//if strings.Compare(tn, LAYOUT_TEMPLATE) != 0 {

		//检查模板是否存在，否则新建
		t := allTemplates.Lookup(tn)
		if t == nil {
			log.Logger().Debug("New template", zap.String("templateName", tn))
			t = allTemplates.New(tn)
		}

		log.Logger().Debug("parse template", zap.String("tn", tn), zap.String("fn", fn))

		/*log.Logger().Debug("parse template",
		zap.String("LAYOUT_TEMPLATE", layoutFN), zap.String("fn", fn))

		//Shared目录下的 _Layout 共享模板需要作为模板第一个页面。
		t1, err := t.ParseFiles(layoutFN, fn)
		if err != nil {
			log.Logger().Error(" ***** parse ERROR:",
				zap.String("filename", fn), zap.Error(err))
		}

		log.Logger().Debug("parsed template", zap.String("t1.Name", t1.Name()))*/

		//编译模板

		actionPageContent, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Logger().Error("parse", zap.String("fn", fn),
				zap.String("tn", tn),
				zap.Error(err))
		}

		var strBody = string(actionPageContent)
		//views/shared目录之外的需要替换模板。
		if strings.HasPrefix(tn, "/shared/") {
			log.Logger().Info("shared file, ", zap.String("templateName", tn))
		} else { //layout 模板文件需要重新加载所有模板。
			strBody = strings.Replace(LAYOUT_DATA,
				LAYOUT_BODY_TAG,
				string(actionPageContent), -1)
		}
		//编译
		_, err = t.Parse(strBody)

		if err != nil {
			log.Logger().Error("parse", zap.String("fn", fn),
				zap.String("tn", tn),
				zap.Error(err))
		}

		/*
			//编译模板
			content, err := ioutil.ReadFile(fn)
			if err != nil {
				debug.Debug("parse", fn, tn, err)
			}

			_, err = t.Parse(string(content))

			if err != nil {
				debug.Debug(" * parse ERROR:", fn, tn, err)
			}*/

	}
}

func getTemplate(name string) *template.Template {
	//return allTemplates[name]
	return allTemplates.Lookup(name)
}

//Render execute the template.
// if template not found , renturn nil,error
func Render(ctlName, actName string, ViewData map[string]interface{}) ([]byte, error) {

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
	err := t.Execute(b, ViewData)
	if err != nil {
		return b.Bytes(), err

	} else {
		return b.Bytes(), nil
	}

}

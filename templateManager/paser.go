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

	"github.com/dereking/grest/log"
	"go.uber.org/zap"

	"github.com/fsnotify/fsnotify"
)

const (
	LAYOUT_BODY_TAG          = "{{ @RenderBody() }}"
	LAYOUT_TEMPLATE_DIRNAME  = "Shared"  //layout模板的文件在views目录下的文件夹名称,区分大小写
	LAYOUT_TEMPLATE_FILENAME = "_Layout" //layout模板的文件名路径,区分大小写，后缀有suffix确定
)

var (
	LAYOUT_TEMPLATE_FILE string //layout模板的文件 路径,区分大小写.文件后缀必须小写。 views/Shared/_Layout.html
)

//var allTemplates map[string]*template.Template
var allTemplates *template.Template

var suffix string      //模板文件后缀，包含.   大写
var TMPLATE_DIR string //模板存放目录，默认运行目录下的views ， 不包含末尾斜杠

var LAYOUT_DATA string //layout文件的内容

var PthSep string //文件路径分隔符

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

	//注册模板函数
	allTemplates = allTemplates.Funcs(template.FuncMap{"html": templateFunc_html})
	allTemplates = allTemplates.Funcs(template.FuncMap{"fileSize": templateFunc_FileSize})
	allTemplates = allTemplates.Funcs(template.FuncMap{"datetime": templateFunc_DateTime})
	allTemplates = allTemplates.Funcs(template.FuncMap{"add": templateFunc_add})

	LAYOUT_TEMPLATE_FILE := fmt.Sprintf("%s%s%s%s%s%s",
		TMPLATE_DIR,
		PthSep,
		LAYOUT_TEMPLATE_DIRNAME,
		PthSep,
		LAYOUT_TEMPLATE_FILENAME,
		strings.ToLower(suffix))

	//读取layout数据
	LAYOUT_DATA = load_layoutfile(LAYOUT_TEMPLATE_FILE)

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

	for i, t := range allTemplates.Templates() {
		log.Logger().Debug("allTemplates", zap.Any("index", i), zap.Any("t", t.Name()))
	}

}

func load_layoutfile(fnRelative string) (ret string) {
	//读取layout数据
	content, err := ioutil.ReadFile(fnRelative)
	if err != nil {
		log.Logger().Error("Loading layout template err:", zap.Error(err),
			zap.String("fnRelative", fnRelative))
	} else {
		ret = string(content)
	}
	return ret
}

//遍历模板目录，并进行模板编译。如果需要动态监视模板变动，进行目录监控。
// templateRoot 模板的当前目录。以/为根目录。
// dirRelative 相对目录
func loadTemplateDir(templateRoot, dirRelative string, bMoniteTemplate bool) {

	fs, err := ioutil.ReadDir(dirRelative)
	if err != nil {
		log.Logger().Error("Loading template err:", zap.Error(err))
		return
	}

	if bMoniteTemplate {
		//moniter this dir
		err = watcher.Add(dirRelative)
		if err != nil {
			log.Logger().Error("watcher.Watch err",
				zap.String("TMPLATE_DIR", TMPLATE_DIR),
				zap.Error(err))
		}
	}

	for _, f := range fs {
		//文件名。相对目录
		fnRelative := dirRelative + PthSep + f.Name()

		if f.IsDir() {
			//递归子目录遍历
			loadTemplateDir(templateRoot+strings.ToLower(f.Name())+"/", fnRelative, bMoniteTemplate)
		} else {
			//处理当前模板，编译并记录。
			parseTemplate(fnRelative)
		}
	}
}

//编译指定模板。 fnRelative 为模板文件相对路径
func parseTemplate(fnRelative string) {
	if !strings.HasSuffix(strings.ToUpper(fnRelative), strings.ToUpper(suffix)) {
		return
	}

	//模板名字 /controller/action.html 全小写
	tn := strings.Replace(fnRelative, TMPLATE_DIR, "", 1) //取相对路径
	tn = strings.ToLower(strings.Replace(tn, "\\", "/", -1))

	//如果是layout文件，直接退出函数，忽略掉。
	// /shared/_layout.html
	LAYOUT_TEMPLATE_Name := strings.ToLower(fmt.Sprintf("/%s/%s%s",
		LAYOUT_TEMPLATE_DIRNAME,
		LAYOUT_TEMPLATE_FILENAME,
		suffix))
	if strings.Compare(tn, LAYOUT_TEMPLATE_Name) == 0 {
		//忽略 layout 文件
		return
	}

	//检查模板是否存在，否则新建
	t := allTemplates.Lookup(tn)
	if t == nil {
		log.Logger().Debug("New template",
			zap.String("templateName", tn),
			zap.String("templateFileName", fnRelative))
		t = allTemplates.New(tn)
	} else {
		log.Logger().Warn("Overrided template",
			zap.String("templateName", tn),
			zap.String("templateFileName", fnRelative))
	}

	//编译模板
	actionPageContent, err := ioutil.ReadFile(fnRelative)
	if err != nil {
		log.Logger().Error("parse",
			zap.String("fnRelative", fnRelative),
			zap.String("tn", tn),
			zap.Error(err))
	}

	var strBody = string(actionPageContent)
	// /views/shared目录之外的需要加载layout模板。
	if strings.HasPrefix(tn, "/shared/") {
		//log.Logger().Info("shared file, ", zap.String("templateName", tn))
	} else {
		// LAYOUT_DATA 不爲空那么 加载 layout 到当前模板
		if len(LAYOUT_DATA) != 0 {
			strBody = strings.Replace(LAYOUT_DATA,
				LAYOUT_BODY_TAG,
				string(actionPageContent), -1)
		}
	}
	//编译
	_, err = t.Parse(strBody)

	if err != nil {
		log.Logger().Error("parse",
			zap.String("templateName", tn),
			zap.String("templateFileName", fnRelative),
			zap.Error(err))
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

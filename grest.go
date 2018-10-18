package grest

/*
静态文件 放在 static 目录下.
	默认 css/img/js 在 static 目录之下.


路由规则: /controllerName/actionName
  controllerName 和 actionName 不区分大小写.
  但是  actionName对应的函数名 必须大写开头(),否则无法找到actionName对应的函数.

例如:
s.AddController("hotel", reflect.TypeOf(controller.TestController{}))
那么:
   /Hotel/IndexFunc
   /Hotel/InDexFunc
   /Hotel/InDexFunc
   /hoteL/IndexFunc
   /HoTel/inDexFunc 都会匹配并执行 controller.HotelController.IndexFunc()

最佳实践:
	s.AddController("hotel", reflect.TypeOf(controller.HotelController{}))

	func ( c* HotelController) HotelOrder(){
		//balabala .....
	}

	那么用户可以通过以下url访问 action 都是可以的:
	1. /hotel/Hotelorder
	2. /hotel/hotelorder
	3. /Hotel/Hotelorder
	4. /hotel/HotelOrder
	5. /hotel/HOTELORDER
	6. /HOTEL/HOTELORDER
	7. .....
*/

import (
	"fmt"
	//"io/ioutil"
	"log"
	"net"
	"net/http"
	//"net/url"
	"reflect"
	"strconv"
	"strings"

	_ "golang.org/x/net/netutil"

	redisCache "github.com/dereking/grest/cache/redis"
	"github.com/dereking/grest/config"
	"github.com/dereking/grest/debug"
	"github.com/dereking/grest/mvc"
	"github.com/dereking/grest/session"
	memsession "github.com/dereking/grest/session/providers/memory"
	"github.com/dereking/grest/templateManager"
	"golang.org/x/net/websocket"
)

const (
	max = 100
)

type GrestServer struct {
	handlerMap map[string]reflect.Type
	listener   net.Listener
}

func NewGrestServer(confName string) *GrestServer {
	s := &GrestServer{handlerMap: make(map[string]reflect.Type, 0)}

	//必须最先初始化
	config.Initialize(confName)
	redisCache.Initialize()

	//initialize memsession first.
	memsession.Initialize()
	session.GetSessionManager()

	templateManager.Initialize()
	return s
}

func (s *GrestServer) Serve() {
	addr := config.AppConfig.StringDefault("addr", ":8000")
	log.Println("addr", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Listen: %v", err)
	}
	defer listener.Close()

	//limit listener
	//l = netutil.LimitListener(l, max)

	http.Handle("/css/", http.FileServer(http.Dir("static")))
	http.Handle("/js/", http.FileServer(http.Dir("static")))
	http.Handle("/img/", http.FileServer(http.Dir("static")))
	http.Handle("/fonts/", http.FileServer(http.Dir("static")))
	http.Handle("/public/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/", s.ServeHTTPDispatch)

	//limit listener
	//listener = netutil.LimitListener(listener, max)

	log.Println("Server running at", listener.Addr())

	http.Serve(listener, nil)

}

func (s *GrestServer) AddController(name string, ctlType reflect.Type) {

	debug.Debug("注册 controller", name, ctlType)

	if s.handlerMap[strings.ToLower(name)] != nil {
		log.Println("WARNING 重复注册 controller", name)
	}

	s.handlerMap[strings.ToLower(name)] = ctlType

	checkController(ctlType)

}

//从url解析出controller 和action 的名称
//controller name : lowcase ; 小写
//action name: with upercase first letter. 首字母大写,其余不变.
func (s *GrestServer) parseControllerAction(path string) (controllerName string, actionName string) {

	//获取controller和action
	list := strings.Split(path, "/")
	pathList := make([]string, 0)
	for _, i := range list {
		tmp := strings.TrimSpace(i)
		if tmp != "" {
			pathList = append(pathList, tmp)
		}
	}

	//如果未指定,使用默认的Home.Index
	controllerName = "home"
	if len(pathList) >= 1 {
		controllerName = strings.ToLower(pathList[0])
	}

	//大小写敏感 case sencitive
	actionName = "Index"
	if len(pathList) >= 2 {
		actionName = strings.ToLower(pathList[1])

		c := actionName[0]
		if c >= 'a' && c <= 'z' {
			c = c - 32
		}
		actionName =
			fmt.Sprintf("%c%s", c, actionName[1:])
	}

	return controllerName, actionName
}

func parseParam(r *http.Request, vals map[string]string) {

	if err := r.ParseForm(); err != nil {
		log.Println("grest parseParam err: ", err)
	}

	log.Println(r.PostForm)

	for k, v := range r.URL.Query() {
		vals[strings.ToLower(k)] = v[0]
	}
	for k, v := range r.PostForm {
		vals[strings.ToLower(k)] = v[0]
	}
	/*
		ct := r.Header.Get("Content-Type")
		// RFC 2616, section 7.2.1 - empty type
		//   SHOULD be treated as application/octet-stream
		if ct == "" {
			ct = "application/octet-stream"
		}
		ct, _, err = mime.ParseMediaType(ct)
		switch {
		case ct == "application/json":
			d, err := ioutil.ReadAll(r.Body)if err != nil {
				log.Println("grest Body ReadAll err: ", err)
			} else {
				//如果
				vals["json"] = string(d)
			}
		}*/

}

//stringToReflectField 把表单\query的数据str转换成对应field类型的值,并赋值给field
func (s *GrestServer) stringToReflectField(field reflect.Value, str string) {
	if field.CanSet() {
		fieldKind := field.Type().Kind()
		switch fieldKind {

		case reflect.String:
			field.Set(reflect.ValueOf(str))

		case reflect.Int:
			num, err := strconv.ParseInt(str, 10, 32)
			if err == nil {
				field.Set(reflect.ValueOf(int(num)))
			}
		case reflect.Int8:
			num, err := strconv.ParseInt(str, 10, 8)
			if err == nil {
				field.Set(reflect.ValueOf(int8(num)))
			}
		case reflect.Int16:
			num, err := strconv.ParseInt(str, 10, 16)
			if err == nil {
				field.Set(reflect.ValueOf(int16(num)))
			}
		case reflect.Int32:
			num, err := strconv.ParseInt(str, 10, 32)
			if err == nil {
				field.Set(reflect.ValueOf(int32(num)))
			}
		case reflect.Int64:
			num, err := strconv.ParseInt(str, 10, 64)
			if err == nil {
				field.Set(reflect.ValueOf(num))
			}

		case reflect.Uint:
			num, err := strconv.ParseUint(str, 10, 64)
			if err == nil {
				field.Set(reflect.ValueOf(uint(num)))
			}
		case reflect.Uint8:
			num, err := strconv.ParseUint(str, 10, 8)
			if err == nil {
				field.Set(reflect.ValueOf(uint8(num)))
			}
		case reflect.Uint16:
			num, err := strconv.ParseUint(str, 10, 16)
			if err == nil {
				field.Set(reflect.ValueOf(uint16(num)))
			}
		case reflect.Uint32:
			num, err := strconv.ParseUint(str, 10, 32)
			if err == nil {
				field.Set(reflect.ValueOf(uint32(num)))
			}
		case reflect.Uint64:
			num, err := strconv.ParseUint(str, 10, 64)
			if err == nil {
				field.Set(reflect.ValueOf(uint64(num)))
			}

		case reflect.Float32:
			num, err := strconv.ParseFloat(str, 32)
			if err == nil {
				field.Set(reflect.ValueOf(float32(num)))
			}
		case reflect.Float64:
			num, err := strconv.ParseFloat(str, 64)
			if err == nil {
				field.Set(reflect.ValueOf(num))
			}

		case reflect.Bool:
			b, err := strconv.ParseBool(str)
			if err == nil {
				field.Set(reflect.ValueOf(b))
			}
		default:
			log.Println("WARNING: Action parameter : Unsupport field type:", fieldKind)
		}
	} else {
		log.Println("WARNING: Action parameter field name must be start with Upper leter")
	}
}

//map HTTP POST/QUERY data to func parameter object.
//把http请求参数映射到 action 对应函数的参数, 返回同参数类型的object.
//如果 action 没有参数, 返回 invalid 的reflect.Value
func (s *GrestServer) formToActionParameter(
	theAction reflect.Value, vals map[string]string) reflect.Value {

	var arg reflect.Value

	//if action has a struct parameter
	//生成 action 的参数列表
	funcType := theAction.Type()
	if funcType.NumIn() == 1 {
		argType := funcType.In(0)
		if argType.Kind() == reflect.Struct {
			//new a struct of the parameter.
			arg = reflect.New(argType).Elem()
			for i := 0; i < argType.NumField(); i++ {

				field := arg.Field(i)
				fieldName := strings.ToLower(argType.Field(i).Name)

				s.stringToReflectField(field, vals[fieldName])
			}
		}
	}

	return arg
}

//call the action of controller.
//调用 action 函数
func (s *GrestServer) callAction(theControllerReflect reflect.Value,
	actionName string, theAction reflect.Value,
	vals map[string]string) (ret mvc.IActionResult) {

	defer func() {
		err := recover()
		if err != nil {
			ret = mvc.HttpInternalError(err)
		}
	}()

	//执行 ActionExecuting 过滤器
	aec := mvc.NewActionExecutingContext(actionName,
		vals)
	executeFilterExecuting(theControllerReflect, aec)

	//过滤器是否返回了结果. 如果是, 则停止执行action
	if aec.Result != nil {
		return aec.Result
	} else { //继续执行action

		arg := s.formToActionParameter(theAction, vals)

		//如果action有参数,才传入参数.
		var args []reflect.Value
		if arg.Kind() != reflect.Invalid {
			args = []reflect.Value{arg}
		}

		//执行 action
		rets := theAction.Call(args)
		if len(rets) != 1 {
			log.Println("Controller 返回值数目错误", rets)
		} else {
			ret = rets[0].Interface().(mvc.IActionResult)
		}
		return ret
	}
}

//call the ws action of controller.
//调用 ws 的 action 函数
func (s *GrestServer) callActionWS(theControllerReflect reflect.Value,
	actionName string, theAction reflect.Value,
	ws *websocket.Conn) {

	defer func() {
		err := recover()
		if err != nil {
			//ret = mvc.HttpInternalError(err)
			ws.Write([]byte(err.(error).Error()))
			ws.Close()
		}
	}()

	args := make([]reflect.Value, 1)
	args[0] = reflect.ValueOf(ws)

	//执行 action
	theAction.Call(args)
}

//http 请求入口, 根据类型派发到ws和http.
func (s *GrestServer) ServeHTTPDispatch(w http.ResponseWriter, r *http.Request) {

	debug.Debug(r)
	debug.Debug(r.Proto)
	//http.Handle("/ws", websocket.Handler(s.ServeHTTPWS))
	if len(r.Header.Get("Sec-WebSocket-Version")) > 0 {
		ws := websocket.Handler(s.serveHTTPWS)
		ws.ServeHTTP(w, r)
		return
	} else {
		s.serveHTTP(w, r)
	}
}

func (s *GrestServer) serveHTTP(w http.ResponseWriter, r *http.Request) {

	//获取controller和action
	controllerName, actionName := s.parseControllerAction(r.URL.Path)
	debug.Debug("controllerName=", controllerName)
	debug.Debug("actionName=", actionName)

	vals := make(map[string]string)
	parseParam(r, vals)
	debug.Debug(" = request data:", len(vals), vals)

	var ar mvc.IActionResult
	var theController mvc.IController //interface{}

	//debug.Debug(s.handlerMap)
	//获取 controllerName 对应的 controller
	ctype := s.handlerMap[controllerName]
	if ctype != nil {

		// new a  controller
		theControllerReflect := reflect.New(ctype)
		theController = theControllerReflect.Interface().(mvc.IController)

		debug.Debug("call controller.Initialize")
		initFunc := theControllerReflect.MethodByName("Initialize")
		initFunc.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)})

		//获取 actionName 对应的 函数.
		//get the func corresponding to the action name.
		theAction := findFuncOfActionInController(theControllerReflect, actionName)

		debug.Debug("find controller." + actionName)
		debug.Debug("executing controller." + actionName)

		if theAction.Kind() != reflect.Invalid {
			//call the action
			ar = s.callAction(theControllerReflect, actionName, theAction, vals)
		} else {
			// action 找不到, 那么返回 template 目录下对应的html 文件
			// If action not found , so render the html file in template dir.
			c := mvc.NewController()
			ar = c.View(controllerName, actionName)
		}
		debug.Debug("executed controller." + actionName)

	} else {
		//controller 不存在, 那么返回 template 目录下对应的html 文件
		// If action not found , so render the html file in template dir.
		debug.Debug("返回 template 目录下对应的html 文件", controllerName, actionName)
		c := mvc.NewController()
		c.Initialize(w, r)
		theController = c
		ar = c.View(controllerName, actionName)
	}

	//render the actionresult.
	//controllerContext := mvc.NewControllerContext(theController, r, w)
	ar.ExecuteResult(theController)
	debug.Debug("ExecuteResult " + actionName)
}

// websocket handle, and dispather.
func (s *GrestServer) serveHTTPWS(ws *websocket.Conn) {
	//io.Copy(ws, ws)
	debug.Debug("ws 已连接;", ws)

	//获取controller和action
	controllerName, actionName := s.parseControllerAction(ws.Request().URL.Path)
	debug.Debug("controllerName=", controllerName)
	debug.Debug("actionName=", actionName)

	//获取 controllerName 对应的 controller
	ctype := s.handlerMap[controllerName]
	if ctype != nil {
		// new a  controller
		theControllerReflect := reflect.New(ctype)

		initFunc := theControllerReflect.MethodByName("InitializeWS")
		initFunc.Call([]reflect.Value{reflect.ValueOf(ws.Request())})

		//获取 actionName 对应的 函数.
		//get the func corresponding to the action name.
		theAction := findFuncOfActionInController(theControllerReflect, actionName)
		if theAction.Kind() != reflect.Invalid {
			//call the action
			s.callActionWS(theControllerReflect, actionName, theAction, ws)
		} else {
			//controller 不存在, 那么返回 404
			ws.WriteClose(404)
		}

	} else {
		//controller 不存在, 那么返回 404
		ws.Write([]byte("Controller not found."))
		ws.WriteClose(404)
	}

}

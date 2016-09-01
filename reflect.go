package grest

import (
	"log"
	"reflect"
	"strings"

	//"github.com/dereking/grest/controller/ActionFilter"
	"github.com/dereking/grest/mvc"
)

//查找controller里面action对应的函数.
//函数名必须大写字母开头(小写开头未导出,无法调用).
func findFuncOfActionInController(theController reflect.Value, actName string) reflect.Value {

	ctlType := theController.Type()

	for i := 0; i < ctlType.NumMethod(); i++ {
		lname := ctlType.Method(i).Name
		if (lname[0] >= 'A') && (lname[0] <= 'Z') {
			if strings.Compare(strings.ToLower(lname), strings.ToLower(actName)) == 0 {
				ret := theController.Method(i)
				return ret
			}
		}

	}

	return reflect.ValueOf(nil)
}

//检查 controller 里是否有重复或异常action.
//比如: action 必须大写开头.
func checkController(ctlType reflect.Type) error {
	funcList := make(map[string]string)

	obj := reflect.New(ctlType)
	ctlType = obj.Type()

	//log.Println(ctlType, ctlType.NumMethod(), obj.NumMethod())

	for i := 0; i < obj.NumMethod(); i++ {
		m := ctlType.Method(i)
		name := m.Name
		lname := strings.ToLower(name)

		funcRetType := ctlType.Method(i).Type
		if funcRetType.NumOut() > 0 {
			retTypeName := funcRetType.Out(0).String()

			//log.Println(retTypeName, name)

			switch retTypeName {
			case "*actionresult.ActionResult",
				"actionresult.IActionResult":

				//log.Println(retTypeName, name)
				if len(funcList[lname]) > 0 {
					log.Panic(ctlType, " action ", name, "跟已有 action 名称冲突: ", funcList[lname])
				} else {
					if (name[0] >= 'A') && (name[0] <= 'Z') {
						funcList[lname] = name
					} else {
						log.Println("WARNING", ctlType, " action name must starts with upper case letter", "[", name, "]")
					}
				}
			}
		}
	}
	return nil
}

//
func executeFilterExecuting(theController reflect.Value, aec *mvc.ActionExecutingContext) {
	OnExecuting := theController.MethodByName("OnExecuting")
	if OnExecuting.Kind() != reflect.Invalid {
		//arg := &ActionFilter.ActionExecutingContext{}
		args := []reflect.Value{reflect.ValueOf(aec)}
		OnExecuting.Call(args)
	}
}

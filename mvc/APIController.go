package mvc

import (
	"errors"
	"log"
	"strings"

	"github.com/dereking/grest/config"
	//"github.com/dereking/grest/controller/ActionFilter"
	"github.com/dereking/grest/utils"
)

//

type APIController struct {
	Controller
}

func (c *APIController) OnExecuting(context *ActionExecutingContext) {

	ip := utils.GetClientIP(c.Request)

	allow, found := config.AppConfig.String("allow.client.ip")

	//配置文件内没有ip限制, 那么允许.
	if !found {
		//context.Result =
		return
	}

	als := strings.Split(allow, ";")
	for _, a := range als {
		if a != "" && a == ip {
			return
		}
	}

	context.Result = c.HttpForbidden()
}

//客户端ip校验. 允许调用则返回nil, 不允许则返回err
func (c *APIController) ClientIPCheck() error {
	req := c.Request
	ip := utils.GetClientIP(req)

	allow, found := config.AppConfig.String("allow.client.ip")

	//配置文件内没有ip限制, 那么允许.
	if !found {
		return nil
	}

	als := strings.Split(allow, ";")
	for _, a := range als {
		if a != "" && a == ip {
			return nil
		}
	}
	log.Println("Client IP not allowed", ip, allow)
	return errors.New("Client IP not allowed")
}

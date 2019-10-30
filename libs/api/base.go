package api

import (
	"../../packs/drivers"
	"../../packs/gin"
	"../../packs/util"
	"log"
	"reflect"
	"strings"
)

type Base struct {
	mod  		map[string]interface{}
	defaultMod 	string
	defaultAct 	string
	db			*drivers.Mysql
}

var baseInstance = &Base{}

func init() {
	log.Println("api base initing")
	baseInstance.defaultMod = "service"
	baseInstance.defaultAct = "/server_list"
	baseInstance.mod = make(map[string]interface{})
	baseInstance.mod["service"] = &Service{}
	baseInstance.db = new (drivers.Mysql)
}

func (this *Base) ResponseJson(c *gin.Context, code int, status int, message interface{}) {
	c.JSON(code, map[string]interface{}{
		"status"	:	status,
		"message"	:	message,
	})
}

func (this *Base) Invoke(c *gin.Context) {
	reqMod := c.Param("mod")
	reqAct := c.Param("act")

	if "" == reqMod && "" == reqAct {
		reqMod = baseInstance.defaultMod
		reqAct = baseInstance.defaultAct
	}else if "" != reqMod && ("" == reqAct || "/" == reqAct) {
		reqAct = baseInstance.defaultAct
	}

	module, exist := baseInstance.mod[reqMod]
	if !exist {
		this.ResponseJson(c, util.HTTPSTATUSCODE_NoContent, util.JSONSTATUSCODE_BadParam, "mod "+reqMod+" not exist")
		return
	}

	first  := strings.ToUpper(reqAct[1:2])	//change the second char into upper
	action := first + reqAct[2:]			//cut the string begin from the third char, first is '/', the second will be replace by it's upper own

	refVal := reflect.ValueOf(module)
	method := refVal.MethodByName(action)

	if method.Kind() == reflect.Invalid {
		this.ResponseJson(c, util.HTTPSTATUSCODE_NoContent, util.JSONSTATUSCODE_BadParam, "act "+action+" not exist")
		return
	}

	c.Set("request_module", module)
	c.Set("request_method", method)

	args := make([]reflect.Value, 1)
	args[0] = reflect.ValueOf(c)
	method.Call(args)

	return
}

func (this *Base) DBInstance() *drivers.Mysql {
	return baseInstance.db
}
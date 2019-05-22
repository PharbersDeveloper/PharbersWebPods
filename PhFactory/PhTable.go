package PhFactory

import (
	"Web/PhDataStorage"
	"Web/PhHandler"
	"Web/PhModel"
	"Web/PhResource"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
)

type PhTable struct{}

var PH_MODEL_FACTORY = map[string]interface{}{
	"PhApplyuser": PhModel.Applyuser{},
}

var PH_STORAGE_FACTORY = map[string]interface{}{
	"PhApplyuserStorage": PhDataStorage.PhApplyuserStorage{},
}

var PH_RESOURCE_FACTORY = map[string]interface{}{
	"PhApplyuserResource": PhResource.PhApplyuserResource{},
}

var PH_FUNCTION_FACTORY = map[string]interface{}{
	"PhCommonPanicHandle":  PhHandler.CommonPanicHandle{},
	"PhSendMailHandler":	PhHandler.PhSendMailHandler{},
	"PhSendBlueBookHandler":PhHandler.PhSendMailHandler{},
	"PhSendSmsHandler":		PhHandler.PhSendSmsHandler{},
	"PhVerifySmsHandler":	PhHandler.PhSendSmsHandler{},

}
var PH_MIDDLEWARE_FACTORY = map[string]interface{}{
}

var PH_DAEMON_FACTORY = map[string]interface{}{
	"BmMongodbDaemon": BmMongodb.BmMongodb{},
	"BmRedisDaemon":   BmRedis.BmRedis{},
}

func (t PhTable) GetModelByName(name string) interface{} {
	return PH_MODEL_FACTORY[name]
}

func (t PhTable) GetResourceByName(name string) interface{} {
	return PH_RESOURCE_FACTORY[name]
}

func (t PhTable) GetStorageByName(name string) interface{} {
	return PH_STORAGE_FACTORY[name]
}

func (t PhTable) GetDaemonByName(name string) interface{} {
	return PH_DAEMON_FACTORY[name]
}

func (t PhTable) GetFunctionByName(name string) interface{} {
	return PH_FUNCTION_FACTORY[name]
}

func (t PhTable) GetMiddlewareByName(name string) interface{} {
	return PH_MIDDLEWARE_FACTORY[name]
}

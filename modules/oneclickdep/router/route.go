package router

import (
	"encoding/json"
	"github.com/ant0ine/go-json-rest/rest"
	"io/ioutil"
	"log"
	"net/http"
	"opensource/chaos/modules/oneclickdep/entity"
	"opensource/chaos/modules/oneclickdep/handler"
	webUtils "opensource/chaos/modules/oneclickdep/utils"
	"runtime/debug"
)

func Route() (rest.App, error) {

	return rest.MakeRouter(
		// 新增某个服务
		rest.Post("/apps", restGuarder(handler.CreateAppsHandler)),
		// 获取所有服务的基本信息
		rest.Get("/apps", restGuarder(handler.GetInfoAppsHandler)),
		// TODO 获取某一个服务的详细信息
		rest.Get("/apps/*appId", restGuarder(handler.GetSingleAppsHandler)),
		// 删除某一个服务的所有实例
		rest.Delete("/apps/*appId", restGuarder(handler.DeleteAppsHandler)),
		// 新增或者更新一批服务
		rest.Post("/apps/updater", restGuarder(handler.CreateOrUpdateAppsHandler)),
		// 回滚一批服务
		rest.Post("/apps/rollback", restGuarder(handler.RollbackAppsHandler)),
		// 新增一批组信息
		rest.Post("/groups", restGuarder(handler.DeployGroupsHandler)),

		// 给新增容器获取一个IP
		rest.Get("/ipholder/*cId", restGuarder(handler.CreateIpForContainer)),
		// 给删除容器释放所占用的IP
		rest.Delete("/ipholder/*cId", restGuarder(handler.DeleteIpForContainer)),

		// 填写某个容器信息
		rest.Post("/containers", restGuarder(handler.CreateContainerInfo)),
		// 更新容器状态
		rest.Post("/containers/:cId/:cState", restGuarder(handler.UpdateStateContainerInfo)),
		// 软删除容器
		rest.Delete("/containers/soft/*cId", restGuarder(handler.MaskContainerInfo)),
	)
}

func restGuarder(method RestFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		// begin := time.Now().UnixNano()
		defer func() {
			if e, ok := recover().(error); ok {
				rest.Error(w, e.Error(), http.StatusInternalServerError)
				log.Println("catchable system error occur: ", e.Error())
				debug.PrintStack()
			}
			// log.Printf("the request: %s cost: %d ms\n", r.URL.RequestURI(), ((time.Now().UnixNano() - begin) / 1000000))
		}()

		pathParams := r.PathParams
		var request entity.CommonRequest
		content, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		webUtils.CheckError(err)
		if len(content) == 0 {
			w.WriteJson(method(pathParams, nil))
			return
		}
		err = json.Unmarshal(content, &request)
		webUtils.CheckError(err)
		switch request.SyncType {
		case "sync":
			log.Println("now use sync mode")
			w.WriteJson(method(pathParams, content))
		case "async":
			log.Println("now use async mode")
			go method(pathParams, content)
			w.WriteJson(map[string]string{"status": "ok"})
		default:
			log.Println("now use default mode(sync)", request.SyncType)
			w.WriteJson(method(pathParams, content))

		}

	}
}

type RestFunc func(pathParams map[string]string, data []byte) interface{}

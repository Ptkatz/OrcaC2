package routers

import (
	"Orca_Server/api"
	"Orca_Server/define/retcode"
	"Orca_Server/servers"
	"net/http"
)

func AccessTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		//检查header是否设置SystemId
		systemId := r.Header.Get("SystemId")
		if len(systemId) == 0 {
			api.Render(w, retcode.FAIL, "系统ID不能为空", []string{})
			return
		}

		//判断是否被注册

		if _, ok := servers.SystemMap.Load(systemId); !ok {
			api.Render(w, retcode.FAIL, "系统ID无效", []string{})
			return
		}

		next.ServeHTTP(w, r)
	})
}

package http

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-transfer/common/g"
	"net/http"
	"os"
	"time"
)

func init() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		RenderJson(w, map[string]interface{}{
			"success": true,
			"msg":     "query success",
			"data":    "ok",
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		RenderJson(w, map[string]interface{}{
			"success": true,
			"msg":     "query success",
			"data":    g.VersionInfo(),
		})
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			var err error

			if err = g.ReloadConfig(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				RenderJson(w, map[string]interface{}{
					"success": false,
					"msg":     err.Error(),
					"data":    nil,
				})
			} else {
				RenderJson(w, map[string]interface{}{
					"success": true,
					"msg":     "reload success",
					"data":    nil,
				})
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			RenderJson(w, map[string]interface{}{
				"success": true,
				"msg":     "exited success",
				"data":    nil,
			})
			go func() {
				time.Sleep(time.Second)
				dlog.Warning("exited..")
				os.Exit(0)
			}()
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})
}

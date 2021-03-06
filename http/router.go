package http

import (
	"encoding/json"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-transfer/common/g"
	"github.com/fanghongbo/ops-transfer/common/model"
	"github.com/fanghongbo/ops-transfer/common/proc"
	"github.com/fanghongbo/ops-transfer/rpc"
	"github.com/fanghongbo/ops-transfer/sender"
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

	http.HandleFunc("/v1/push", func(w http.ResponseWriter, r *http.Request) {
		var (
			metric  []*model.MetricValue
			err     error
			decoder *json.Decoder
		)

		if r.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "body is blank",
				"data":    nil,
			})
			return
		}

		decoder = json.NewDecoder(r.Body)
		err = decoder.Decode(&metric)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     err.Error(),
				"data":    nil,
			})
			return
		}

		if err = rpc.RecvMetricValues(metric, &model.TransferResponse{}, "http"); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     err.Error(),
				"data":    nil,
			})
		} else {
			RenderJson(w, map[string]interface{}{
				"success": true,
				"msg":     "push success",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/counter", func(w http.ResponseWriter, r *http.Request) {
		RenderJson(w, proc.GetAll())
	})

	http.HandleFunc("/statistics", func(w http.ResponseWriter, r *http.Request) {
		RenderJson(w, proc.GetAll())
	})

	http.HandleFunc("/proc/step", func(w http.ResponseWriter, r *http.Request) {
		RenderJson(w, map[string]interface{}{
			"min_step": sender.MinStep,
		})
	})
}

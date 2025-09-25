package servlet

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"monitor/internal/config"
	"monitor/internal/utils"
	"net/http"
	"os"
	"path/filepath"
)

type Info struct {
	Check config.Check
	Alert map[string]any
}

func Server(ctx context.Context, c *config.ApplicationConfig) error {
	w := c.Webserver
	tmpl := filepath.Join(c.ConfigDir, w.TemplateFile)
	_, err := os.Stat(tmpl)
	if err != nil {
		utils.Error("file %s is not present", tmpl)
		return err
	}
	tg, err := template.New(w.TemplateFile).ParseFiles(tmpl)
	if err != nil {
		utils.Error(fmt.Sprintf("%v", err))
		return err

	}
	if w.RootDir == "" {
		return errors.New("the specified patttern for the root_dir is invalid ")
	}
	utils.Debug("Handler root %v\n", w.RootDir)
	http.HandleFunc(w.RootDir, handleAlert(c, tg))
	server := &http.Server{Addr: fmt.Sprintf("%s:%d", w.Host, w.Port), Handler: nil}
	go func() {
		e := server.ListenAndServe()
		if e != nil {
			utils.Error(fmt.Sprintf("%v", e))
		}
	}()

	<-ctx.Done()

	return server.Shutdown(ctx)
}

func handleAlert(cfgp *config.ApplicationConfig, t *template.Template) func(w http.ResponseWriter, r *http.Request) {
	tmpl := *t
	cfg := *cfgp
	utils.Debug("Pointers %v %v\n", t, cfgp)
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Debug("Handler params %v %v\n", w, r)

		if r.Method != "GET" {
			utils.Error(fmt.Sprintf("Got non GET rq: %s", r.URL.String()))
			http.Error(w, "Must GET alert", http.StatusBadRequest)
			return
		}
		checkid := r.URL.Query().Get("checkid")
		if checkid == "" {
			utils.Error(fmt.Sprintf("Missing checkid on request to servlet rq: %s", r.URL.String()))
			http.Error(w, "Missing checkid parameter", http.StatusBadRequest)
			return
		}
		alert := r.URL.Query().Get("alert")
		if alert != "" {

		}
		var c Info
		var chk *config.Check
		for i := range cfg.Checks {
			if cfg.Checks[i].ID == checkid {
				chk = &cfg.Checks[i]
			}
		}
		if chk == nil {
			utils.Error(fmt.Sprintf("Missing or invalid check in request to servlet rq: %s", r.URL.String()))
			http.Error(w, "Missing check or invalid checkid parameter", http.StatusBadRequest)
			return
		}
		c.Check = *chk
		utils.Debug("Fetching form storage %v %v", c.Check.Name, alert)
		a, e := cfgp.Storage.Fetch(r.Context(), []byte(c.Check.Name), alert)
		if e == nil {
			c.Alert = a
		}
		var buffer bytes.Buffer

		b := bufio.NewWriter(&buffer)
		utils.Debug(fmt.Sprintf("Template: %s\n", tmpl.Name()))
		err := tmpl.ExecuteTemplate(b, cfg.Webserver.TemplateFile, c)
		if err != nil {
			http.Error(w, fmt.Sprintf("Mapping failed %v", err), http.StatusBadRequest)
		}
		b.Flush()
		w.Write(buffer.Bytes())
	}
}

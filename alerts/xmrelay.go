package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"

	_ "github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
)

type CommandLineOptions struct {
	LogLevel        string
	Host            string
	Subscription    string
	ResourceGroup   string
	AppGroupEmail   string
	AssignmentGroup string
	IPAddr          string
	TargetURL       string
	TemplatePath    string
	Org             string
	Self            string
	Port            int
	Live            bool
}

func getCLO() CommandLineOptions {
	hostname, error := os.Hostname()
	if error != nil {
		panic(error)
	}
	hostname = strings.Split(hostname, ".")[0]
	commandLineOptions := CommandLineOptions{}
	flags := flag.NewFlagSet("fs", flag.ExitOnError)

	flags.StringVar(&commandLineOptions.LogLevel, "log-level", "info", "log level - 'debug', 'info', 'warn' or 'error'")
	flags.StringVar(&commandLineOptions.Host, "host", "0.0.0.0", "webserver host")
	flags.StringVar(&commandLineOptions.Subscription, "subscription", "sub-001", "AZ Subscriiption")
	flags.StringVar(&commandLineOptions.ResourceGroup, "resourcegroup", "rez", "Az Resource Group")
	flags.StringVar(&commandLineOptions.AppGroupEmail, "email", "fis.wm-azure.admins@fisglobal.com", "App Group eMail")
	flags.StringVar(&commandLineOptions.AssignmentGroup, "assign", "TSG - Wealthware", "App Assignment Group")
	flags.StringVar(&commandLineOptions.IPAddr, "ipaddr", "10.73.87.73", "Influx IP address")
	flags.StringVar(&commandLineOptions.TargetURL, "targeturl", "https://fisglobal.xmatters.com/api/integration/1/functions/0295db39-831a-4fad-a30d-d2d821974232/triggers?apiKey=5f13a481-aa3b-479c-af7c-4430fc983801", "xMatter URL")

	flags.StringVar(&commandLineOptions.TemplatePath, "templatepath", ".", "Path to gf template")
	flags.StringVar(&commandLineOptions.Org, "org", "fb67cbd5fa6747e0", "Influx Org ID")
	flags.StringVar(&commandLineOptions.Self, "self", hostname, "My Hostname")

	flags.IntVar(&commandLineOptions.Port, "port", 80, "webserver port")
	flags.BoolVar(&commandLineOptions.Live, "live", false, "Send live alerts")
	flags.Parse(os.Args[1:])
	log.Println(commandLineOptions.Port)
	return commandLineOptions
}

func walk(j interface{}) map[string]string {
	slog.Debug(fmt.Sprintf("Walking: %v\n", j))
	r := make(map[string]string)

	switch j := j.(type) {
	case []string:
		slog.Debug("[string]\n")

		for i := range j {
			w := walk(i)
			for k, v := range w {
				r[k] = v
			}
		}
	case map[string]string:
		slog.Debug("map[string]string\n")

		for s, v := range j {
			slog.Debug(fmt.Sprintf("STRING: %s = %s\n", s, v))
			r[s] = v
		}
	case map[string]interface{}:
		slog.Debug("map[string]interface\n")

		for s, v := range j {
			switch v := v.(type) {
			case string:
				slog.Debug(fmt.Sprintf("STRING: %s = %s\n", s, v))
				r[s] = v

			default:
				w := walk(v)
				for k, v := range w {
					r[k] = v
				}
			}
		}

	case string:
		slog.Debug(fmt.Sprintf("STRING: %s\n", j))
		r["unnamed"] = j
	}
	return r
}

func sendalert(url string, buf bytes.Buffer) {
	resp, err := http.Post(url, "application/json", &buf)
	if err != nil {
		slog.Debug(fmt.Sprintf("Failed to send alert: %v\n", err))
	}
	b, err := io.ReadAll(resp.Body)
	if err == nil {

		slog.Debug(fmt.Sprintf("Got response: %v\n", string(b)))
	}
}

func handleAlert(t *template.Template) func(w http.ResponseWriter, r *http.Request) {
	tmpl := t
	return func(w http.ResponseWriter, r *http.Request) {
		if tmpl == nil {
			panic("Nil Template Pointer")
		}
		if r.Method != "POST" {
			slog.Error(fmt.Sprintf("Got non POST rq: %s", r.URL.String()))
			http.Error(w, "Must POST alert", http.StatusBadRequest)
			return
		}
		slog.Debug(fmt.Sprintf("Got a request: %s", r.URL.String()))

		bodyContent, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Request body invalid", http.StatusBadRequest)

			return
		}
		slog.Debug(fmt.Sprintf("Got POST rq: %s\n with body:\n%s\n----------------\n\n", r.URL.String(), string(bodyContent)))

		var j interface{}
		json.Unmarshal(bodyContent, &j)
		content := walk(j)
		lv := content["_level"]
		content["status"] = "firing"
		switch lv {
		case "crit":
			content["severity"] = "1"
		case "warn":
			content["severity"] = "2"
		case "info":
			content["severity"] = "3"
		default:
			content["severity"] = "4"
			content["status"] = "resolved"
		}
		content["self"] = o.Self
		content["uuid"] = uuid.NewString()
		content["resourceGroup"] = o.ResourceGroup
		content["subscription"] = o.Subscription
		if content["environment"] == "" {
			content["environment"] = "Test"
		}
		if content["appGroupEmail"] == "" {
			content["appGroupEmail"] = o.AppGroupEmail
		}
		if content["assignmentGroup"] == "" {
			content["assignmentGroup"] = o.AssignmentGroup
		}
		if content["ipaddr"] == "" {
			content["ipaddr"] = o.IPAddr
		}
		if content["recipients"] == "" {
			content["recipients"] = content["assignmentGroup"]
		}
		if content["org"] == "" {
			content["org"] = o.Org
		}
		var buffer bytes.Buffer
		b := bufio.NewWriter(&buffer)
		slog.Debug(fmt.Sprintf("Template: %s\n", tmpl.Name()))
		err = tmpl.ExecuteTemplate(b, "gf.tmpl", content)
		if err != nil {
			http.Error(w, "Mapping failed", http.StatusBadRequest)
		}
		b.Flush()
		//	target := o.TargetURL + "&recipients=" + content["recipients"]
		target, err := url.Parse(o.TargetURL)
		parameters := url.Values{}
		for k, v := range target.Query() {
			fmt.Printf("%s = %s\n", k, v[0])
			parameters.Add(k, v[0])
		}
		parameters.Add("recipients", content["recipients"])
		target.RawQuery = parameters.Encode()
		if err != nil {
			http.Error(w, "URL Build failed", http.StatusBadRequest)
		}
		slog.Debug(fmt.Sprintf("Target: %s\n", target.String()))
		if o.Live {
			sendalert(target.String(), buffer)
		}
		if !o.Live || o.LogLevel == "debug" {
			fmt.Println(buffer.String())
		}
		w.Write([]byte("OK"))
	}
}

var o CommandLineOptions

func main() {

	o = getCLO()
	tg := template.Must(template.New("gf.tmpl").ParseFiles(o.TemplatePath + "/gf.tmpl"))
	if o.LogLevel == "debug" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	http.HandleFunc("/alert", handleAlert(tg))
	e := http.ListenAndServe(fmt.Sprintf("%s:%d", o.Host, o.Port), nil)
	if e != nil {
		panic(e)
	}

}

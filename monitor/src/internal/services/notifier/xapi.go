package notifier

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"monitor/internal/utils"
	"net/http"
	"net/url"
	"text/template"

	"github.com/pkg/errors"
)

type XAPI struct {
	TemplateName string
	Content      Content
	Template     *template.Template
	TargetURL    string
}

func (x *XAPI) SendAlert(content map[string]string) {
	t, err := x.GetTemplate()
	if err != nil {
		utils.Error("Template Build failed")
		return
	}
	buf, err := render(content, *t)
	if err != nil {
		utils.Error("Template render")
		return
	}
	target, err := url.Parse(x.GetTargetURL())

	parameters := url.Values{}
	for k, v := range target.Query() {
		fmt.Printf("%s = %s\n", k, v[0])
		parameters.Add(k, v[0])
	}
	parameters.Add("recipients", content["recipients"])
	target.RawQuery = parameters.Encode()
	if err != nil {
		utils.Error("URL Build failed")
	}
	utils.Info("Target: %s\n", target.String())
	utils.Info("Payload: %s\n", buf.String())
	resp, err := http.Post(target.String(), "application/json", &buf)
	if err != nil {
		utils.Warn("Failed to send alert: %v\n", err)
	} else {
		b, err := io.ReadAll(resp.Body)
		if err == nil {

			utils.Error("Got response: %v\n", string(b))
		}
	}

}

func (x *XAPI) GetTemplate() (*template.Template, error) {
	if x.Template == nil {
		return nil, errors.New("Nil template")
	} else {
		return x.Template, nil
	}
}

func (x *XAPI) GetTargetURL() string {
	return x.TargetURL
}

func (x *XAPI) GetSrcEmail() string {
	return ""
}
func (x *XAPI) GetDstEmail() string {
	return ""
}

func render(content Content, tmpl template.Template) (buffer bytes.Buffer, err error) {

	b := bufio.NewWriter(&buffer)
	utils.Info("Template: %s\n", tmpl.Name())
	err = tmpl.Execute(b, content)
	if err != nil {
		return
	}
	b.Flush()
	return
}

package notifier

import (
	"log"
	"monitor/internal/utils"

	"text/template"

	"github.com/pkg/errors"
)

type Logger struct {
	TemplateName string
	Template     *template.Template
	SrcEmail     string
	DstEmail     string
}

func (x *Logger) SendAlert(content Content) {
	t, err := x.GetTemplate()
	if err != nil {
		utils.Error("Template Build failed %v", err)
		return
	}
	buf, err := render(content, *t)
	if err != nil {
		utils.Error("Template render failed %v", err)
		return
	}

	log.Println(buf.String())

}

func (x *Logger) GetTemplate() (*template.Template, error) {
	if x.Template == nil {
		return nil, errors.New("Nil template")
	} else {
		return x.Template, nil
	}

}
func (x *Logger) GetTargetURL() string {
	return ""
}
func (x *Logger) GetSrcEmail() string {
	return ""
}
func (x *Logger) GetDstEmail() string {
	return ""
}

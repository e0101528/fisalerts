package notifier

import (
	"io"
	"monitor/internal/utils"
	"os/exec"
	"text/template"

	"github.com/pkg/errors"
)

type Postfix struct {
	TemplateName string
	Template     *template.Template
	SrcEmail     string
	DstEmail     string
}

func (x *Postfix) SendAlert(content Content) {
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
	src := x.GetSrcEmail()
	dst := x.GetDstEmail()

	cmd := exec.Command("/usr/bin/mailx", "-s", "Dummy", "-r", src, dst)
	// cmd := exec.Command("cat")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		utils.Error("postfix sednalert error: %v", err)
		return
	}
	utils.Debug("%s %s %s %s %s %s\n", "mailx", "-s", "'I am a Dummy'", "-r", src, dst)
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, buf.String())
	}()
	err = cmd.Start()
	if err != nil {
		utils.Error("postfix error - %v", err)
	} else {
		utils.Debug("Mailx PID: %v", cmd.Process.Pid)
		err = cmd.Wait()
		utils.Debug("Command finished with error: %v", err)
	}

	utils.Debug("Done eMails\n")
}

func (x *Postfix) GetTemplate() (*template.Template, error) {
	if x.Template == nil {
		return nil, errors.New("Nil template")
	} else {
		return x.Template, nil
	}

}
func (x *Postfix) GetTargetURL() string {
	return ""
}
func (x *Postfix) GetSrcEmail() string {
	return x.SrcEmail
}
func (x *Postfix) GetDstEmail() string {
	return x.DstEmail
}

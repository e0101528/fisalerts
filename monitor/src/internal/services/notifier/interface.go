package notifier

import (
	"text/template"
)

type INotifier interface {
	SendAlert(map[string]string)
	GetTemplate() (*template.Template, error)
	GetTargetURL() string
}

type Content map[string]string

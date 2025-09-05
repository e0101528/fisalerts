package notifier

import (
	"text/template"
)

type INotifier interface {
	SendAlert(Content)
	GetTemplate() (*template.Template, error)
	GetTargetURL() string
}

type Content struct {
	Fields map[string]string
	Labels map[string]string
}

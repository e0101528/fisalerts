package config

import (
	"crypto/sha256"
	"flag"
	"fmt"
	repositories "monitor/internal/repositories/bolt"
	"monitor/internal/services/notifier"
	"monitor/internal/utils"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"
)

func LoadAppConfig() (appConfig ApplicationConfig, err error) {
	commandLineOptions, _ := getCommandLineOptions()
	configFilePath := filepath.Join(commandLineOptions.ConfigDir, "config.yaml")
	utils.SetLevel(commandLineOptions.LogLevel)
	fileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(fileBytes, &appConfig)
	if err != nil {
		return
	}
	appConfig.MaintenanceMode = !commandLineOptions.Live
	for i, _ := range appConfig.Checks {
		if appConfig.Checks[i].FluxFile != "" {
			fluxfile := relative(appConfig.Checks[i].FluxFile, commandLineOptions.ConfigDir)

			utils.Info("Reading flux file %s\n", fluxfile)

			fileBytes, err := os.ReadFile(fluxfile)
			if err != nil {
				utils.Error("Flux file read error for %s\n", appConfig.Checks[i].FluxFile)
			} else {
				appConfig.Checks[i].Flux = string(fileBytes)
				appConfig.Checks[i].ID = fmt.Sprintf("%x", sha256.Sum256(fileBytes))
			}

		}
	}

	appConfig.Storage, err = repositories.InitBolt()

	appConfig.Notifiers = make(map[string]notifier.INotifier)
	var tmpl *template.Template
	for i, t := range appConfig.Targets {
		fileBytes, e := os.ReadFile(relative(t.TemplateFile, commandLineOptions.ConfigDir))
		if e == nil {
			tmpl, e = template.New(i).Parse(string(fileBytes))
		} else {
			utils.Error("Template file load failed with %v\n", e)
		}

		if e == nil {

			switch t.Method {
			case "postfix":
				appConfig.Notifiers[i] = &notifier.Postfix{
					TemplateName: i,
					Template:     tmpl,
					SrcEmail:     t.SrcEmail,
					DstEmail:     t.DstEmail,
				}
			case "xapi":
				appConfig.Notifiers[i] = &notifier.XAPI{
					TemplateName: i,
					Template:     tmpl,
					TargetURL:    t.TargetURL,
				}
			case "logger":
				appConfig.Notifiers[i] = &notifier.Logger{
					TemplateName: i,
					Template:     tmpl,
				}
			}
		} else {
			utils.Error("Template creation failed with %v\n", e)

		}
	}

	return
}

func getCommandLineOptions() (options CommandLineOptions, err error) {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	flags.StringVar(&options.ConfigDir, "config-dir", "", "path to config directory")
	flags.StringVar(&options.LogFormat, "log-format", "json", "log format - 'json' or 'pretty'")
	flags.IntVar(&options.LogLevel, "log-level", 0, "0-4 (error, warn, info, debug, dump)")
	flags.BoolVar(&options.Live, "live", false, "enable alerts to live incident management")

	err = flags.Parse(os.Args[1:])
	return
}

func relative(filename string, defaultpath string) string {
	if filename[0] != '/' {
		utils.Debug("Prepending file path with default path %s\n", defaultpath)
		filename = filepath.Join(defaultpath, filename)
	}
	return filename
}

package config

import (
	repositories "monitor/internal/repositories/bolt"
	"monitor/internal/services/notifier"
	"monitor/internal/services/workcalendar"
)

type CommandLineOptions struct {
	ConfigDir                   string `yaml:"config_dir"`
	LogFormat                   string `yaml:"log_format"`
	LogLevel                    int    `yaml:"log_level"`
	WebserverHost               string `yaml:"webserver_host"`
	WebserverPort               int    `yaml:"webserver_port"`
	WebserverRootDir            string `yaml:"webserver_root_dir"`
	WebserverRequestTimeoutSecs int    `yaml:"webserver_request_timeout_secs"`
	WebserverTLSCertFilepath    string `yaml:"webserver_tls_cert_filepath"`
	WebserverTLSKeyFilepath     string `yaml:"webserver_tls_key_filepath"`
	Telemetry                   string `yaml:"telemetry"`
	Live                        bool   `yaml:"live"`
}

// -----------------
type WebServerConfig struct {
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	RootDir            string `yaml:"root_dir"`
	RequestTimeoutSecs int    `yaml:"request_timeout_secs"`
	TLSCertFilepath    string `yaml:"tls_cert_filepath"`
	TLSKeyFilepath     string `yaml:"tls_key_filepath"`
}

// -----------------

type Check struct {
	Name       string   `yaml:"name"`
	Detail     string   `yaml:"detail"`
	Calendar   string   `yaml:"calendar"`
	FluxFile   string   `yaml:"fluxfile"`
	Flux       string   `yaml:"flux"`
	Severity   int      `yaml:"severity"` //1..4 = crit..ok
	Threshold  float64  `yaml:"threshold"`
	Comparison string   `yaml:"comparison"`
	Match      string   `yaml:"match"`
	Tags       []string `yaml:"tags"`
	Interval   int      `yaml:"interval"`
	Target     string   `yaml:"target"`
	ID         string   `yaml:"id"` //overwritten by checksum of flux
}

type Target struct {
	Method       string `yaml:"method"`
	Address      string `yaml:"address"`
	TemplateFile string `yaml:"templatefile"`
	TargetURL    string `yaml:"targeturl"`
	SrcEmail     string `yaml:"srcemail"`
	DstEmail     string `yaml:"dstemail"`
}

type InfluxDBConfig struct {
	ServerURL string `yaml:"serverurl"`
	AuthToken string `yaml:"authtoken"`
	Org       string `yaml:"org"`
}

type ApplicationConfig struct {
	Name            string                           `yaml:"name"`
	Webserver       WebServerConfig                  `yaml:"go_server"`
	Calendars       map[string]workcalendar.Calendar `yaml:"calendars"`
	Checks          []Check
	Targets         map[string]Target `yaml:"targets"`
	Notifiers       map[string]notifier.INotifier
	CalendarMap     *workcalendar.CalendarMap
	Storage         *repositories.TimeSeries
	Organization    string         `yaml:"organization"`
	AzureSub        string         `yaml:"azuresub"`
	AzureRG         string         `yaml:"azurerg"`
	AssignmentGroup string         `yaml:"assignmentgroup"`
	AppGroupEmail   string         `yaml:"appgroupemail"`
	IPAddr          string         `yaml:"ipaddr"`
	InfluxDB        InfluxDBConfig `yaml:"influxdb"`
}

type Metric struct {
	Key string
}

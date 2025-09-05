# Monitor

Monitor is a golang service that parses a set of YAML check configurations, imports associated flux queries and then executed them on the defined schedule and uses alerting pathways as configured.

To simplify config management, the large single YAML file is separated into a base config and multiple check configs. These are they merged into the unified YAML for use in the production system config. 

The resources directory contains those configs. To merge the partials into the main config file:


```
  cd resources/production
  cat config_bare.yaml check/*.check > config.yaml

```

For each check there is a corresponding flux query.
The indentation in the check is **CRITICAL**. They are simply parts of a YAML map[string]struct under the 'checks' key.

The daemon code is under src/

To build the daemon 

```
  cd src/cmd
  go build -o monitor

```

## Command Line flags 

### -config-dir string
        path to config directory
###  -log-level int
        0-4 (error, warn, info, debug, dump)
###  -live
        enable development mode


## YAML Configuration File

### Section 1 - globals 
These are basically needed to form the API endpoint and populate the JSON template for xMatters. 

```
organization: fb67cbd5fa6747e0
self: cibccaprodadm2
azuresub: ca049151-8cb9-421b-b536-251fa31c9a62
azurerg: wealthware_cibc-prod-sql-rg
ipaddr: 10.73.87.73
assignmentgroup: "TSG - Wealthware" 
appgroupemail: fis.wm-azure.admins@fisglobal.com

```
### Section 1 - influxdb data source

The server URL http://{server}:{port}, the org of interest in InfluxDB and a READ-ONLY token must be provided. 
 
```
influxdb:
  serverurl: 
  org: 
  authtoken: 


```


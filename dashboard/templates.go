package main

import "html/template"

var (
	IndexTemplate          = template.Must(template.ParseFiles("dashboard/templates/index.html"))
	DeviceDetailsTemplate  = template.Must(template.ParseFiles("dashboard/templates/device.html"))
	FilesTemplate          = template.Must(template.ParseFiles("dashboard/templates/files.html"))
	BlockedDomainsTemplate = template.Must(template.ParseFiles("dashboard/templates/blocked_domains.html"))
)

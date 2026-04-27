package main

import "html/template"

var (
	IndexTemplate         = template.Must(template.ParseFiles("dashboard/templates/index.html"))
	DeviceDetailsTemplate = template.Must(template.ParseFiles("dashboard/templates/device.html"))
)

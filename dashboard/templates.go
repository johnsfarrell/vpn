package main

import "html/template"

var (
	IndexTemplate         = template.Must(template.ParseFiles("templates/index.html"))
	DeviceDetailsTemplate = template.Must(template.ParseFiles("templates/device.html"))
)

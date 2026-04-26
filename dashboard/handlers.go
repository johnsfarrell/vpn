package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

func RegisterRoutes() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/download-config/", handleDownloadConfig)
	http.HandleFunc("/add-device/", handleAddDevice)
	http.HandleFunc("/delete-device/", handleDeleteDevice)
	http.HandleFunc("/device/", handleDeviceDetails)
}

func handleIndex(w http.ResponseWriter, _ *http.Request) {
	names, err := ListDeviceNames()
	if err != nil {
		http.Error(w, "failed to read devices", http.StatusInternalServerError)
		return
	}

	err = IndexTemplate.Execute(w, names)
	if err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		return
	}
}

func handleDownloadConfig(w http.ResponseWriter, r *http.Request) {
	deviceName, err := pathDeviceName(r.URL.Path, "/download-config/")
	if err != nil {
		http.Error(w, "invalid device name", http.StatusBadRequest)
		return
	}

	config, err := ReadClientConfig(deviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\"client.conf\"")
	_, _ = w.Write(config)
}

func handleAddDevice(w http.ResponseWriter, r *http.Request) {
	deviceName, err := pathDeviceName(r.URL.Path, "/add-device/")
	if err != nil {
		http.Error(w, "invalid device name", http.StatusBadRequest)
		return
	}

	err = SetupClient(deviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Device added: %s", deviceName)
	http.Redirect(w, r, "/device/"+url.PathEscape(deviceName), http.StatusSeeOther)
}

func handleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	deviceName, err := pathDeviceName(r.URL.Path, "/delete-device/")
	if err != nil {
		http.Error(w, "invalid device name", http.StatusBadRequest)
		return
	}

	err = DeleteDevice(deviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Device deleted: %s", deviceName)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleDeviceDetails(w http.ResponseWriter, r *http.Request) {
	deviceName, err := pathDeviceName(r.URL.Path, "/device/")
	if err != nil {
		http.Error(w, "invalid device name", http.StatusBadRequest)
		return
	}

	deviceDetails, err := GetDeviceDetails(deviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = DeviceDetailsTemplate.Execute(w, deviceDetails)
	if err != nil {
		http.Error(w, "failed to render device details page", http.StatusInternalServerError)
		return
	}
}

func pathDeviceName(path, prefix string) (string, error) {
	raw := strings.TrimPrefix(path, prefix)
	raw = strings.TrimSuffix(raw, "/")
	return url.PathUnescape(raw)
}

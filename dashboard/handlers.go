package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type FilesPageData struct {
	Files              []string
	FreeDiskSpaceBytes int64
}

func RegisterRoutes() {
	http.HandleFunc("/", handleIndex)

	http.HandleFunc("/device/", handleDeviceDetails)
	http.HandleFunc("/download-device-config/", handleDownloadDeviceConfig)
	http.HandleFunc("/add-device/", handleAddDevice)
	http.HandleFunc("/delete-device/", handleDeleteDevice)
	
	http.HandleFunc("/files", handleFiles)
	http.HandleFunc("/upload-file", handleUploadFile)
	http.HandleFunc("/download-file/", handleDownloadFile)
	http.HandleFunc("/delete-file/", handleDeleteFile)
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

func handleDownloadDeviceConfig(w http.ResponseWriter, r *http.Request) {
	deviceName, err := pathDeviceName(r.URL.Path, "/download-device-config/")
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

func handleFiles(w http.ResponseWriter, r *http.Request) {
	freeDiskSpaceBytes, err := VMFreeDiskSpaceGB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files, err := ListFiles()
	if err != nil {
		http.Error(w, "failed to list files", http.StatusInternalServerError)
		return
	}
	
	data := FilesPageData{
		Files:              files,
		FreeDiskSpaceBytes: freeDiskSpaceBytes,
	}

	err = FilesTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to render files page", http.StatusInternalServerError)
		return
	}
}

func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to get file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	err = UploadFile(file, header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/files", http.StatusSeeOther)
}

func handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	filename, err := pathDeviceName(r.URL.Path, "/download-file/")
	if err != nil {
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("files", filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	_, _ = w.Write(content)
}

func handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	filename, err := pathDeviceName(r.URL.Path, "/delete-file/")
	if err != nil {
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	err = DeleteFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/files", http.StatusSeeOther)
}
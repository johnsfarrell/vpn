package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

type DeviceDetails struct {
	Name          string `json:"name"`
	PublicKey     string `json:"publicKey"`
	PrivateKey    string `json:"privateKey"`
	DNS           string `json:"dns"`
	IP            string `json:"ip"`
	Endpoint      string `json:"endpoint"`
	LastHandshakeUnix int64 `json:"lastHandshakeUnix"`
	TransferRxBytes   int64 `json:"transferRxBytes"`
	TransferTxBytes   int64 `json:"transferTxBytes"`
	QRCodeDataURL template.URL
	DownloadPath  string
}

func ListDeviceNames() ([]string, error) {
	entries, err := os.ReadDir(filepath.Join("wireguard", "clients"))
	if err != nil {
	  return nil, fmt.Errorf("unable to read clients directory: %w", err)
	}
  
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
	  if entry.IsDir() {
		names = append(names, entry.Name())
	  }
	}
	sort.Strings(names)
	return names, nil
}

func GetDeviceDetails(deviceName string) (*DeviceDetails, error) {
	cmd := exec.Command("bash", "./wireguard/scripts/client_details.sh", deviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("client_details.sh failed: %w: %s", err, string(output))
	}

	var response DeviceDetails
	err = json.Unmarshal(output, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse client_details.sh output: %w", err)
	}

	qrCodeDataURL, err := ClientConfQRCode(deviceName)
	if err != nil {
		return nil, err
	}

	response.QRCodeDataURL = qrCodeDataURL
	response.DownloadPath = "/download-config/" + deviceName
	
	return &response, nil
}

func DeleteDevice(deviceName string) error {
	cmd := exec.Command("bash", "./wireguard/scripts/client_delete.sh", deviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("client_delete.sh failed: %w: %s", err, string(output))
	}
	return nil
}

func SetupClient(deviceName string) error {
	cmd := exec.Command("bash", "./wireguard/scripts/client_create.sh", deviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("client_create.sh failed: %w: %s", err, string(output))
	}
	return nil
}

func ClientConfQRCode(deviceName string) (template.URL, error) {
	configPath := filepath.Join("wireguard", "clients", deviceName, "client.conf")

	cmd := exec.Command("qrencode", "-t", "SVG", "-r", configPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("qrencode failed: %w: %s", err, string(output))
	}

	return template.URL("data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(output)), nil
}

func ReadClientConfig(deviceName string) ([]byte, error) {
	configPath := filepath.Join("wireguard", "clients", deviceName, "client.conf")
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read client config: %w", err)
	}
	return content, nil
}

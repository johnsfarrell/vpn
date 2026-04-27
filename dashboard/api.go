package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DeviceDetails struct {
	Name          string `json:"name"`
	PublicKey     string `json:"publicKey"`
	PrivateKey    string `json:"privateKey"`
	DNS           string `json:"dns"`
	IP            string `json:"ip"`
	Endpoint      string `json:"endpoint"`
	LastHandshakeUnix int64 `json:"lastHandshakeUnix"`
	LastHandshakeAgo  string `json:"-"`
	TransferRxBytes   int64 `json:"transferRxBytes"`
	TransferTxBytes   int64 `json:"transferTxBytes"`
	DownloadMB    float64 `json:"-"`
	UploadMB      float64 `json:"-"`
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
	response.DownloadPath = "/download-device-config/" + deviceName
	response.DownloadMB = float64(response.TransferRxBytes) / 1000 / 1000
	response.UploadMB = float64(response.TransferTxBytes) / 1000 / 1000
	response.LastHandshakeAgo = "Never"
	if response.LastHandshakeUnix > 0 {
		handshakeTime := time.Unix(response.LastHandshakeUnix, 0)
		duration := time.Since(handshakeTime).Round(time.Second)
		response.LastHandshakeAgo = duration.String() + " ago"
	}
	
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

func VMFreeDiskSpaceGB() (int64, error) {
	cmd := exec.Command("df", "--output=avail", "-B1", "/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("df failed: %w: %s", err, string(output))
	}
	fields := strings.Fields(string(output))
	if len(fields) < 2 {
		return 0, fmt.Errorf("unexpected df output: %q", string(output))
	}
	freeBytes, err := strconv.ParseInt(fields[len(fields)-1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse free bytes: %w", err)
	}
	return freeBytes / 1000 / 1000 / 1000, nil
}

func ListFiles() ([]string, error) {
	entries, err := os.ReadDir("files")
	if err != nil {
		return nil, fmt.Errorf("unable to read files directory: %w", err)
	}
	
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.Name() != ".gitkeep" {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func UploadFile(file io.Reader, filename string) error {
	filePath := filepath.Join("files", filepath.Base(filename))
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func DeleteFile(filename string) error {
	filePath := filepath.Join("files", filename)
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func BlockDomain(domainName string) error {
	cmd := exec.Command("bash", "./dns/scripts/block_domain.sh", domainName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("block_domain.sh failed: %w: %s", err, string(output))
	}
	return nil
}

func UnblockDomain(domainName string) error {
	cmd := exec.Command("bash", "./dns/scripts/unblock_domain.sh", domainName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unblock_domain.sh failed: %w: %s", err, string(output))
	}
	return nil
}

func ListBlockedDomains() ([]string, error) {
	cmd := exec.Command("bash", "./dns/scripts/list_blocked_domains.sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("list_blocked_domains.sh failed: %w: %s", err, string(output))
	}
	domains := strings.Fields(string(output))
	sort.Strings(domains)
	return domains, nil
}
package common

import (
	"encoding/json"

	"github.com/filanov/bm-inventory/models"
)

func GetCurrentHostName(host *models.Host) (string, error) {
	var inventory models.Inventory
	if host.RequestedHostname != "" {
		return host.RequestedHostname, nil
	}
	err := json.Unmarshal([]byte(host.Inventory), &inventory)
	if err != nil {
		return "", err
	}
	return inventory.Hostname, nil
}

func GetHostnameForMsg(host *models.Host) string {
	hostName, err := GetCurrentHostName(host)
	if err != nil || hostName == "" {
		return host.ID.String()
	}
	return hostName
}

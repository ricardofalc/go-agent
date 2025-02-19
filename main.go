package main

import (
	"fmt"
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
	"log"
	"os/user"
	"strings"
	"time"
)

// FusionNic represents a network interface configuration with required fields
type FusionNic struct {
	IPAddress   string `json:"ipAddress"`
	Description string `json:"description"` //TODO: Get description of NIC e.g. the name of the NIC / manufacturer
	Gateway     string `json:"gateway"`     //TODO: Get the default gateway
	MacAddress  string `json:"macAddress"`
}

// FusionDevice represents the MVP device payload with required and optional fields
type FusionDevice struct {
	// Required fields
	Hostname        string `json:"hostname"`
	Model           string `json:"model"` //TODO: Get model of device e.g. virtual machine or not
	OS              string `json:"os"`
	OSVersion       string `json:"osVersion"`
	OSKernelVersion string `json:"osKernelVersion"`
	Architecture    string `json:"architecture"`

	// Optional fields (using pointers for nullable values)
	NetworkInterfaces     *[]FusionNic `json:"networkInterfaces,omitempty"`
	Processor             *string      `json:"processor,omitempty"` //TODO: Get processor name, speed, etc.
	PhysicalMemory        *int64       `json:"physicalMemory,omitempty"`
	SerialNumber          *string      `json:"serialNumber,omitempty"`          //TODO: Get serial number
	ExternalIPAddresses   *[]string    `json:"externalIpAddress,omitempty"`     //TODO: Actually get external IP instead of internal with pinging or something?
	InstalledApplications *[]string    `json:"installedApplications,omitempty"` //TODO: Get installed applications
}

var (
	Version   string
	BuildTime string
)

func main() {
	// Format build time for display (replace _ with space)
	displayBuildTime := strings.ReplaceAll(BuildTime, "_", " ")

	// Print version and build info
	fmt.Printf("Fusion Agent Version: %s\n", Version)
	fmt.Printf("Build Time (UTC): %s\n", displayBuildTime)

	// Print current time in UTC
	currentTime := time.Now().UTC().Format("2006-01-02 15:04:05")
	fmt.Printf("Current Date and Time (UTC): %s\n", currentTime)

	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		log.Printf("Error getting current user: %v", err)
	} else {
		fmt.Printf("Current User: %s\n\n", currentUser.Username)
	}
	// Get host information
	host, err := sysinfo.Host()
	if err != nil {
		log.Fatalf("Error getting host info: %v", err)
	}

	// Create and populate device info
	device := createDeviceInfo(host, currentUser)

	// Print the structured data
	printStructuredData(device)

}

func createDeviceInfo(host types.Host, currentUser *user.User) FusionDevice {
	info := host.Info()
	memory, _ := host.Memory()

	// Create base device
	device := FusionDevice{
		Hostname:        info.Hostname,
		Model:           "", //TODO: This might need to be retrieved differently
		OS:              info.OS.Name,
		OSVersion:       info.OS.Version,
		OSKernelVersion: info.KernelVersion,
		Architecture:    info.Architecture,
	}

	// Set memory info
	memoryTotal := int64(memory.Total)
	device.PhysicalMemory = &memoryTotal

	// Set network interfaces
	if len(info.MACs) > 0 || len(info.IPs) > 0 {
		nics := make([]FusionNic, 0)
		for i, mac := range info.MACs {
			nic := FusionNic{
				MacAddress: mac,
			}
			if i < len(info.IPs) {
				nic.IPAddress = info.IPs[i]
			}
			nics = append(nics, nic)
		}
		device.NetworkInterfaces = &nics
	}

	// Set external IPs
	if len(info.IPs) > 0 {
		ips := make([]string, 0)
		for _, ip := range info.IPs {
			ips = append(ips, ip)
		}
		device.ExternalIPAddresses = &ips
	}

	return device
}

func printStructuredData(device FusionDevice) {
	fmt.Println("\n=== FusionDevice Structure ===")
	fmt.Printf("Hostname: %s\n", device.Hostname)
	fmt.Printf("OS: %s\n", device.OS)
	fmt.Printf("OS Version: %s\n", device.OSVersion)
	fmt.Printf("OS Kernel Version: %s\n", device.OSKernelVersion)
	fmt.Printf("Architecture: %s\n", device.Architecture)

	if device.PhysicalMemory != nil {
		gbMemory := float64(*device.PhysicalMemory) / (1024 * 1024 * 1024)
		fmt.Printf("Physical Memory: %.2f GB\n", gbMemory)
	}

	if device.NetworkInterfaces != nil {
		fmt.Println("\nNetwork Interfaces:")
		for i, nic := range *device.NetworkInterfaces {
			fmt.Printf("Interface %d:\n", i+1)
			fmt.Printf("  MAC Address: %s\n", nic.MacAddress)
			if nic.IPAddress != "" {
				fmt.Printf("  IP Address: %s\n", nic.IPAddress)
			}
			if nic.Gateway != "" {
				fmt.Printf("  Gateway: %s\n", nic.Gateway)
			}
			if nic.Description != "" {
				fmt.Printf("  Description: %s\n", nic.Description)
			}
		}
	}

	if device.ExternalIPAddresses != nil {
		fmt.Println("\nExternal IP Addresses:")
		for _, ip := range *device.ExternalIPAddresses {
			fmt.Printf("  %s\n", ip)
		}
	}
}

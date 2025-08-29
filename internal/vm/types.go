package vm

// LaunchConfig contains configuration for launching a new VM
type LaunchConfig struct {
	Image     string // e.g., "debian:bookworm"
	Name      string // VM instance name
	CPUs      int    // Number of CPUs
	Memory    string // Memory size (e.g., "2G")
	Disk      string // Disk size (e.g., "20G")
	CloudInit string // Path to cloud-init file
}

// ImageInfo represents information about a cloud image
type ImageInfo struct {
	Distribution string `json:"distribution"` // e.g., "debian"
	Version      string `json:"version"`      // e.g., "bookworm"
	Architecture string `json:"architecture"` // e.g., "amd64"
	URL          string `json:"url"`          // Download URL
	Checksum     string `json:"checksum"`     // SHA256 checksum
	Size         int64  `json:"size"`         // File size in bytes
}

// NetworkConfig represents network configuration for a VM
type NetworkConfig struct {
	Bridge    string `json:"bridge"`     // Bridge interface name
	IPAddress string `json:"ip_address"` // Static IP (optional)
	MAC       string `json:"mac"`        // MAC address
}
package kvm

import "time"

// VMState represents the possible states of a VM
type VMState string

const (
	StateRunning   VMState = "Running"
	StateStopped   VMState = "Stopped"
	StateSuspended VMState = "Suspended"
	StateStarting  VMState = "Starting"
	StateStopping  VMState = "Stopping"
	StateDeleted   VMState = "Deleted"
	StateUnknown   VMState = "Unknown"
)

// Instance represents a virtual machine instance
type Instance struct {
	Name   string `json:"name"`
	State  string `json:"state"`  // Running, Stopped, Suspended, etc.
	IPv4   string `json:"ipv4"`   // IP address
	Image  string `json:"image"`  // Source image
	CPUs   int    `json:"cpus"`   // Number of CPUs
	Memory string `json:"memory"` // Memory allocation
	Disk   string `json:"disk"`   // Disk size
}

// InstanceMetadata represents the metadata stored for each VM instance
type InstanceMetadata struct {
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	CPUs      int       `json:"cpus"`
	Memory    string    `json:"memory"`
	Disk      string    `json:"disk"`
	DiskPath  string    `json:"disk_path"`
	CloudInit string    `json:"cloud_init,omitempty"`
	State     string    `json:"state"`
	PID       int       `json:"pid,omitempty"`
	IPv4      string    `json:"ipv4,omitempty"`
	MAC       string    `json:"mac,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// QEMUProcess represents a running QEMU process
type QEMUProcess struct {
	PID     int    `json:"pid"`
	Name    string `json:"name"`
	Command string `json:"command"`
}

// NetworkInterface represents a network interface configuration
type NetworkInterface struct {
	Name    string `json:"name"`
	Type    string `json:"type"`    // user, bridge, tap
	Bridge  string `json:"bridge"`  // bridge name if type is bridge
	MAC     string `json:"mac"`     // MAC address
	IPv4    string `json:"ipv4"`    // IPv4 address
	Gateway string `json:"gateway"` // Gateway IP
}

// DiskConfig represents disk configuration
type DiskConfig struct {
	Path   string `json:"path"`   // Path to disk image
	Format string `json:"format"` // qcow2, raw, etc.
	Size   string `json:"size"`   // Disk size
	Bus    string `json:"bus"`    // virtio, ide, scsi
}

// CloudInitConfig represents cloud-init configuration
type CloudInitConfig struct {
	UserData     string            `json:"user_data,omitempty"`
	MetaData     string            `json:"meta_data,omitempty"`
	NetworkData  string            `json:"network_data,omitempty"`
	SSHKeys      []string          `json:"ssh_keys,omitempty"`
	Packages     []string          `json:"packages,omitempty"`
	RunCommands  []string          `json:"run_commands,omitempty"`
	WriteFiles   []CloudInitFile   `json:"write_files,omitempty"`
	Users        []CloudInitUser   `json:"users,omitempty"`
}

// CloudInitFile represents a file to be written by cloud-init
type CloudInitFile struct {
	Path        string `json:"path"`
	Content     string `json:"content"`
	Permissions string `json:"permissions,omitempty"`
	Owner       string `json:"owner,omitempty"`
}

// CloudInitUser represents a user to be created by cloud-init
type CloudInitUser struct {
	Name              string   `json:"name"`
	SSHAuthorizedKeys []string `json:"ssh_authorized_keys,omitempty"`
	Sudo              string   `json:"sudo,omitempty"`
	Shell             string   `json:"shell,omitempty"`
	Groups            []string `json:"groups,omitempty"`
}
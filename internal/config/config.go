package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// Config holds the application configuration
type Config struct {
	// Directories
	DataDir      string `yaml:"data_dir"`
	InstancesDir string `yaml:"instances_dir"`
	ImagesDir    string `yaml:"images_dir"`
	KeysDir      string `yaml:"keys_dir"`

	// KVM/QEMU settings
	QEMUBinary    string `yaml:"qemu_binary"`
	QEMUImgBinary string `yaml:"qemu_img_binary"`
	BridgeName    string `yaml:"bridge_name"`

	// SSH settings
	SSHKeyPath   string `yaml:"ssh_key_path"`
	SSHUser      string `yaml:"ssh_user"`
	SSHPort      int    `yaml:"ssh_port"`
	SSHTimeout   int    `yaml:"ssh_timeout"`

	// Default VM settings
	DefaultCPUs   int    `yaml:"default_cpus"`
	DefaultMemory string `yaml:"default_memory"`
	DefaultDisk   string `yaml:"default_disk"`

	// Image repositories
	ImageRepositories map[string]string `yaml:"image_repositories"`
}

// Load loads the configuration from file or returns defaults
func Load() *Config {
	cfg := getDefaults()

	// TODO: Load from config file if exists
	// For now, just return defaults

	return cfg
}

// getDefaults returns the default configuration
func getDefaults() *Config {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".slackpass")

	// Ensure data directory exists
	os.MkdirAll(dataDir, 0755)
	os.MkdirAll(filepath.Join(dataDir, "instances"), 0755)
	os.MkdirAll(filepath.Join(dataDir, "images"), 0755)
	os.MkdirAll(filepath.Join(dataDir, "keys"), 0755)

	cfg := &Config{
		// Directories
		DataDir:      dataDir,
		InstancesDir: filepath.Join(dataDir, "instances"),
		ImagesDir:    filepath.Join(dataDir, "images"),
		KeysDir:      filepath.Join(dataDir, "keys"),

		// KVM/QEMU settings
		QEMUBinary:    getQEMUBinary(),
		QEMUImgBinary: getQEMUImgBinary(),
		BridgeName:    "slackpass0",

		// SSH settings
		SSHKeyPath: filepath.Join(dataDir, "keys", "slackpass_rsa"),
		SSHUser:    "ubuntu", // Default cloud user
		SSHPort:    22,
		SSHTimeout: 30,

		// Default VM settings
		DefaultCPUs:   1,
		DefaultMemory: "1G",
		DefaultDisk:   "10G",

		// Image repositories
		ImageRepositories: getDefaultRepositories(),
	}

	return cfg
}

// getQEMUBinary returns the path to the QEMU binary
func getQEMUBinary() string {
	switch runtime.GOOS {
	case "linux":
		// Try common locations
		paths := []string{
			"/usr/bin/qemu-system-x86_64",
			"/usr/local/bin/qemu-system-x86_64",
			"/opt/qemu/bin/qemu-system-x86_64",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
		return "qemu-system-x86_64" // Fallback to PATH
	default:
		return "qemu-system-x86_64"
	}
}

// getQEMUImgBinary returns the path to the qemu-img binary
func getQEMUImgBinary() string {
	switch runtime.GOOS {
	case "linux":
		paths := []string{
			"/usr/bin/qemu-img",
			"/usr/local/bin/qemu-img",
			"/opt/qemu/bin/qemu-img",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
		return "qemu-img"
	default:
		return "qemu-img"
	}
}

// getDefaultRepositories returns the default image repositories
func getDefaultRepositories() map[string]string {
	return map[string]string{
		"debian":     "https://cloud.debian.org/images/cloud",
		"fedora":     "https://download.fedoraproject.org/pub/fedora/linux/releases",
		"almalinux":  "https://repo.almalinux.org/almalinux",
		"rockylinux": "https://download.rockylinux.org/pub/rocky",
		"centos":     "https://cloud.centos.org/centos",
		"opensuse":   "https://download.opensuse.org/repositories/Cloud:/Images",
		"gentoo":     "https://bouncer.gentoo.org/fetch/root/all/releases",
		"slackware":  "https://mirrors.slackware.com/slackware",
	}
}
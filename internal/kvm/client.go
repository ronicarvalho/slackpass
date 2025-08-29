package kvm

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/slackpass/slackpass/internal/config"
)

// Client handles KVM/QEMU operations
type Client struct {
	config *config.Config
}

// NewClient creates a new KVM client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
	}
}

// VMConfig represents the configuration for a virtual machine
type VMConfig struct {
	Name      string
	ImagePath string
	CPUs      int
	Memory    string
	Disk      string
	CloudInit string
}

// Create creates a new virtual machine
func (c *Client) Create(config *VMConfig) error {
	instanceDir := filepath.Join(c.config.InstancesDir, config.Name)

	// Create disk image
	diskPath := filepath.Join(instanceDir, "disk.qcow2")
	if err := c.createDiskImage(config.ImagePath, diskPath, config.Disk); err != nil {
		return fmt.Errorf("failed to create disk image: %w", err)
	}

	// Generate cloud-init ISO if needed
	var cloudInitPath string
	if config.CloudInit != "" {
		cloudInitPath = filepath.Join(instanceDir, "cloud-init.iso")
		if err := c.createCloudInitISO(config.CloudInit, cloudInitPath); err != nil {
			return fmt.Errorf("failed to create cloud-init ISO: %w", err)
		}
	}

	// Save instance metadata
	metadata := &InstanceMetadata{
		Name:      config.Name,
		Image:     config.ImagePath,
		CPUs:      config.CPUs,
		Memory:    config.Memory,
		Disk:      config.Disk,
		DiskPath:  diskPath,
		CloudInit: cloudInitPath,
		CreatedAt: time.Now(),
		State:     string(StateStopped),
	}

	metadataPath := filepath.Join(instanceDir, "metadata.json")
	return c.saveMetadata(metadata, metadataPath)
}

// Start starts a virtual machine
func (c *Client) Start(name string) error {
	metadata, err := c.loadMetadata(name)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	if metadata.State == string(StateRunning) {
		return fmt.Errorf("instance '%s' is already running", name)
	}

	// Build QEMU command
	cmd := c.buildQEMUCommand(metadata)

	// Start the VM
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start VM: %w", err)
	}

	// Update state
	metadata.State = string(StateRunning)
	metadata.PID = cmd.Process.Pid
	metadataPath := filepath.Join(c.config.InstancesDir, name, "metadata.json")
	return c.saveMetadata(metadata, metadataPath)
}

// Stop stops a virtual machine
func (c *Client) Stop(name string, force bool) error {
	metadata, err := c.loadMetadata(name)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	if metadata.State != string(StateRunning) {
		return fmt.Errorf("instance '%s' is not running", name)
	}

	// Send shutdown signal to QEMU process
	if metadata.PID > 0 {
		process, err := os.FindProcess(metadata.PID)
		if err == nil {
			if force {
				process.Kill()
			} else {
				process.Signal(os.Interrupt)
			}
		}
	}

	// Update state
	metadata.State = string(StateStopped)
	metadata.PID = 0
	metadataPath := filepath.Join(c.config.InstancesDir, name, "metadata.json")
	return c.saveMetadata(metadata, metadataPath)
}

// Delete deletes a virtual machine
func (c *Client) Delete(name string, purge, force bool) error {
	// Stop the VM first if running
	metadata, err := c.loadMetadata(name)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	if metadata.State == string(StateRunning) {
		if err := c.Stop(name, force); err != nil {
			return fmt.Errorf("failed to stop VM: %w", err)
		}
	}

	// Remove instance directory
	instanceDir := filepath.Join(c.config.InstancesDir, name)
	return os.RemoveAll(instanceDir)
}

// List returns all virtual machine instances
func (c *Client) List() ([]*Instance, error) {
	instances := make([]*Instance, 0)

	entries, err := os.ReadDir(c.config.InstancesDir)
	if err != nil {
		return instances, nil // Return empty list if directory doesn't exist
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		metadata, err := c.loadMetadata(entry.Name())
		if err != nil {
			continue // Skip invalid instances
		}

		instance := &Instance{
			Name:   metadata.Name,
			State:  metadata.State,
			IPv4:   metadata.IPv4,
			Image:  metadata.Image,
			CPUs:   metadata.CPUs,
			Memory: metadata.Memory,
			Disk:   metadata.Disk,
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

// Info displays detailed information about a virtual machine
func (c *Client) Info(name string) error {
	metadata, err := c.loadMetadata(name)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	fmt.Printf("Name:           %s\n", metadata.Name)
	fmt.Printf("State:          %s\n", metadata.State)
	fmt.Printf("IPv4:           %s\n", metadata.IPv4)
	fmt.Printf("Image:          %s\n", metadata.Image)
	fmt.Printf("CPUs:           %d\n", metadata.CPUs)
	fmt.Printf("Memory:         %s\n", metadata.Memory)
	fmt.Printf("Disk:           %s\n", metadata.Disk)
	fmt.Printf("Created:        %s\n", metadata.CreatedAt.Format(time.RFC3339))

	return nil
}

// Helper methods

func (c *Client) createDiskImage(sourcePath, targetPath, size string) error {
	// For now, create an empty disk image
	// TODO: Implement proper image copying/resizing
	cmd := exec.Command(c.config.QEMUImgBinary, "create", "-f", "qcow2", targetPath, size)
	return cmd.Run()
}

func (c *Client) createCloudInitISO(configPath, isoPath string) error {
	// TODO: Implement cloud-init ISO creation
	return nil
}

func (c *Client) buildQEMUCommand(metadata *InstanceMetadata) *exec.Cmd {
	args := []string{
		"-name", metadata.Name,
		"-machine", "type=pc,accel=kvm",
		"-cpu", "host",
		"-smp", strconv.Itoa(metadata.CPUs),
		"-m", metadata.Memory,
		"-drive", fmt.Sprintf("file=%s,format=qcow2,if=virtio", metadata.DiskPath),
		"-netdev", "user,id=net0,hostfwd=tcp::0-:22",
		"-device", "virtio-net-pci,netdev=net0",
		"-nographic",
		"-daemonize",
	}

	if metadata.CloudInit != "" {
		args = append(args, "-drive", fmt.Sprintf("file=%s,format=raw,if=virtio", metadata.CloudInit))
	}

	return exec.Command(c.config.QEMUBinary, args...)
}

func (c *Client) loadMetadata(name string) (*InstanceMetadata, error) {
	metadataPath := filepath.Join(c.config.InstancesDir, name, "metadata.json")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	var metadata InstanceMetadata
	err = json.Unmarshal(data, &metadata)
	return &metadata, err
}

func (c *Client) saveMetadata(metadata *InstanceMetadata, path string) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
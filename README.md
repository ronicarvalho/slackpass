# Slackpass

Slackpass is a virtual machine manager for Linux distributions, inspired by Multipass. It creates and manages Linux virtual machines using KVM/QEMU with a simple command-line interface.

## Features

- **Multiple Linux Distributions**: Support for Slackware, Debian, Fedora, AlmaLinux, Rocky Linux, CentOS, Gentoo, and openSUSE
- **KVM/QEMU Backend**: Uses KVM virtualization with QEMU for high performance
- **SSH Integration**: Automatic SSH access to virtual machines
- **Cloud Images**: Downloads and uses official cloud images
- **Simple CLI**: Easy-to-use command-line interface similar to Multipass

## Supported Distributions

- **Debian**: bookworm (12), bullseye (11)
- **Fedora**: 39, 38, 37
- **AlmaLinux**: 9, 8
- **Rocky Linux**: 9, 8
- **CentOS**: Stream 9, Stream 8
- **openSUSE**: Tumbleweed, Leap 15.5
- **Gentoo**: Latest
- **Slackware**: 15.0, 14.2

## Prerequisites

- Linux operating system
- KVM virtualization support
- QEMU installed (`qemu-system-x86_64`, `qemu-img`)
- SSH client
- Go 1.21+ (for building from source)

### Installing Prerequisites (Ubuntu/Debian)

```bash
sudo apt update
sudo apt install qemu-kvm qemu-utils libvirt-daemon-system libvirt-clients bridge-utils
sudo usermod -aG kvm,libvirt $USER
```

### Installing Prerequisites (Fedora/CentOS/RHEL)

```bash
sudo dnf install qemu-kvm qemu-img libvirt libvirt-daemon-config-network libvirt-daemon-kvm
sudo usermod -aG kvm,libvirt $USER
```

## Installation

### From Source

```bash
git clone https://github.com/slackpass/slackpass.git
cd slackpass
go build -o slackpass
sudo mv slackpass /usr/local/bin/
```

## Quick Start

1. **Launch a virtual machine**:
   ```bash
   slackpass launch debian
   ```

2. **List running instances**:
   ```bash
   slackpass list
   ```

3. **Open a shell session**:
   ```bash
   slackpass shell butterfly-effect
   ```

4. **Execute a command**:
   ```bash
   slackpass exec butterfly-effect -- sudo apt update
   ```

5. **Stop and delete an instance**:
   ```bash
   slackpass stop butterfly-effect
   slackpass delete butterfly-effect
   ```

## Commands

### Core Commands

- `slackpass launch [image] [name]` - Create and start a new virtual machine
- `slackpass list` - List all virtual machine instances
- `slackpass shell [name]` - Open a shell session to a virtual machine
- `slackpass exec [name] -- [command]` - Execute a command on a virtual machine
- `slackpass info [name]` - Display detailed information about instances
- `slackpass start [name...]` - Start virtual machine instances
- `slackpass stop [name...]` - Stop virtual machine instances
- `slackpass delete [name...]` - Delete virtual machine instances

### Image Management

- `slackpass find [image-name]` - Display available images to launch

### Examples

```bash
# Launch different distributions
slackpass launch debian:bookworm mydebian
slackpass launch fedora:39 myfedora --cpus 2 --memory 2G
slackpass launch almalinux:9 myalma --disk 20G

# Find available images
slackpass find
slackpass find debian

# Manage instances
slackpass list
slackpass info mydebian
slackpass shell mydebian
slackpass exec mydebian -- ls -la /home

# Stop and clean up
slackpass stop mydebian myfedora
slackpass delete mydebian --purge
```

## Configuration

Slackpass stores its configuration and data in `~/.slackpass/`:

- `~/.slackpass/instances/` - Virtual machine instances
- `~/.slackpass/images/` - Downloaded cloud images
- `~/.slackpass/keys/` - SSH keys

## Architecture

Slackpass is built with a modular architecture:

- **CLI Layer**: Cobra-based command-line interface
- **VM Manager**: High-level virtual machine operations
- **KVM Client**: Low-level QEMU/KVM integration
- **SSH Client**: SSH connection and command execution
- **Image Manager**: Cloud image discovery and management
- **Configuration**: Application settings and defaults

## Development

### Building

```bash
go mod download
go build -o slackpass
```

### Testing

```bash
go test ./...
```

### Project Structure

```
slackpass/
├── main.go                # Application entry point
├── cmd/                   # CLI commands
│   ├── root.go            # Root command and configuration
│   ├── launch.go          # Launch command
│   ├── list.go            # List command
│   ├── shell.go           # Shell command
│   ├── exec.go            # Exec command
│   ├── info.go            # Info command
│   ├── start.go           # Start/Stop commands
│   ├── delete.go          # Delete command
│   └── find.go            # Find command
├── internal/              # Internal packages
│   ├── vm/                # Virtual machine management
│   ├── kvm/               # KVM/QEMU integration
│   ├── ssh/               # SSH client
│   ├── images/            # Image management
│   └── config/            # Configuration
└── go.mod                 # Go module definition
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [Multipass](https://multipass.run/) by Canonical
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Uses [Viper](https://github.com/spf13/viper) for configuration management

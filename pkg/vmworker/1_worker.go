package vmworker

import (
	"fmt"

	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

// VmWorkerConfig holds configuration for VM worker connection.
//
// Driver specifies the virtualization driver to use.
// Allowed values: "qemu", "xen", "vbox", "test".
//
// Transport specifies the transport protocol for libvirt connection.
// Allowed values: "ssh", "tls", "tcp", "unix", "test".
//
// Hostname is required when using ssh, tls, or tcp transports.
// It specifies the remote host to connect to.
//
// Username is required only for ssh transport.
// It is the user used to login on the remote host.
//
// Port is optional and must be between 0 and 65535 if specified.
// It is used for tcp and tls transports.
//
// Scope is optional and specifies the libvirt connection scope.
// Allowed values: "system" or "session".
//
// SocketPath is optional and used when Transport is "unix".
// It specifies the path to the unix domain socket.
type VmWorkerConfig struct {
	log        *zap.Logger
	Driver     string // required, one of: qemu, xen, vbox, test
	Transport  string // required, one of: ssh, tls, tcp, unix, test
	Hostname   string // required if transport is ssh/tls/tcp
	Username   string // required if transport is ssh
	Port       int    // optional, must be 0-65535 if set
	Scope      string // optional, one of: system, session
	SocketPath string // optional, used if transport is unix

	client *libvirt.Connect
}

// NewVmWorkerConfig creates and validates a new VmWorkerConfig.
// It performs:
//   - Logger assignment
//   - Field population from parameters
//   - Validation of required combinations
//   - ICMP reachability check (warn-only)
//
// Note: It requires a zap.Logger for structured logging.
func NewVmWorkerConfig(
	logger *zap.Logger,
	driver string,
	transport string,
	hostname string,
	username string,
	port int,
	scope string,
	socketPath string,
) (*VmWorkerConfig, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	cfg := &VmWorkerConfig{
		log:        logger,
		Driver:     driver,
		Transport:  transport,
		Hostname:   hostname,
		Username:   username,
		Port:       port,
		Scope:      scope,
		SocketPath: socketPath,
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if warning := cfg.CheckReachability(); !warning {
		cfg.log.Warn("Unable to reach the host", zap.String("host", cfg.Hostname))
	}

	return cfg, nil
}

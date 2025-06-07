package vmworker

import (
	"errors"
	"fmt"
	"strings"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"go.uber.org/zap"
)

// CheckReachability performs an ICMP ping to the configured host to check reachability.
// Returns true if the host responds to ping, false otherwise.
//
// Logs a warning if:
// - The hostname is empty.
// - The ICMP pinger cannot be created.
// - The host does not respond to the ping request.
func (cfg *VmWorkerConfig) CheckReachability() bool {
	if strings.TrimSpace(cfg.Hostname) == "" {
		cfg.log.Warn(
			"Hostname is empty; skipping reachability check",
		)
		return false
	}

	ping, err := probing.NewPinger(cfg.Hostname)
	if err != nil {
		cfg.log.Warn(
			"Could not create ICMP ping client for host",
			zap.String("host", cfg.Hostname),
		)
		return false
	}

	ping.Count = 2
	ping.Timeout = 3 * time.Second
	ping.Interval = 500 * time.Millisecond
	ping.SetPrivileged(false)

	err = ping.Run()
	if err != nil {
		cfg.log.Warn(
			"ICMP ping to host failed",
			zap.String("host", cfg.Hostname),
		)
		return false
	}

	stats := ping.Statistics()
	if stats.PacketsRecv == 0 {
		cfg.log.Warn(
			"Host did not respond to ICMP ping (may be firewalled or offline)",
			zap.String("host", cfg.Hostname),
		)
		return false
	}

	return true
}

// Validate performs sanity and cross-field checks on VmWorkerConfig.
// Returns an error if validation fails, otherwise nil.
//
// Checks include:
// - Driver and Transport are among allowed values.
// - Hostname is required for ssh, tls, and tcp transports.
// - Username is required for ssh transport.
// - Port must be within valid range if specified.
// - Scope, if specified, must be either "system" or "session".
// - SocketPath is currently not enforced
func (cfg *VmWorkerConfig) Validate() error {
	// backup if i forgot (trim spaces)
	cfg.Driver = strings.TrimSpace(cfg.Driver)
	cfg.Transport = strings.TrimSpace(cfg.Transport)
	cfg.Hostname = strings.TrimSpace(cfg.Hostname)
	cfg.Username = strings.TrimSpace(cfg.Username)
	cfg.Scope = strings.TrimSpace(cfg.Scope)
	cfg.SocketPath = strings.TrimSpace(cfg.SocketPath)

	switch cfg.Driver {
	case "qemu", "xen", "vbox", "test":

	default:
		return fmt.Errorf("invalid driver: %s (allowed: qemu, xen, vbox, test)", cfg.Driver)
	}

	switch cfg.Transport {
	case "ssh", "tls", "tcp", "unix", "test":
	default:
		return fmt.Errorf("invalid transport: %s (allowed: ssh, tls, tcp, unix, test)", cfg.Transport)
	}

	if cfg.Transport == "ssh" || cfg.Transport == "tls" || cfg.Transport == "tcp" {
		if cfg.Hostname == "" {
			return fmt.Errorf("hostname required for transport %s", cfg.Transport)
		}
	}

	if cfg.Transport == "ssh" && cfg.Username == "" {
		return errors.New("username required for ssh transport")
	}

	if cfg.Port < 0 || cfg.Port > 65535 {
		return fmt.Errorf("port must be between 0 and 65535, got %d", cfg.Port)
	}

	if cfg.Scope != "" && cfg.Scope != "system" && cfg.Scope != "session" {
		return fmt.Errorf("invalid scope: %s (allowed: system, session)", cfg.Scope)
	}
	return nil
}

func defaultIfEmpty(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

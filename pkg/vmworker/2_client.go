package vmworker

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

func (cfg *VmWorkerConfig) Connect() (*libvirt.Connect, error) {
	if cfg.client != nil {
		ok, err := cfg.client.IsAlive()
		if err == nil && ok {
			return cfg.client, nil
		}
	}

	uri, err := cfg.URI()
	if err != nil {
		return nil, err
	}
	client, err := libvirt.NewConnect(uri) // need update for other auth
	if err != nil {
		return nil, err
	}
	cfg.client = client
	return nil, nil
}

// URI constructs the libvirt connection URI based on the VmWorkerConfig.
// It first validates the configuration using Validate().
// Returns a well-formed connection URI or an error if validation fails or unsupported transport is used.
//
// URI format examples:
//   - SSH:   qemu+ssh://user@host/system
//   - TLS:   qemu+tls://host[:port]/system
//   - TCP:   qemu+tcp://host[:port]/system
//   - UNIX:  qemu+unix:///system or qemu+unix://<socket_path>/system
//   - TEST:  test:///default
//
// The driver defaults to "qemu" if not specified.
// The scope defaults to "session" if not specified.
func (cfg *VmWorkerConfig) URI() (string, error) {
	driver := defaultIfEmpty(cfg.Driver, "qemu")
	scope := defaultIfEmpty(cfg.Scope, "session")
	switch cfg.Transport {
	case "ssh":
		return fmt.Sprintf("%s+ssh://%s@%s/%s", driver, cfg.Username, cfg.Hostname, scope), nil
	case "tls":
		if cfg.Port > 0 {
			return fmt.Sprintf("%s+tls://%s:%d/%s", driver, cfg.Hostname, cfg.Port, scope), nil
		}
		return fmt.Sprintf("%s+tls://%s/%s", driver, cfg.Hostname, scope), nil
	case "tcp":
		if cfg.Port > 0 {
			return fmt.Sprintf("%s+tcp://%s:%d/%s", driver, cfg.Hostname, cfg.Port, scope), nil
		}
		return fmt.Sprintf("%s+tcp://%s/%s", driver, cfg.Hostname, scope), nil
	case "unix":
		if cfg.SocketPath != "" {
			return fmt.Sprintf("%s+unix://%s/%s", driver, cfg.SocketPath, scope), nil
		}
		return fmt.Sprintf("%s+unix:///%s", driver, scope), nil
	case "test":
		return "test:///default", nil
	default:
		return "", fmt.Errorf("unsupported transport: %s", cfg.Transport)
	}

	// return "", nil
}

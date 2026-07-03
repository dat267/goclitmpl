package diagnose

import (
	"log/slog"
	"net"
	"time"
)

// CheckAddress validates connection reachability to a specific address within the given timeout.
func CheckAddress(address string, timeout time.Duration) error {
	slog.Debug("checking network reachability", slog.String("address", address), slog.Duration("timeout", timeout))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

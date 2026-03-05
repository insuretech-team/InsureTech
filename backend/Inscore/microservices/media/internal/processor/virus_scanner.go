package processor

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// VirusScanner performs virus scanning using ClamAV via TCP socket.
// Uses the clamd protocol to communicate with a running ClamAV daemon.
// No CGO required - communicates via standard TCP sockets.
type VirusScanner struct {
	enabled      bool
	clamavAddr   string
}

// NewVirusScanner creates a new virus scanner.
// enabled: set to true to enable scanning against ClamAV daemon
// clamavAddr: address of ClamAV daemon, e.g., "localhost:3310"
func NewVirusScanner(enabled bool, clamavAddr string) *VirusScanner {
	return &VirusScanner{
		enabled:    enabled,
		clamavAddr: clamavAddr,
	}
}

// Scan scans data for viruses using ClamAV.
// If enabled=false, returns clean=true as a graceful no-op.
// If enabled=true, connects to ClamAV via TCP and streams data using INSTREAM protocol.
// Returns: clean (bool), virusName (string if infected), error
func (vs *VirusScanner) Scan(ctx context.Context, data []byte) (clean bool, virusName string, err error) {
	if !vs.enabled {
		return true, "", nil
	}

	if vs.clamavAddr == "" {
		return false, "", fmt.Errorf("ClamAV address not configured")
	}

	// Establish connection to ClamAV daemon
	conn, err := net.DialTimeout("tcp", vs.clamavAddr, 0)
	if err != nil {
		return false, "", fmt.Errorf("failed to connect to ClamAV: %w", err)
	}
	defer conn.Close()

	// Handle context cancellation
	if ctx.Err() != nil {
		return false, "", ctx.Err()
	}

	// Send INSTREAM command
	if _, err := conn.Write([]byte("zINSTREAM\x00")); err != nil {
		return false, "", fmt.Errorf("failed to send INSTREAM command: %w", err)
	}

	// Stream data in chunks (default 4KB chunks)
	chunkSize := 4096
	offset := 0
	for offset < len(data) {
		end := offset + chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunk := data[offset:end]
		chunkLen := uint32(len(chunk))

		// Send 4-byte big-endian length + chunk data
		lenBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(lenBuf, chunkLen)

		if _, err := conn.Write(lenBuf); err != nil {
			return false, "", fmt.Errorf("failed to send chunk length: %w", err)
		}

		if _, err := conn.Write(chunk); err != nil {
			return false, "", fmt.Errorf("failed to send chunk data: %w", err)
		}

		offset = end
	}

	// Send 4-byte zero to end stream
	endBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(endBuf, 0)
	if _, err := conn.Write(endBuf); err != nil {
		return false, "", fmt.Errorf("failed to send stream terminator: %w", err)
	}

	// Read response line
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, "", fmt.Errorf("failed to read ClamAV response: %w", err)
	}

	response = strings.TrimSpace(response)

	// Parse clamd response
	// Format: 'stream: OK' = clean
	// Format: 'stream: VIRUS_NAME FOUND' = infected
	if strings.Contains(response, "OK") {
		return true, "", nil
	}

	if strings.Contains(response, "FOUND") {
		parts := strings.Split(response, ":")
		if len(parts) >= 2 {
			virusPart := strings.TrimSpace(parts[1])
			virusName = strings.TrimSuffix(virusPart, "FOUND")
			virusName = strings.TrimSpace(virusName)
		}
		return false, virusName, nil
	}

	// Unknown response format
	return false, "", fmt.Errorf("unexpected ClamAV response: %s", response)
}

// IsEnabled returns whether virus scanning is enabled.
func (vs *VirusScanner) IsEnabled() bool {
	return vs.enabled
}

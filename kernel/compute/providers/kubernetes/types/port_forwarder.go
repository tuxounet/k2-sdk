package types

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type PortForwarder struct {
	mu          sync.Mutex
	Record      PortsForwardRecord
	mounting    bool
	mounted     bool
	kubeConfig  string
	hostAddress string
	log         runtimeTypes.ILogger
	cmd         *exec.Cmd
	cancel      context.CancelFunc
}

func NewPortForwarder(record PortsForwardRecord, kubeConfig string, parentLog runtimeTypes.ILogger, hostAddress string) *PortForwarder {

	subLogger := parentLog.CreateSubLogger("PortForwarder/" + record.ServiceName)
	return &PortForwarder{
		Record:      record,
		kubeConfig:  kubeConfig,
		log:         subLogger,
		hostAddress: hostAddress,
		mounting:    false,
		mounted:     false,
	}
}

func (p *PortForwarder) IsReady() bool {
	if p.mounting {
		return false
	}
	if !p.mounted {
		return false
	}
	return true
}

func (p *PortForwarder) ForwardRequest(c *gin.Context) error {
	if p.mounting {
		p.log.DebugF("Port forwarder is mounting, please wait...")
		c.Status(http.StatusBadGateway)
		return nil
	}

	if !p.mounted {
		if err := p.Mount(); err != nil {
			p.log.ErrorF("failed to mount port forwarder: %v", err)
			c.Status(http.StatusBadGateway)
			return nil
		}
	}

	// Create a reverse proxy targeting the locally forwarded port
	targetURL, err := url.Parse(fmt.Sprintf("http://%s:%d", p.hostAddress, p.Record.LocalPort))
	if err != nil {
		p.log.ErrorF("Failed to parse target URL: %v", err)
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("invalid proxy target: %w", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if err != http.ErrAbortHandler {
			p.log.ErrorF("proxy error: %v", err)
		}
	}

	proxy.Director = func(r *http.Request) {
		r.URL.Scheme = targetURL.Scheme
		r.URL.Host = targetURL.Host
		r.Host = "kube.k2"
	}

	// Recover from http.ErrAbortHandler panics during SSE/long-polling
	defer func() {
		if err := recover(); err != nil {
			if err != http.ErrAbortHandler {
				panic(err)
			}
		}
	}()

	proxy.ServeHTTP(c.Writer, c.Request)

	return nil
}

func (p *PortForwarder) Mount() error {
	if p.mounted {
		p.log.DebugF("Port forwarder is already mounted")
		return nil
	}
	if p.mounting {
		return fmt.Errorf("port forwarder is already mounting")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.mounting = true
	p.log.DebugF("Mounting port forwarder for %s/%s:%d -> localhost:%d",
		p.Record.ServiceNamespace, p.Record.ServiceName,
		p.Record.ServicePort, p.Record.LocalPort)

	ctx, cancel := context.WithCancel(context.Background())

	args := []string{
		"port-forward",
		fmt.Sprintf("service/%s", p.Record.ServiceName),
		fmt.Sprintf("%d:%d", p.Record.LocalPort, p.Record.ServicePort),
		"-n", p.Record.ServiceNamespace,
	}
	if p.kubeConfig != "" {
		args = append(args, "--kubeconfig", p.kubeConfig)
	}

	cmd := exec.CommandContext(ctx, "kubectl", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		p.mounting = false
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		p.mounting = false
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		cancel()
		p.mounting = false
		return fmt.Errorf("failed to start kubectl port-forward: %w", err)
	}

	// Detect errors from stderr
	errCh := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			p.log.DebugF("kubectl stderr: %s", line)
			if strings.Contains(strings.ToLower(line), "error") {
				errCh <- fmt.Errorf("kubectl: %s", line)
				return
			}
		}
	}()

	// Wait for "Forwarding from" on stdout to know port-forward is ready
	readyCh := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			p.log.DebugF("kubectl: %s", line)
			if strings.Contains(line, "Forwarding from") {
				close(readyCh)
			}
		}
	}()

	// Wait for ready signal, error, or timeout
	select {
	case <-readyCh:
		p.log.InfoF("port forwarding established: %s:%d -> %s/%s:%d",
			p.hostAddress, p.Record.LocalPort,
			p.Record.ServiceNamespace, p.Record.ServiceName, p.Record.ServicePort)
		p.mounted = true
		p.mounting = false
		p.cmd = cmd
		p.cancel = cancel
	case err := <-errCh:
		cancel()
		p.mounting = false
		return fmt.Errorf("kubectl port-forward failed: %w", err)
	case <-time.After(30 * time.Second):
		cancel()
		p.mounting = false
		return fmt.Errorf("timeout waiting for kubectl port-forward to be ready")
	}

	return nil
}

func (p *PortForwarder) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}
	if p.cmd != nil && p.cmd.Process != nil {
		// Give kubectl a chance to exit cleanly, then force kill
		done := make(chan error, 1)
		go func() {
			done <- p.cmd.Wait()
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			p.cmd.Process.Kill()
		}
		p.cmd = nil
	}
	p.mounted = false
	p.mounting = false
	p.log.DebugF("port forwarding stopped: %s:%d -> %s/%s:%d",
		p.hostAddress, p.Record.LocalPort,
		p.Record.ServiceNamespace, p.Record.ServiceName, p.Record.ServicePort)
	return nil
}

package types

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type PortForwarder struct {
	Record      PortsForwardRecord
	mounting    bool
	mounted     bool
	kubeConfig  string
	portForward *portforward.PortForwarder
	stopchan    chan struct{}
	hostAddress string
	log         runtimeTypes.ILogger
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
		log.Fatalf("Failed to parse target URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Reset the request URL path to remove the /proxy prefix
	c.Request.URL.Path = c.Param("any")
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
	p.mounting = true
	p.log.DebugF("Mounting port forwarder for %s:%d", p.Record.ServiceName, p.Record.ServicePort)

	// Build the configuration from kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", p.kubeConfig)
	if err != nil {
		p.mounting = false
		return fmt.Errorf("failed to build kubeconfig: %v", err)
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		p.mounting = false
		return fmt.Errorf("failed to create clientset: %v", err)

	}

	// Retrieve the service to obtain its label selector
	svc, err := clientset.CoreV1().Services(p.Record.ServiceNamespace).Get(context.TODO(), p.Record.ServiceName, metav1.GetOptions{})
	if err != nil {
		p.mounting = false
		return fmt.Errorf("failed to get service: %v", err)
	}
	// List pods matching the service's selector
	selector := labels.Set(svc.Spec.Selector)
	listOpts := metav1.ListOptions{
		LabelSelector: selector.AsSelector().String(),
	}
	podList, err := clientset.CoreV1().Pods(p.Record.ServiceNamespace).List(context.TODO(), listOpts)
	if err != nil {
		p.mounting = false
		return fmt.Errorf("failed to list pods: %v", err)
	}
	if len(podList.Items) == 0 {
		p.mounting = false
		return fmt.Errorf("no pods found for service: %q", p.Record.ServiceName)
	}

	// Choose the first running pod
	var podName string
	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodRunning {
			podName = pod.Name
			break
		}
	}
	if podName == "" {

		p.mounting = false
		return fmt.Errorf("no running pod found for service: %q", p.Record.ServiceName)
	}
	p.log.DebugF("Selected pod %q for port-forwarding", podName)

	// Setup port-forwarding using client-go's tools
	stopChan := make(chan struct{}, 1)
	readyChan := make(chan struct{})

	// Build the request URL for port-forwarding to the pod
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(p.Record.ServiceNamespace).
		Name(podName).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {

		p.mounting = false
		return fmt.Errorf("failed to create round tripper: %v", err)
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", req.URL())

	// Format: localPort:remotePort (e.g., "8080:80")
	ports := []string{fmt.Sprintf("%d:%d", p.Record.LocalPort, p.Record.ServicePort)}

	pf, err := portforward.New(dialer, ports, stopChan, readyChan, os.Stdout, os.Stderr)
	if err != nil {

		p.mounting = false
		return fmt.Errorf("failed to create port forwarder: %v", err)
	}

	p.portForward = pf
	p.stopchan = stopChan

	// Start the port-forwarding in a goroutine
	go func() {

		if err := pf.ForwardPorts(); err != nil {
			p.log.ErrorF("port forwarding failed: %v", err)
			p.mounting = false
			return
		}
	}()

	// Wait until port-forwarding is ready (or timeout after 10 seconds)
	select {
	case <-readyChan:
		p.log.InfoF("port forwarding established: %s:%d -> pod:%d\n", p.hostAddress, p.Record.LocalPort, p.Record.ServicePort)
		p.mounted = true
		p.mounting = false
	case <-time.After(10 * time.Second):
		close(stopChan)
		p.mounted = false
		p.mounting = false
		return fmt.Errorf("timeout waiting for port-forward to be ready")
	}

	return nil
}

func (p *PortForwarder) Stop() error {

	if p.stopchan != nil {
		close(p.stopchan)
		p.stopchan = nil
	}
	if p.portForward != nil {
		p.portForward.Close()
		p.portForward = nil
	}
	p.mounted = false
	p.mounting = false
	p.log.DebugF("port forwarding stopped: %s:%d -> pod:%d\n", p.hostAddress, p.Record.LocalPort, p.Record.ServicePort)
	return nil
}

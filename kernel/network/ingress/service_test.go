package ingress

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	"github.com/tuxounet/k2-sdk/kernel/storage/paths"
	"github.com/tuxounet/k2-sdk/testutils"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

func newTestIngressService(t *testing.T) *Service {
	t.Helper()
	kernel := testutils.NewMockKernel("test-app", "1.0.0")

	app := testutils.NewMockApp("test-app", "1.0.0")
	kernel.SetApp(app)

	// Create config service with needed ingress config
	configSvc := config.NewService(kernel).(*config.Service)
	kernel.SetService(configSvc)
	configSvc.SetData("records", map[string]any{
		"host": map[string]any{
			"ingress": map[string]any{
				"root":    "http://localhost:8080",
				"address": "0.0.0.0",
				"port":    8080,
				"auth": map[string]any{
					"default_access": "public",
				},
			},
		},
	})

	// Create paths service
	pathsSvc := paths.NewService(kernel)
	kernel.SetService(pathsSvc)

	svc := NewService(kernel).(*Service)
	return svc
}

func TestIngressService_NewService(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel)

	assert.Equal(t, ServiceKey, svc.GetName())
	assert.Equal(t, kernel, svc.GetKernel())
}

func TestIngressService_GetServer_NoInit(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*Service)

	server := svc.GetServer()
	assert.Nil(t, server)
}

func TestIngressService_GetRouter_NoInit(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*Service)

	router := svc.GetRouter()
	assert.Nil(t, router)
}

func TestIngressService_RegisterIngress(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*Service)

	ingress := &types.IngressDefinition{
		IngressPath:  "/api/v1",
		AccessPolicy: runtimeTypes.AccessPolicyPublic,
		ServicePort:  3000,
		ServiceHost:  "localhost",
	}
	err := svc.RegisterIngress(ingress)
	require.NoError(t, err)

	records := svc.getIngressesRecords()
	assert.Len(t, records, 1)
	assert.Equal(t, "/api/v1", records[0].IngressPath)
}

func TestIngressService_RegisterIngress_Multiple(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*Service)

	for i, path := range []string{"/api/v1", "/api/v2", "/health"} {
		err := svc.RegisterIngress(&types.IngressDefinition{
			IngressPath: path,
			ServicePort: 3000 + i,
			ServiceHost: "localhost",
		})
		require.NoError(t, err)
	}

	records := svc.getIngressesRecords()
	assert.Len(t, records, 3)
}

func TestIngressService_Register_CreatesServer(t *testing.T) {
	svc := newTestIngressService(t)

	err := svc.Register()
	require.NoError(t, err)

	server := svc.GetServer()
	assert.NotNil(t, server)
}

func TestIngressService_Register_WithComponents(t *testing.T) {
	kernel := testutils.NewMockKernel("test-app", "1.0.0")

	app := testutils.NewMockApp("test-app", "1.0.0")
	comp := testutils.NewMockComponent(app, "api", 0)
	ctrl := testutils.NewMockController(comp, "users", 0)
	comp.SetControllers([]runtimeTypes.IAppController{ctrl})
	app.SetComponents([]runtimeTypes.IAppComponent{comp})
	kernel.SetApp(app)

	configSvc := config.NewService(kernel).(*config.Service)
	kernel.SetService(configSvc)
	configSvc.SetData("records", map[string]any{
		"host": map[string]any{
			"ingress": map[string]any{
				"root":    "http://localhost:8080",
				"address": "0.0.0.0",
				"port":    8080,
				"auth": map[string]any{
					"default_access": "public",
				},
			},
		},
	})

	pathsSvc := paths.NewService(kernel)
	kernel.SetService(pathsSvc)

	svc := NewService(kernel).(*Service)

	err := svc.Register()
	require.NoError(t, err)

	server := svc.GetServer()
	assert.NotNil(t, server)
}

func TestIngressService_Register_AuthenticatedAccess(t *testing.T) {
	kernel := testutils.NewMockKernel("test-app", "1.0.0")

	app := testutils.NewMockApp("test-app", "1.0.0")
	kernel.SetApp(app)

	configSvc := config.NewService(kernel).(*config.Service)
	kernel.SetService(configSvc)
	configSvc.SetData("records", map[string]any{
		"host": map[string]any{
			"ingress": map[string]any{
				"root":    "http://localhost:8080",
				"address": "0.0.0.0",
				"port":    8080,
				"auth": map[string]any{
					"default_access": "authenticated",
				},
			},
		},
	})

	pathsSvc := paths.NewService(kernel)
	kernel.SetService(pathsSvc)

	svc := NewService(kernel).(*Service)

	err := svc.Register()
	require.NoError(t, err)
	assert.NotNil(t, svc.GetServer())
}

func TestIngressService_Stop_WithServer(t *testing.T) {
	svc := newTestIngressService(t)

	err := svc.Register()
	require.NoError(t, err)

	err = svc.Stop()
	require.NoError(t, err)
}

func TestIngressService_IngressDefinition(t *testing.T) {
	def := types.IngressDefinition{
		IngressPath:  "/api",
		AccessPolicy: runtimeTypes.AccessPolicyAuthenticated,
		ServicePort:  9090,
		ServiceHost:  "backend",
	}

	assert.Equal(t, "/api", def.IngressPath)
	assert.Equal(t, runtimeTypes.AccessPolicyAuthenticated, def.AccessPolicy)
	assert.Equal(t, 9090, def.ServicePort)
	assert.Equal(t, "backend", def.ServiceHost)
}

func TestIngressService_IngressDefinition_WithRewrite(t *testing.T) {
	rewrite := "/v2"
	def := types.IngressDefinition{
		IngressPath: "/api",
		RewritePath: &rewrite,
	}

	assert.Equal(t, "/v2", *def.RewritePath)
}

func TestIngressService_Lifecycle(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel)

	assert.NoError(t, svc.Start())
}

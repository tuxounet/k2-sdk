package testutils

import (
	"context"
	"embed"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/swag"
	"github.com/tuxounet/k2-sdk/types"
)

// --- MockLogger ---

type MockLogger struct {
	name string
}

func NewMockLogger(name string) *MockLogger {
	return &MockLogger{name: name}
}

func (m *MockLogger) GetName() string { return m.name }
func (m *MockLogger) CreateSubLogger(name string) types.ILogger {
	return NewMockLogger(m.name + "." + name)
}
func (m *MockLogger) Scope(name string, handler func(log types.ILogger) error) error {
	return handler(m.CreateSubLogger(name))
}
func (m *MockLogger) ScopeWithReturn(name string, handler func(log types.ILogger) (any, error)) (any, error) {
	return handler(m.CreateSubLogger(name))
}
func (m *MockLogger) GetBaseLogger() *logrus.Entry              { return logrus.NewEntry(logrus.New()) }
func (m *MockLogger) Trace(message string)                      {}
func (m *MockLogger) TraceF(format string, args ...interface{}) {}
func (m *MockLogger) Debug(message string)                      {}
func (m *MockLogger) DebugF(format string, args ...interface{}) {}
func (m *MockLogger) Info(message string)                       {}
func (m *MockLogger) InfoF(format string, args ...interface{})  {}
func (m *MockLogger) Warn(message string)                       {}
func (m *MockLogger) WarnF(format string, args ...interface{})  {}
func (m *MockLogger) Error(message string)                      {}
func (m *MockLogger) ErrorF(format string, args ...interface{}) {}
func (m *MockLogger) Panic(message string)                      {}
func (m *MockLogger) PanicF(format string, args ...interface{}) {}

// --- MockKernel ---

type MockKernel struct {
	name        string
	version     string
	runDir      string
	unsecure    bool
	log         types.ILogger
	app         types.IApp
	rootContext context.Context
}

func NewMockKernel(name, version string) *MockKernel {
	return &MockKernel{
		name:        name,
		version:     version,
		runDir:      "/tmp/k2-test",
		unsecure:    false,
		log:         NewMockLogger("kernel"),
		rootContext: context.Background(),
	}
}

func (m *MockKernel) GetLogger() types.ILogger        { return m.log }
func (m *MockKernel) IsUnsecure() bool                { return m.unsecure }
func (m *MockKernel) GetApp() types.IApp              { return m.app }
func (m *MockKernel) GetRunDirectory() string         { return m.runDir }
func (m *MockKernel) GetRootContext() context.Context { return m.rootContext }
func (m *MockKernel) GetService(key types.KernelServiceContextKey) types.IKernelService {
	svc := m.rootContext.Value(key)
	if svc == nil {
		return nil
	}
	return svc.(types.IKernelService)
}
func (m *MockKernel) SetService(service types.IKernelService) {
	key := types.KernelServiceContextKey(service.GetName())
	m.rootContext = context.WithValue(m.rootContext, key, service)
}
func (m *MockKernel) SetApp(app types.IApp) { m.app = app }
func (m *MockKernel) SetRunDir(dir string)  { m.runDir = dir }
func (m *MockKernel) Init() error           { return nil }
func (m *MockKernel) Register() error       { return nil }
func (m *MockKernel) Start() error          { return nil }
func (m *MockKernel) ListenAndServe() error { return nil }

// --- MockApp ---

type MockApp struct {
	name       string
	version    string
	log        types.ILogger
	kernel     types.IKernel
	components []types.IAppComponent
}

func NewMockApp(name, version string) *MockApp {
	return &MockApp{
		name:    name,
		version: version,
		log:     NewMockLogger("app"),
	}
}

func (m *MockApp) GetLogger() types.ILogger                  { return m.log }
func (m *MockApp) SetLogger(logger types.ILogger)            { m.log = logger }
func (m *MockApp) GetName() string                           { return m.name }
func (m *MockApp) GetVersion() string                        { return m.version }
func (m *MockApp) GetDocs() *swag.Spec                       { return nil }
func (m *MockApp) GetUI() *embed.FS                          { return nil }
func (m *MockApp) GetConfig() *embed.FS                      { return nil }
func (m *MockApp) GetExternals() *embed.FS                   { return nil }
func (m *MockApp) GetComponents() []types.IAppComponent      { return m.components }
func (m *MockApp) SetComponents(comps []types.IAppComponent) { m.components = comps }
func (m *MockApp) GetComponent(name string) types.IAppComponent {
	for _, c := range m.components {
		if c.GetName() == name {
			return c
		}
	}
	return nil
}
func (m *MockApp) AddComponent(ctor types.AppComponentCtor) {}
func (m *MockApp) GetKernel() types.IKernel                 { return m.kernel }
func (m *MockApp) SetKernel(kernel types.IKernel)           { m.kernel = kernel }

// --- MockComponent ---

type MockComponent struct {
	name         string
	order        int
	app          types.IApp
	log          types.ILogger
	accessPolicy types.IAccessPolicy
	controllers  []types.IAppController
}

func NewMockComponent(app types.IApp, name string, order int) *MockComponent {
	var log types.ILogger
	if app.GetLogger() != nil {
		log = app.GetLogger().CreateSubLogger(name)
	} else {
		log = NewMockLogger(name)
	}
	return &MockComponent{
		name:         name,
		order:        order,
		app:          app,
		log:          log,
		accessPolicy: types.AccessPolicyPublic,
	}
}

func (m *MockComponent) GetName() string                             { return m.name }
func (m *MockComponent) GetOrder() int                               { return m.order }
func (m *MockComponent) GetApp() types.IApp                          { return m.app }
func (m *MockComponent) GetLogger() types.ILogger                    { return m.log }
func (m *MockComponent) GetConfig() *embed.FS                        { return nil }
func (m *MockComponent) GetDocs() *swag.Spec                         { return nil }
func (m *MockComponent) GetUI() *embed.FS                            { return nil }
func (m *MockComponent) GetAccessPolicy() types.IAccessPolicy        { return m.accessPolicy }
func (m *MockComponent) GetControllers() []types.IAppController      { return m.controllers }
func (m *MockComponent) SetControllers(ctrls []types.IAppController) { m.controllers = ctrls }
func (m *MockComponent) GetController(name string) types.IAppController {
	for _, c := range m.controllers {
		if c.GetName() == name {
			return c
		}
	}
	return nil
}
func (m *MockComponent) Init() error                       { return nil }
func (m *MockComponent) Register(r *gin.RouterGroup) error { return nil }
func (m *MockComponent) Start() error                      { return nil }
func (m *MockComponent) Stop() error                       { return nil }

// --- MockController ---

type MockController struct {
	name         string
	order        int
	component    types.IAppComponent
	log          types.ILogger
	accessPolicy types.IAccessPolicy
	schedules    []types.IAppSchedule
}

func NewMockController(comp types.IAppComponent, name string, order int) *MockController {
	var log types.ILogger
	if comp.GetLogger() != nil {
		log = comp.GetLogger().CreateSubLogger(name)
	} else {
		log = NewMockLogger(name)
	}
	return &MockController{
		name:         name,
		order:        order,
		component:    comp,
		log:          log,
		accessPolicy: types.AccessPolicyPublic,
		schedules:    make([]types.IAppSchedule, 0),
	}
}

func (m *MockController) GetName() string                      { return m.name }
func (m *MockController) GetOrder() int                        { return m.order }
func (m *MockController) GetComponent() types.IAppComponent    { return m.component }
func (m *MockController) GetLogger() types.ILogger             { return m.log }
func (m *MockController) GetConfig() *embed.FS                 { return nil }
func (m *MockController) GetAccessPolicy() types.IAccessPolicy { return m.accessPolicy }
func (m *MockController) GetSchedules() []types.IAppSchedule   { return m.schedules }
func (m *MockController) Init() error                          { return nil }
func (m *MockController) Register(r *gin.RouterGroup) error    { return nil }
func (m *MockController) Start() error                         { return nil }
func (m *MockController) Stop() error                          { return nil }

func (m *MockController) AddTestSchedule(name, cron string, handler types.AppScheduleHandler) {
	m.schedules = append(m.schedules, &mockSchedule{name: name, cron: cron, handler: handler})
}

// --- mockSchedule ---

type mockSchedule struct {
	name    string
	cron    string
	handler types.AppScheduleHandler
}

func (s *mockSchedule) GetName() string                          { return s.name }
func (s *mockSchedule) GetCron() string                          { return s.cron }
func (s *mockSchedule) GetTaskHandler() types.AppScheduleHandler { return s.handler }

// --- MockKernelService ---

type MockKernelService struct {
	name   string
	kernel types.IKernel
	log    types.ILogger
	config map[string]string
	data   map[string]any
}

func NewMockKernelService(kernel types.IKernel, name string) *MockKernelService {
	return &MockKernelService{
		name:   name,
		kernel: kernel,
		log:    NewMockLogger(name),
		config: make(map[string]string),
		data:   make(map[string]any),
	}
}

func (m *MockKernelService) GetName() string                    { return m.name }
func (m *MockKernelService) GetKernel() types.IKernel           { return m.kernel }
func (m *MockKernelService) GetLogger() types.ILogger           { return m.log }
func (m *MockKernelService) GetConfig(key string) string        { return m.config[key] }
func (m *MockKernelService) SetConfig(key string, value string) { m.config[key] = value }
func (m *MockKernelService) GetData(key string) any             { return m.data[key] }
func (m *MockKernelService) SetData(key string, value any)      { m.data[key] = value }
func (m *MockKernelService) Init() error                        { return nil }
func (m *MockKernelService) Register() error                    { return nil }
func (m *MockKernelService) Start() error                       { return nil }
func (m *MockKernelService) Stop() error                        { return nil }

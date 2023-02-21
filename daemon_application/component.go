package daemon_application

type ComponentID string
type ComponentStatus uint16

const (
	ComponentBaseInitStatus = iota
	ComponentInitStatus
	ComponentStartStatus
	ComponentStopStatus
)

type IComponent interface {
	Initialize(args ...interface{}) error
	GetID() ComponentID
	Start() error
	Stop() error
	setAppProxy(appProxy IApplication)
	GetAppProxy() IApplication
}

type BaseComponent struct {
	id       ComponentID
	status   ComponentStatus
	appProxy IApplication
}

func (t *BaseComponent) Initialize(args ...interface{}) error {
	return nil
}

func (t *BaseComponent) GetID() ComponentID {
	return t.id
}

func (t *BaseComponent) GetStatus() ComponentStatus {
	return t.status
}

func (t *BaseComponent) setPartStatus(status ComponentStatus) {
	t.status = status
}

func (t *BaseComponent) Start() error {
	return nil
}

func (t *BaseComponent) Stop() error {
	return nil
}

func (t *BaseComponent) setAppProxy(appProxy IApplication) {
	t.appProxy = appProxy
}

func (t *BaseComponent) GetAppProxy() IApplication {
	return t.appProxy
}

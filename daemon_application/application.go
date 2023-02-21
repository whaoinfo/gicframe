package daemon_application

import (
	"fmt"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/go-box/ossignal"
	"os"
	"syscall"
)

const (
	APPDefaultSignChanSize = 1
)

type ApplicationID string

type IApplication interface {
	baseInitialize(appID ApplicationID) error
	Initialize(args ...interface{}) error
	GetID() ApplicationID
	importComponent(componentID ComponentID, newFunc NewComponentFunc, args []interface{}) error
	start() error
	OnStart() error
	stop() error
	OnStop() error
	forever()
}

type BaseApplication struct {
	id            ApplicationID
	signalHandler *ossignal.SignalHandler
	componentMap  map[ComponentID]IComponent
}

func (t *BaseApplication) baseInitialize(appID ApplicationID) error {
	t.id = appID
	t.componentMap = make(map[ComponentID]IComponent)
	t.signalHandler = &ossignal.SignalHandler{}
	if err := t.signalHandler.InitSignalHandler(APPDefaultSignChanSize); err != nil {
		return err
	}
	for _, sig := range []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT} {
		t.signalHandler.RegisterSignal(sig, func() {
			t.signalHandler.CloseSignalHandler()
		})
	}

	return nil
}

func (t *BaseApplication) Initialize(args ...interface{}) error {
	return nil
}

func (t *BaseApplication) GetID() ApplicationID {
	return t.id
}

func (t *BaseApplication) start() error {
	for componentID, component := range t.componentMap {
		logger.InfoFmt("Start %v component...", componentID)
		if err := component.Start(); err != nil {
			return fmt.Errorf("failed to start %v component, err: %v", componentID, err)
		}
		logger.InfoFmt("Started %v component", componentID)
	}
	return nil
}

func (t *BaseApplication) OnStart() error {
	return nil
}

func (t *BaseApplication) stop() error {
	for componentID, component := range t.componentMap {
		logger.InfoFmt("Stop %v component...", componentID)
		if err := component.Stop(); err != nil {
			return fmt.Errorf("failed to stop component, component ID: %v, err: %v", componentID, err)
		}
		logger.InfoFmt("Stopped %v component", componentID)
	}
	return nil
}

func (t *BaseApplication) OnStop() error {
	return nil
}

func (t *BaseApplication) forever() {
	t.signalHandler.ListenSignal()
}

func (t *BaseApplication) importComponent(componentID ComponentID, newFunc NewComponentFunc, args []interface{}) error {
	component := newFunc()
	if err := component.Initialize(args...); err != nil {
		return err
	}

	component.setAppProxy(t)
	t.componentMap[componentID] = component
	return nil
}

package daemon_application

import (
	"github.com/whaoinfo/go-box/logger"
)

type NewApplicationFunc func() (IApplication, ApplicationID, map[ComponentID]NewComponentFunc)

type NewComponentFunc func() IComponent

// LaunchDaemonApplication
// 创建运行应用程序
func LaunchDaemonApplication(workPath string, newAPPFunc NewApplicationFunc,
	appArgs []interface{}, conf IConfig, enabledDevMode bool, disableDaemon bool) error {

	// 创建应用
	app, appID, newComponents := newAPPFunc()
	logger.InfoFmt("Create %v application, workPath: %v, enabledDevMode: %v",
		appID, workPath, enabledDevMode)

	// 初始化应用
	if err := app.baseInitialize(appID); err != nil {
		return err
	}
	if err := app.Initialize(appArgs...); err != nil {
		return err
	}
	logger.Info("Initialize succeeded")

	// 初始化配置
	if err := configInitialize(conf, workPath, appID, enabledDevMode); err != nil {
		return err
	}
	logger.Info("Config initialize succeeded")

	// 导入应用部件
	for componentID, f := range newComponents {
		if err := app.importComponent(componentID, f, appArgs); err != nil {
			return err
		}
		logger.InfoFmt("Imported %v component", componentID)
	}

	// 启动应用
	logger.Info("Start application...")
	if err := app.start(); err != nil {
		return err
	}
	if err := app.OnStart(); err != nil {
		return err
	}
	logger.Info("Started application")

	if !disableDaemon {
		app.forever()
	}

	// 停止应用
	logger.Info("Stop application...")
	if err := app.stop(); err != nil {
		return err
	}
	if err := app.OnStop(); err != nil {
		return err
	}
	logger.Info("Stopped application")

	return nil
}

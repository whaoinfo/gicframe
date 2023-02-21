package daemon_application

import (
	"encoding/json"
	"fmt"
	"github.com/whaoinfo/go-box/logger"
	"io/ioutil"
	"path"
)

type IConfig interface {
	setDirPath(string)
	GetDirPath() string
	OnParse() error
	SetData(data []byte)
	GetData() []byte
}

type BaseConfig struct {
	dirPath string
	data    []byte
}

func (t *BaseConfig) setDirPath(dirPath string) {
	t.dirPath = dirPath
}

func (t *BaseConfig) GetDirPath() string {
	return t.dirPath
}

func (t *BaseConfig) OnParse() error {
	return nil
}

func (t *BaseConfig) SetData(data []byte) {
	dataLen := len(data)
	if dataLen <= 0 {
		return
	}

	t.data = make([]byte, dataLen)
	copy(t.data, data)
}

func (t *BaseConfig) GetData() []byte {
	return t.data
}

func configInitialize(conf IConfig, workPath string, appID ApplicationID, enabledDevMode bool) error {
	if conf == nil {
		return nil
	}

	confData := conf.GetData()
	if len(confData) > 0 {
		if err := json.Unmarshal(confData, conf); err != nil {
			return err
		}
		return conf.OnParse()
	}

	dirPath := ""
	if enabledDevMode {
		dirPath = path.Join(workPath, "config_template", fmt.Sprintf("%v", appID))
	} else {
		dirPath = path.Join(workPath, "config")
	}
	confPath := path.Join(dirPath, "config.json")

	if err := loadConfig(confPath, conf); err != nil {
		return err
	}

	conf.setDirPath(dirPath)
	return conf.OnParse()
}

func loadConfig(configPath string, conf IConfig) error {
	if configPath != "" && configPath != "unused" {
		logger.InfoFmt("Config file path: %v", configPath)
		fileData, readFileErr := ioutil.ReadFile(configPath)
		if readFileErr != nil {
			return readFileErr
		}

		if err := json.Unmarshal(fileData, conf); err != nil {
			return err
		}
	}

	return nil
}

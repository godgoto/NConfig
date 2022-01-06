package NConfig

import (
	"errors"
	"fmt"
	"github.com/Unknwon/goconfig"
	"path/filepath"
	"reflect"
	"time"
)

type LocalIniConfig struct {
	nc                *LocalIniConfig
	config            interface{}
	path              string
	env               string
	callback          CallBackHandler
	errCallback       errorCallBackHandler
	updateMillisecond int64
	bCallback         bool
	listStruct        []*BaseType
}

func (p *LocalIniConfig) save(nc *LocalIniConfig) *LocalIniConfig {
	p.nc = nc
	p.bCallback = false
	return p.nc
}

func (p *LocalIniConfig) Source(path string) *LocalIniConfig {
	p.path = path
	return p.nc
}

func (p *LocalIniConfig) Env(env string) *LocalIniConfig {
	p.env = env
	p.envPath()
	return p.nc
}
func (p *LocalIniConfig) envPath() {
	if len(p.path) > 0 && len(p.env) > 0 {
		dir, filename := filepath.Split(p.path)
		suffix := filepath.Ext(p.path)
		iLen := len(filename) - len(suffix)
		name := filename[0:iLen]
		tempPath := fmt.Sprintf("%v%v-%v%v", dir, name, p.env, suffix)
		p.path = tempPath
	}
}

func (p *LocalIniConfig) UpdateTime(millisecond int64) *LocalIniConfig {
	p.nc.updateMillisecond = millisecond
	return p.nc
}

func (p *LocalIniConfig) CallBack(call CallBackHandler) *LocalIniConfig {
	p.callback = call
	return p.nc
}

func (p *LocalIniConfig) Error(call errorCallBackHandler) *LocalIniConfig {
	p.errCallback = call
	return p.nc
}

func (p *LocalIniConfig) Sync(cfg interface{}) (*LocalIniConfig, error) {
	if len(p.path) == 0 {
		return p.nc, errors.New("path cannot be empty!")
	}

	if !FileExist(p.path) {
		return p.nc, errors.New(p.path + "file does not exist!")
	}

	if err := p.analysis(); err != nil {
		return p.nc, err
	}
	p.read(cfg)
	p.loop()

	return p.nc, nil
}

func (p *LocalIniConfig) analysis() error {
	list, err := analysisStruct(p.config, "nk")
	p.listStruct = list
	return err
}

func (p *LocalIniConfig) loop() {
	go func() {
		p.bCallback = true
		for p.bCallback {
			if p.updateMillisecond > 0 && p.callback != nil {
				loopCfg := reflect.New(reflect.TypeOf(p.config)).Interface()
				loopErr := p.read(loopCfg)
				if loopErr != nil && p.errCallback != nil {
					p.errCallback(loopErr)
				}
				p.callback(loopCfg)
			}
			time.Sleep(time.Duration(p.updateMillisecond) * time.Millisecond)
		}
	}()
}

func (p *LocalIniConfig) Stop() {
	p.bCallback = false
}

func (p *LocalIniConfig) Save() error {

	return nil
}

func (p *LocalIniConfig) read(cfg interface{}) error {
	section(p.listStruct)
	loalCfg, err := goconfig.LoadConfigFile(p.path)
	if err != nil {
		return err
	}
	for _, v := range p.listStruct {
		v.Value, err = loalCfg.GetValue(v.Section, v.Key)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		saveLocalValue(cfg, v)
	}
	return nil
}

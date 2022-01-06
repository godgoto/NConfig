package NConfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

type LocalJsonConfig struct {
	nc                *LocalJsonConfig
	config            interface{}
	path              string
	env               string
	callback          CallBackHandler
	errCallback       errorCallBackHandler
	updateMillisecond int64
	bCallback         bool
	section           string
	listStruct        []*BaseType
}

func (p *LocalJsonConfig) save(nc *LocalJsonConfig) *LocalJsonConfig {
	p.nc = nc
	p.bCallback = false
	return p.nc
}

func (p *LocalJsonConfig) Source(path string) *LocalJsonConfig {
	p.path = path
	p.envPath()
	return p.nc
}

func (p *LocalJsonConfig) Env(env string) *LocalJsonConfig {
	p.env = env
	p.envPath()
	return p.nc
}
func (p *LocalJsonConfig) envPath() {
	if len(p.path) > 0 && len(p.env) > 0 {
		dir, filename := filepath.Split(p.path)
		suffix := filepath.Ext(p.path)
		iLen := len(filename) - len(suffix)
		name := filename[0:iLen]
		tempPath := fmt.Sprintf("%v%v-%v%v", dir, name, p.env, suffix)
		p.path = tempPath
	}
}

func (p *LocalJsonConfig) Section(section string) *LocalJsonConfig {
	p.section = section
	return p.nc
}

func (p *LocalJsonConfig) UpdateTime(millisecond int64) *LocalJsonConfig {
	p.nc.updateMillisecond = millisecond
	return p.nc
}

func (p *LocalJsonConfig) CallBack(call CallBackHandler) *LocalJsonConfig {
	p.callback = call
	return p.nc
}

func (p *LocalJsonConfig) Error(call errorCallBackHandler) *LocalJsonConfig {
	p.errCallback = call
	return p.nc
}



func (p *LocalJsonConfig) Sync(cfg interface{}) (*LocalJsonConfig, error) {
	if len(p.path) == 0 {
		return p.nc, errors.New("path cannot be empty!")
	}

	if !FileExist(p.path) {
		return p.nc, errors.New(p.path + " file does not exist!")
	}

	if err := p.analysis(); err != nil {
		return p.nc, err
	}
	
	p.read(cfg)
	p.loop()
	return p.nc, nil
}

func (p *LocalJsonConfig) analysis() error {
	list, err := analysisStruct(p.config, "nk")
	p.listStruct = list
	return err
}

func (p *LocalJsonConfig) loop() {
	go func() {
		p.bCallback = true
		for p.bCallback {
			if p.updateMillisecond > 0 && p.callback != nil {
				loopCfg := reflect.New(reflect.TypeOf(p.config)).Interface()
				p.read(loopCfg)
				p.callback(loopCfg)
			}
			time.Sleep(time.Duration(p.updateMillisecond) * time.Millisecond)
		}
	}()
}

func (p *LocalJsonConfig) read(cfg interface{}) {
	jsonClient, err := NewJSonConfig(p.path)
	for _, v := range p.listStruct {
		if len(p.section) > 0 {
			v.Section = p.section
		}
		v.Value, err = jsonClient.GetValue(v.Section, v.Key)
		if err != nil {
			if p.errCallback != nil {
				p.errCallback(err)
			}
			continue
		}
		saveLocalValue(cfg, v)
	}
}
func (p *LocalJsonConfig) Stop() {
	p.bCallback = false
}

type JSonConfig struct {
	path    string
	file    *os.File
	content *[]byte
}

func NewJSonConfig(path string) (*JSonConfig, error) {
	p := &JSonConfig{}
	file, _ := os.Open(path)
	p.file = file
	defer p.file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return p, err
	}
	p.content = &data
	return p, nil
}

func (p *JSonConfig) GetValue(section, key string) (interface{}, error) {
	if len(key) == 0 {
		return "", nil
	}
	var confighash = make(map[string]interface{})
	err := json.Unmarshal(*p.content, &confighash)

	if err == nil {
		if len(section) > 0 {
			if _, ok := confighash[section]; ok {
				objecthash := confighash[section].(map[string]interface{})
				if _, ok := objecthash[key]; ok {
					values := objecthash[key]
					return values, nil
				} else {
					return "", errors.New("There is no key : " + key)
				}
			}
		} else {
			if _, ok := confighash[key]; ok {
				val := confighash[key]
				return val, err
			}
		}
		return "", errors.New(" There is no section : " + section)
	}
	errorStr := fmt.Sprintf("data is not in JSON format! section:%v key:%v", section, key)
	return "", errors.New(errorStr)
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

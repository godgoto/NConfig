package NConfig

import (
	"errors"
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"reflect"
)

type ApolloConfig struct {
	nc                *ApolloConfig
	config            interface{}
	callback          CallBackHandler
	errCallback       errorCallBackHandler
	updateMillisecond int64
	bCallback         bool
	c                 *config.AppConfig
	client            agollo.Client
	listStruct        []*BaseType
}

func (p *ApolloConfig) save(nc *ApolloConfig) *ApolloConfig {
	p.nc = nc
	p.bCallback = false
	return p.nc
}

func (p *ApolloConfig) UpdateTime(millisecond int64) *ApolloConfig {
	p.nc.updateMillisecond = millisecond
	return p.nc
}

func (p *ApolloConfig) CallBack(call CallBackHandler) *ApolloConfig {
	p.callback = call
	return p.nc
}

func (p *ApolloConfig) Error(call errorCallBackHandler) *ApolloConfig {
	p.errCallback = call
	return p.nc
}

func (p *ApolloConfig) Sync(cfg interface{}) (*ApolloConfig, error) {
	if err := p.analysis(); err != nil {
		return p.nc, err
	}
	p.newApolloClient()
	p.read(cfg)
	p.loop()
	return p.nc, nil
}

func (p *ApolloConfig) analysis() error {
	list, err := analysisStruct(p.config, "nk")
	p.listStruct = list
	return err
}

func (p *ApolloConfig) newApolloClient() {
	agollo.SetLogger(&DefaultLogger{p.errCallback})
	p.client, _ = agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return p.c, nil
	})
}
func (p *ApolloConfig) read(cfg interface{}) {
	cache := p.client.GetConfigCache(p.c.NamespaceName)
	for _, v := range p.listStruct {
		val, err := cache.Get(v.Key)
		if err != nil && p.errCallback != nil {
			p.errCallback(err)
			continue
		}
		var str = fmt.Sprintf("%v", val)
		v.Value = str
		saveApolloValue(cfg, v)
	}
}

func (p *ApolloConfig) loop() {

	p.client.AddChangeListener(&ChangeListener{callback: p.callback, config: p.config, list: p.listStruct})

}

func (p *ApolloConfig) Source(cfg interface{}) *ApolloConfig {
	var appConfig config.AppConfig
	list, _ := analysisStruct(cfg, "ac")
	for _, v := range list {
		if "AppID" == v.Key {
			appConfig.AppID = v.Value.(string)
		}
		if "Cluster" == v.Key {
			appConfig.Cluster = v.Value.(string)
		}
		if "IP" == v.Key {
			appConfig.IP = v.Value.(string)
		}
		if "NamespaceName" == v.Key {
			appConfig.NamespaceName = v.Value.(string)
		}
		if "IsBackupConfig" == v.Key {
			appConfig.IsBackupConfig = true //v.Value
		}
		if "Secret" == v.Key {
			appConfig.Secret = v.Value.(string)
		}
	}
	p.c = &appConfig
	return p.nc
}

type ChangeListener struct {
	callback CallBackHandler
	config   interface{}
	list     []*BaseType
}

//OnChange 增加变更监控
func (c *ChangeListener) OnChange(event *storage.ChangeEvent) {
}

//OnNewestChange 监控最新变更
func (c *ChangeListener) OnNewestChange(event *storage.FullChangeEvent) () {
	cfg := reflect.New(reflect.TypeOf(c.config)).Interface()
	for _, v := range c.list {
		if _, ok := event.Changes[v.Key]; ok {
			v.Value = event.Changes[v.Key].(string)
			saveApolloValue(cfg, v)
		}
	}
	if c.callback != nil {
		c.callback(cfg)
	}
}

//DefaultLogger 默认日志实现
type DefaultLogger struct {
	errCallback errorCallBackHandler
}

//Debugf debug 格式化
func (d *DefaultLogger) Debugf(format string, params ...interface{}) {

}

//Infof 打印info
func (d *DefaultLogger) Infof(format string, params ...interface{}) {

}

//Warnf warn格式化
func (d *DefaultLogger) Warnf(format string, params ...interface{}) {
	if d.errCallback != nil {
		str := fmt.Sprintf(format, params)
		d.errCallback(errors.New(str))
	}
}

//Errorf error格式化
func (d *DefaultLogger) Errorf(format string, params ...interface{}) {
	str := fmt.Sprintf(format, params)
	if d.errCallback != nil {
		d.errCallback(errors.New(str))
	}
}

//Debug 打印debug
func (d *DefaultLogger) Debug(v ...interface{}) {

}

//Info 打印Info
func (d *DefaultLogger) Info(v ...interface{}) {

}

//Warn 打印Warn
func (d *DefaultLogger) Warn(v ...interface{}) {
}

//Error 打印Error
func (d *DefaultLogger) Error(v ...interface{}) {
}

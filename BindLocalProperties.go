package NConfig

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
)

type PropertiesConfig struct {
	nc                *PropertiesConfig
	config            interface{}
	path              string
	env               string
	callback          CallBackHandler
	errCallback       errorCallBackHandler
	updateMillisecond int64
	bCallback         bool
	listStruct        []*BaseType
}

func (p *PropertiesConfig) save(nc *PropertiesConfig) *PropertiesConfig {
	p.nc = nc
	p.bCallback = false
	return p.nc
}

func (p *PropertiesConfig) Source(path string) *PropertiesConfig {
	p.path = path
	return p.nc
}

func (p *PropertiesConfig) Env(env string) *PropertiesConfig {
	p.env = env
	p.envPath()
	return p.nc
}

func (p *PropertiesConfig) Stop() {
	p.bCallback = false
}

func (p *PropertiesConfig) envPath() {
	if len(p.path) > 0 && len(p.env) > 0 {
		dir, filename := filepath.Split(p.path)
		suffix := filepath.Ext(p.path)
		iLen := len(filename) - len(suffix)
		name := filename[0:iLen]
		tempPath := fmt.Sprintf("%v%v-%v%v", dir, name, p.env, suffix)
		p.path = tempPath
	}
}

func (p *PropertiesConfig) UpdateTime(millisecond int64) *PropertiesConfig {
	p.nc.updateMillisecond = millisecond
	return p.nc
}

func (p *PropertiesConfig) CallBack(call CallBackHandler) *PropertiesConfig {
	p.callback = call
	return p.nc
}

func (p *PropertiesConfig) Error(call errorCallBackHandler) *PropertiesConfig {
	p.errCallback = call
	return p.nc
}

func (p *PropertiesConfig) Sync(cfg interface{}) (*PropertiesConfig, error) {

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

func (p *PropertiesConfig) analysis() error {
	list, err := analysisStruct(p.config, "nk")
	p.listStruct = list
	return err
}

func (p *PropertiesConfig) loop() {
	go func() {
		for true {
			if p.bCallback {
				loopCfg := reflect.New(reflect.TypeOf(p.config)).Interface()
				if p.updateMillisecond > 0 && p.callback != nil {
					p.read(loopCfg)
					p.callback(loopCfg)
				}
				time.Sleep(time.Duration(p.updateMillisecond) * time.Millisecond)
			}
		}
	}()
}

func (p *PropertiesConfig) read(cfg interface{}) error {
	pro := NewProperties()
	err := pro.LoadFromFile(p.path)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, v := range p.listStruct {
		if len(v.Key) == 0 {
			continue
		}
		val, isExist := pro.Property(v.Key)
		v.Value = val
		if !isExist {
			errStr := v.Key + " not exist"
			if p.errCallback != nil {
				p.errCallback(errors.New(errStr))
			}
		}
		saveApolloValue(cfg, v)
	}
	return nil
}


// Create an empty property list.
func NewProperties() *Properties {
	return &Properties{
		object: make(map[string]string),
	}
}

type Properties struct {
	m      sync.Mutex
	object map[string]string
}

// 用指定的键在此属性列表中搜索属性。
//
// Searches for the property with the specified key in this property list.
func (p *Properties) Property(key string) (value string, isExist bool) {
	if value, ok := p.object[key]; !ok {
		return "", false
	} else {
		return value, true
	}
}

// 用指定的键在此属性列表中搜索属性，把","连接的多个属性转换为切片返回。
//
// Search for attributes in this attribute list using the specified key to return multiple attributes connected by "," converted to slices.
func (p *Properties) PropertySlice(key string) (values []string, isExist bool) {
	strs, isExist := p.Property(key)
	if !isExist {
		return nil, false
	}
	return strings.Split(strs, ","), true
}

// 为指定的键设置多个属性，把多个属性值转换成“，”连接的属性字符串。
// Set  multiple attributes for the specified key, converting multiple attribute values into a "," concatenated attribute string.
func (p *Properties) SetPropertySlice(key string, values ...string) {
	value := strings.Join(values, ",")
	p.SetProperty(key, value)
}

// 更新指定的键和属性,如果键不存在就新建。
//
// Update the specified key and properties. If the key does not exist, create a new one.
func (p *Properties) SetProperty(key, value string) {
	p.m.Lock()
	p.object[key] = value
	p.m.Unlock()
}

// 从输入流读取属性列表。
//
// Reads a property list (key and element pairs) from the input character stream in a simple line-oriented format.
func (p *Properties) Load(r io.Reader) error {
	br := bufio.NewReader(r)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		p.parseLine(line)
	}
	return nil
}

// 从文件中读取属性列表。
//
// Reads a property list from a file
func (p *Properties) LoadFromFile(filePath string) error {
	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("os.Open failed: %v", err)
	}
	return p.Load(f)
}

// 将属性列表写入输出流。
//
// Writes this property list (key and element pairs) in this Properties table to the output stream in a format suitable for loading into a Properties table using the Load() method.
func (p *Properties) Store(w io.Writer) error {
	var buf bytes.Buffer
	for k, v := range p.object {
		p.line(k, v, &buf)
	}
	_, err := w.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("w.Write failed: %v", err)
	}
	return nil
}

// 将属性列表写入文件。
//
// Writes a list of property to a file.
func (p *Properties) StoreToFile(filePath string) error {
	f, err := os.Create(filePath)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("os.Open failed: %v", err)
	}
	return p.Store(f)

}

func (p *Properties) line(key, value string, buf *bytes.Buffer) {
	buf.WriteString(key)
	buf.WriteString(" = ")
	buf.WriteString(value)
	buf.WriteByte('\n')
}

func (p *Properties) parseLine(line []byte) {
	lineStr := strings.TrimSpace(string(line))
	if isCommentline(lineStr) {
		return
	}
	splitStrs := strings.Split(lineStr, "=")
	key := strings.TrimSpace(splitStrs[0])
	value := strings.TrimSpace(splitStrs[1])
	p.SetProperty(key, value)
}

// 返回属性列表中所有键的枚举。
//
// Returns an enumeration of all keys in the property list.
func (p *Properties) PropertyNames() []string {
	var names []string
	for k := range p.object {
		names = append(names, k)
	}
	return names
}

// 过滤注释的行
func isCommentline(line string) bool {
	if len(line) == 0 {
		return true
	}
	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
		return true
	}
	return false
}

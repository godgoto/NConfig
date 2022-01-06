### nconfg 一个配置工具类

#### 1. 功能

|  格式   | 是否支持  |
|  ----  | ----  |
| ini  | 支持 |
| json  | 支持 |
| properties  | 支持 |
| apollo  | 支持 |

|  开发中   |
|  ----  | 
| 1. redis 配置的支持  | 
| 2. struct 嵌套的支持  |
| 3. Default 的支持  | 
| 4. 回写功能 |

#### 2. struct 定义

```cgo
type LocalConfig struct {
	AppID            string  `nk:"appId" ac:"AppID"`
	Cluster          string  `nk:"cluster" ac:"Cluster"`
	NamespaceName    string  `nk:"namespaceName" ac:"NamespaceName"`
	IP               string  `nk:"ip" ac:"IP"`
	NextTryConnTime  int64   ``
	IsBackupConfig   bool    `nk:"isBackupConfig" ac:"isBackupConfig" default:"true"`
	BackupConfigPath string  `nk:"backupConfigPath"`
	Secret           string  `nk:"secret"`
	RetryCount       int64   `nk:"retryCount"`
	Pi               float64 `nk:"pi"`
}
    
```

    nk : 配置文件的key
    ac : apoollo 对的配置的key 用户自apollo 自动查找配置
    NConfig.BindApolloModels(MySqlConfig{}).Source(localCfg)   Source 可以直接根据 ac的 key 读取apollo的配置

#### 3.1 读取Ini

```cgo

    var qc = &IniConfig{}
	path := "./test.ini"
	callBackHandler := func(config interface{}) {
		qc = config.(*IniConfig)
		str, _ := json.Marshal(qc)
		fmt.Println("回调 : ", string(str))
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	client, err := NConfig.BindLocalIniModels(IniConfig{}).Source(path).CallBack(callBackHandler).UpdateTime(3000).Error(errHandler).Sync(qc)
```

#### 3.2读取 json

```cgo

    var jc = &LocalConfig{}
	path := "./app.json"
	callBackHandler := func(config interface{}) {
		jc = config.(*LocalConfig)
		str, _ := json.Marshal(jc)
		fmt.Println("CallBack : ", string(str))
	}
	errHandler := func(err error) {
		fmt.Println("error:", err)
	}

	client, err := NConfig.BindLocalJsonModels(LocalConfig{}).Section("dev").Source(path).CallBack(callBackHandler).UpdateTime(3000).Error(errHandler).Sync(jc)
	if err != nil {
		fmt.Println("error:", err.Error())
	}

	str, _ := json.Marshal(jc)
	fmt.Println(string(str))
```

#### 3.3读取 Properties

```cgo

    localCfg := &LocalConfig{}
	path := "./app.properties"

	callBackHandler := func(config interface{}) {
		localCfg =  config.(*LocalConfig)
		str, _ := json.Marshal(localCfg)
		fmt.Println("回调  : ", string(str))
	}

	_, err := NConfig.BindPropertiesModels(LocalConfig{}).Source(path).Env("dev").CallBack(callBackHandler).UpdateTime(3000).Sync(localCfg)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	str, _ := json.Marshal(localCfg)
	fmt.Println(string(str))
```

#### 3.4 读取 Apollo

```cgo

    var mysqlCfg MySqlConfig
	callBackHandler := func(config interface{}) {
		cfg := * config.(*MySqlConfig)
		str, _ := json.Marshal(cfg)
		fmt.Println("回调  : ", string(str))
	}

	errHandler := func(err error) {
		fmt.Println("error:", err)
	}

	NConfig.BindApolloModels(MySqlConfig{}).Source(localCfg).CallBack(callBackHandler).Error(errHandler).Sync(&mysqlCfg)

	str2, _ := json.Marshal(mysqlCfg)
	fmt.Println("获取 : ", string(str2))

	for true {
		time.Sleep(time.Second * 3)
	}

```

#### 4 引用 & 鸣谢

    github.com/apolloconfig/agollo
    github.com/obity/properties
    github.com/Unknwon/goconfig
    


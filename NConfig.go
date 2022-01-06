package NConfig

type NConfig struct {
}

type CallBackHandler func(config interface{})
type errorCallBackHandler func(err error)

func BindLocalIniModels(config interface{}) *LocalIniConfig {
	p := &LocalIniConfig{config: config, callback: nil, errCallback: nil, updateMillisecond: 5000}
	p.save(p)
	return p
}

func BindLocalJsonModels(config interface{}) *LocalJsonConfig {
	p := &LocalJsonConfig{nc: &LocalJsonConfig{}, config: config, callback: nil, errCallback: nil, updateMillisecond: 5000}
	p.save(p)
	return p
}

func BindPropertiesModels(config interface{}) *PropertiesConfig {
	p := &PropertiesConfig{nc: &PropertiesConfig{}, config: config, callback: nil, errCallback: nil, updateMillisecond: 5000}
	p.save(p)
	return p
}

func BindApolloModels(config interface{}) *ApolloConfig {
	p := &ApolloConfig{nc: &ApolloConfig{}, config: config}
	p.save(p)
	return p
}

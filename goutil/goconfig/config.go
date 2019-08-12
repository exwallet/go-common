/*
 * @Author: kidd
 * @Date: 5/8/19 1:28 PM
 */

package goconfig

import "regexp"

var regexSeps = regexp.MustCompile(`[;,]`)

type Configer interface {
	Set(key, val string) error   //support section::key type in given key when using ini type.
	String(key string) string    //support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	Strings(key string, seps... string) []string //get string slice
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DefaultString(key string, defaultVal string) string      // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	DefaultStrings(key string, defaultVal []string) []string //get string slice
	DefaultInt(key string, defaultVal int) int
	DefaultInt64(key string, defaultVal int64) int64
	DefaultBool(key string, defaultVal bool) bool
	DefaultFloat(key string, defaultVal float64) float64
	DIY(key string) (interface{}, error)
	GetSection(section string) (map[string]string, error)
	SaveConfigFile(filename string) error
}

// Config is the main struct for BConfig
type Config interface {
	Parse(key string) (Configer, error)
	ParseData(data []byte) (Configer, error)
}

func init() {

}

func NewIniConfig(filename string) (Configer, error) {
	config := &IniConfig{}
	return config.Parse(filename)
}

func NewJsonConfig(filename string) (Configer, error) {
	config := &JSONConfig{}
	return config.Parse(filename)
}

//
//
//// Register makes a config adapter available by the adapter name.
//// If Register is called twice with the same name or if driver is nil,
//// it panics.
//func Register(name string, adapter Config) {
//	if adapter == nil {
//		panic("config: Register adapter is nil")
//	}
//	if _, ok := adapters[name]; ok {
//		panic("config: Register called twice for adapter " + name)
//	}
//	adapters[name] = adapter
//}

//// NewConfig adapterName is ini/json/xml/yaml.
//// filename is the config file path.
//func NewConfig(adapterName, filename string) (Configer, error) {
//	adapter, ok := adapters[adapterName]
//	if !ok {
//		return nil, fmt.Errorf("config: unknown adaptername %q (forgotten import?)", adapterName)
//	}
//	return adapter.Parse(filename)
//}

//// NewConfigData adapterName is ini/json/xml/yaml.
//// data is the config data.
//func NewConfigData(adapterName string, data []byte) (Configer, error) {
//	adapter, ok := adapters[adapterName]
//	if !ok {
//		return nil, fmt.Errorf("config: unknown adaptername %q (forgotten import?)", adapterName)
//	}
//	return adapter.ParseData(data)
//}

//
//// now only support ini, next will support json.
//func parseConfig(appConfigPath string) (err error) {
//	AppConfig, err = newAppConfig(appConfigProvider, appConfigPath)
//	if err != nil {
//		return err
//	}
//	return assignConfig(AppConfig)
//}
//
//
//
//// now only support ini, next will support json.
//func parseConfig(appConfigPath string) (err error) {
//	AppConfig, err = newAppConfig(appConfigProvider, appConfigPath)
//	return err
//}

//func assignConfig(ac config.Configer) error {
//	for _, i := range []interface{}{BConfig, &BConfig.Listen, &BConfig.WebConfig, &BConfig.Log, &BConfig.WebConfig.Session} {
//		assignSingleConfig(i, ac)
//	}
//	// set the run mode first
//	if envRunMode := os.Getenv("BEEGO_RUNMODE"); envRunMode != "" {
//		BConfig.RunMode = envRunMode
//	} else if runMode := ac.String("RunMode"); runMode != "" {
//		BConfig.RunMode = runMode
//	}
//
//	if sd := ac.String("StaticDir"); sd != "" {
//		BConfig.WebConfig.StaticDir = map[string]string{}
//		sds := strings.Fields(sd)
//		for _, v := range sds {
//			if url2fsmap := strings.SplitN(v, ":", 2); len(url2fsmap) == 2 {
//				BConfig.WebConfig.StaticDir["/"+strings.Trim(url2fsmap[0], "/")] = url2fsmap[1]
//			} else {
//				BConfig.WebConfig.StaticDir["/"+strings.Trim(url2fsmap[0], "/")] = url2fsmap[0]
//			}
//		}
//	}
//
//	if sgz := ac.String("StaticExtensionsToGzip"); sgz != "" {
//		extensions := strings.Split(sgz, ",")
//		fileExts := []string{}
//		for _, ext := range extensions {
//			ext = strings.TrimSpace(ext)
//			if ext == "" {
//				continue
//			}
//			if !strings.HasPrefix(ext, ".") {
//				ext = "." + ext
//			}
//			fileExts = append(fileExts, ext)
//		}
//		if len(fileExts) > 0 {
//			BConfig.WebConfig.StaticExtensionsToGzip = fileExts
//		}
//	}
//
//	if lo := ac.String("LogOutputs"); lo != "" {
//		// if lo is not nil or empty
//		// means user has set his own LogOutputs
//		// clear the default setting to BConfig.Log.Outputs
//		BConfig.Log.Outputs = make(map[string]string)
//		los := strings.Split(lo, ";")
//		for _, v := range los {
//			if logType2Config := strings.SplitN(v, ",", 2); len(logType2Config) == 2 {
//				BConfig.Log.Outputs[logType2Config[0]] = logType2Config[1]
//			} else {
//				continue
//			}
//		}
//	}
//
//	//init log
//	logs.Reset()
//	for adaptor, config := range BConfig.Log.Outputs {
//		err := logs.SetLogger(adaptor, config)
//		if err != nil {
//			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s with the config %q got err:%s", adaptor, config, err.Error()))
//		}
//	}
//	logs.SetLogFuncCall(BConfig.Log.FileLineNum)
//
//	return nil
//}

//func assignSingleConfig(p interface{}, ac Configer) {
//	pt := reflect.TypeOf(p)
//	if pt.Kind() != reflect.Ptr {
//		return
//	}
//	pt = pt.Elem()
//	if pt.Kind() != reflect.Struct {
//		return
//	}
//	pv := reflect.ValueOf(p).Elem()
//
//	for i := 0; i < pt.NumField(); i++ {
//		pf := pv.Field(i)
//		if !pf.CanSet() {
//			continue
//		}
//		name := pt.Field(i).Name
//		switch pf.Kind() {
//		case reflect.String:
//			pf.SetString(ac.DefaultString(name, pf.String()))
//		case reflect.Int, reflect.Int64:
//			pf.SetInt(ac.DefaultInt64(name, pf.Int()))
//		case reflect.Bool:
//			pf.SetBool(ac.DefaultBool(name, pf.Bool()))
//		case reflect.Struct:
//		default:
//			//do nothing here
//		}
//	}
//
//}
//
//// LoadAppConfig allow developer to apply a config file
//func LoadAppConfig(adapterName, configPath string) error {
//	absConfigPath, err := filepath.Abs(configPath)
//	if err != nil {
//		return err
//	}
//
//	if !utils.FileExists(absConfigPath) {
//		return fmt.Errorf("the target config file: %s don't exist", configPath)
//	}
//
//	appConfigPath = absConfigPath
//	appConfigProvider = adapterName
//
//	return parseConfig(appConfigPath)
//}
//
//type beegoAppConfig struct {
//	innerConfig Configer
//}

//func newAppConfig(appConfigProvider, appConfigPath string) (*beegoAppConfig, error) {
//	ac, err := NewConfig(appConfigProvider, appConfigPath)
//	if err != nil {
//		return nil, err
//	}
//	return &beegoAppConfig{ac}, nil
//}
//
//func (b *beegoAppConfig) Set(key, val string) error {
//	if err := b.innerConfig.Set(BConfig.RunMode+"::"+key, val); err != nil {
//		return err
//	}
//	return b.innerConfig.Set(key, val)
//}
//
//func (b *beegoAppConfig) String(key string) string {
//	if v := b.innerConfig.String(BConfig.RunMode + "::" + key); v != "" {
//		return v
//	}
//	return b.innerConfig.String(key)
//}
//
//func (b *beegoAppConfig) Strings(key string) []string {
//	if v := b.innerConfig.Strings(BConfig.RunMode + "::" + key); len(v) > 0 {
//		return v
//	}
//	return b.innerConfig.Strings(key)
//}
//
//func (b *beegoAppConfig) Int(key string) (int, error) {
//	if v, err := b.innerConfig.Int(BConfig.RunMode + "::" + key); err == nil {
//		return v, nil
//	}
//	return b.innerConfig.Int(key)
//}
//
//func (b *beegoAppConfig) Int64(key string) (int64, error) {
//	if v, err := b.innerConfig.Int64(BConfig.RunMode + "::" + key); err == nil {
//		return v, nil
//	}
//	return b.innerConfig.Int64(key)
//}
//
//func (b *beegoAppConfig) Bool(key string) (bool, error) {
//	if v, err := b.innerConfig.Bool(BConfig.RunMode + "::" + key); err == nil {
//		return v, nil
//	}
//	return b.innerConfig.Bool(key)
//}
//
//func (b *beegoAppConfig) Float(key string) (float64, error) {
//	if v, err := b.innerConfig.Float(BConfig.RunMode + "::" + key); err == nil {
//		return v, nil
//	}
//	return b.innerConfig.Float(key)
//}
//
//func (b *beegoAppConfig) DefaultString(key string, defaultVal string) string {
//	if v := b.String(key); v != "" {
//		return v
//	}
//	return defaultVal
//}
//
//func (b *beegoAppConfig) DefaultStrings(key string, defaultVal []string) []string {
//	if v := b.Strings(key); len(v) != 0 {
//		return v
//	}
//	return defaultVal
//}
//
//func (b *beegoAppConfig) DefaultInt(key string, defaultVal int) int {
//	if v, err := b.Int(key); err == nil {
//		return v
//	}
//	return defaultVal
//}
//
//func (b *beegoAppConfig) DefaultInt64(key string, defaultVal int64) int64 {
//	if v, err := b.Int64(key); err == nil {
//		return v
//	}
//	return defaultVal
//}
//
//func (b *beegoAppConfig) DefaultBool(key string, defaultVal bool) bool {
//	if v, err := b.Bool(key); err == nil {
//		return v
//	}
//	return defaultVal
//}
//
//func (b *beegoAppConfig) DefaultFloat(key string, defaultVal float64) float64 {
//	if v, err := b.Float(key); err == nil {
//		return v
//	}
//	return defaultVal
//}
//
//func (b *beegoAppConfig) DIY(key string) (interface{}, error) {
//	return b.innerConfig.DIY(key)
//}
//
//func (b *beegoAppConfig) GetSection(section string) (map[string]string, error) {
//	return b.innerConfig.GetSection(section)
//}
//
//func (b *beegoAppConfig) SaveConfigFile(filename string) error {
//	return b.innerConfig.SaveConfigFile(filename)
//}
//
//

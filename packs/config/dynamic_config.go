package config

import (
	"log"
	"reflect"
	"sync"
	"../util"
)

const (
	ConstConfigEnv 			= "dev"					// dev or pro; dev will show some debug print.
	ConstConfigIniFile 		= "conf/static.ini"
	ConstConfigLogPath		= "/logs/access.log"
	ConstConfigDefaultValue	= ""
)

type DynamicSetting struct {
	DbType					string
	MysqlMaxIdelConns		string
	MysqlMaxOpenConns		string
	MysqlHost				string
	MysqlUser				string
	MysqlPass				string
	MysqlPort				string
	MysqlName				string
	MysqlChar				string
	HttpPort				string
	HttpPageSize			string
	lock					sync.Mutex
}

var DynamicSettingInstance = &DynamicSetting{}

func init() {
	log.Println("util.dynamic_config.init")
	//init ini's config
	IniConfig, err := util.NewFileReader(ConstConfigIniFile)
	if err != nil {
		log.Println("error: fail to load ini config file: " + ConstConfigIniFile)
		return
	}
	if !IniConfig.HasSection(ConstConfigEnv) {
		log.Println("error: fail to find ini config file: " + ConstConfigIniFile + "'s sectiin: " + ConstConfigEnv)
		return
	}
	DynamicSettingInstance.SetString("DbType", IniConfig.GetString(ConstConfigEnv, "db.type", "mysql"))
	DynamicSettingInstance.SetString("MysqlMaxIdelConns", IniConfig.GetString(ConstConfigEnv, "mysql.max_idle_conns", "100"))
	DynamicSettingInstance.SetString("MysqlMaxOpenConns", IniConfig.GetString(ConstConfigEnv, "mysql.max_open_conns", "200"))
	DynamicSettingInstance.SetString("MysqlHost", IniConfig.GetString(ConstConfigEnv, "mysql.host", "127.0.0.1"))
	DynamicSettingInstance.SetString("MysqlUser", IniConfig.GetString(ConstConfigEnv, "mysql.user", "root"))
	DynamicSettingInstance.SetString("MysqlPass", IniConfig.GetString(ConstConfigEnv, "mysql.passwd", ""))
	DynamicSettingInstance.SetString("MysqlPort", IniConfig.GetString(ConstConfigEnv, "mysql.port", "3306"))
	DynamicSettingInstance.SetString("MysqlName", IniConfig.GetString(ConstConfigEnv, "mysql.name", "center"))
	DynamicSettingInstance.SetString("MysqlChar", IniConfig.GetString(ConstConfigEnv, "mysql.char", "utf8"))
	DynamicSettingInstance.SetString("HttpPort",  IniConfig.GetString(ConstConfigEnv, "http.port", ":8080"))
	DynamicSettingInstance.SetString("HttpPageSize",  IniConfig.GetString(ConstConfigEnv, "http.pageSize", "10"))
}

func (ds *DynamicSetting) SetString(key, value string) {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	ref := reflect.ValueOf(ds).Elem()
	if ref.Kind() == reflect.Invalid {
		return
	}else{
		ref.FieldByName(key).SetString(value)
	}
}

func (ds *DynamicSetting) SetInt(key string, value int64) {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	ref := reflect.ValueOf(ds).Elem()
	if ref.Kind() == reflect.Invalid {
		return
	}else{
		ref.FieldByName(key).SetInt(value)
	}
}

/*func (ds *DynamicSetting) Set(key string, value interface{}) {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	ref := reflect.ValueOf(ds).Elem()
	if ref.Kind() == reflect.Invalid {
		return
	}else{
		ref.FieldByName(key).Set(value)
	}
}*/

func (ds *DynamicSetting) GetString(key, def string) string {
	ret := reflect.ValueOf(ds).Elem().FieldByName(key)

	if ret.Kind() == reflect.Invalid {
		return def
	}else{
		return ret.String()
	}
}

func (ds *DynamicSetting) GetInt(key string, def int) int {
	ret := reflect.ValueOf(ds).Elem().FieldByName(key)

	if ret.Kind() == reflect.Invalid {
		return def
	}else{
		return int(ret.Int())
	}
}

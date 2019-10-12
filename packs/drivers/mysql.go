package drivers

import (
	"database/sql"
	"log"
	"sync"
	_ "../../packs/gin/plugins/mysql"
	"../../packs/util"
)

type Mysql struct {
	instance *sql.DB
	lock	sync.Mutex
}

var this Mysql

// this init func will be executed when this package is imported
func init() {
	log.Println("driver.mysql.init")
	var err error

	db_user := util.DynamicSettingInstance.GetString("MysqlUser", util.ConstConfigDefaultValue)
	db_pwd  := util.DynamicSettingInstance.GetString("MysqlPass", util.ConstConfigDefaultValue)
	db_host := util.DynamicSettingInstance.GetString("MysqlHost", util.ConstConfigDefaultValue)
	db_port := util.DynamicSettingInstance.GetString("MysqlPort", util.ConstConfigDefaultValue)
	db_name := util.DynamicSettingInstance.GetString("MysqlName", util.ConstConfigDefaultValue)
	db_char := util.DynamicSettingInstance.GetString("MysqlChar", util.ConstConfigDefaultValue)

	dns := db_user + ":" + db_pwd + "@tcp(" + db_host + ":" + db_port + ")/" + db_name + "?charset=" + db_char + "&loc=Asia%2FShanghai"

	this.instance, err = sql.Open(util.DynamicSettingInstance.GetString("DbType", util.ConstConfigDefaultValue), dns)

	err_dping := this.instance.Ping()

	util.CheckErr(err_dping)

	db_idel	:= util.Str2int(util.DynamicSettingInstance.GetString("MysqlMaxIdelConns", util.ConstConfigDefaultValue))
	util.CheckErr(err)

	db_open := util.Str2int(util.DynamicSettingInstance.GetString("MysqlMaxOpenConns", util.ConstConfigDefaultValue))
	util.CheckErr(err)

	this.instance.SetMaxIdleConns(db_idel)


	this.instance.SetMaxOpenConns(db_open)
}

func (_ *Mysql) GetInstance() *sql.DB{
	return this.instance
}

func (_ *Mysql) GetAll(querySql string) (int, []map[string]string, error) {
	var rows *sql.Rows
	var err error
	var rnum int
	var container []map[string]string

	rows, err = this.instance.Query(querySql)
	defer rows.Close()

	if nil !=	err {
		return rnum, container, err
	}

	var rcol []string

	rcol, err = rows.Columns()
	if nil != err {
		return rnum, container, err
	}

	cnum := len(rcol)
	scaner	:=	make([]interface{}, cnum)
	value	:=	make([]interface{}, cnum)
	for j := range value {
		scaner[j] = &value[j]
	}

	index := 0
	for rows.Next() {
		err = rows.Scan(scaner...)
		container = append(container, make(map[string]string))
		for i, col := range value {
			rnum += 1
			container[index][rcol[i]] = string(col.([]byte))
		}
		index ++
	}

	return	rnum, container, nil
}

func (_ *Mysql) GetRow(querySql string) (*map[string]string, error) {
	ret := make(map[string]string)

	row, err := this.instance.Query(querySql)
	defer row.Close()
	if nil != err {
		return &ret, err
	}

	columns, err := row.Columns()
	if nil != err {
		return &ret, err
	}

	cnum := len(columns)
	scaner := make([]interface{}, cnum)
	values := make([]interface{}, cnum)

	for j := range values {
		scaner[j] = &values[j]
	}

	row.Next()
	err = row.Scan(scaner...)
	for i, col := range values {
		if nil != col {
			ret[columns[i]] = string(col.([]byte))
		}
	}

	return  &ret, nil
}

func (_ *Mysql) GetOne(querySql string) (string, error) {
	var tmp interface{}
	ret := ""

	err := this.instance.QueryRow(querySql).Scan(&tmp)

	if nil != err {
		return ret, err
	}

	ret = string(tmp.([]byte))

	return ret, err
}

func (_ *Mysql) DefGetOne(querySql, defaultStr string) (string) {
	var tmp interface{}

	err := this.instance.QueryRow(querySql).Scan(&tmp)

	if nil != err {
		return defaultStr
	}

	return string(tmp.([]byte))
}

func (_ *Mysql) Exec(querySql string) (sql.Result, error) {
	return this.instance.Exec(querySql)
}

package driver

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"sync"
)

var (
	MySqlDb = make(map[string]*MySqlDriver)
	mysqlMu sync.Mutex
)

type MySqlDriver struct {
	Connected bool
	Engine    *xorm.Engine
	cfg       *MysqlConfig
}

type MysqlConfig struct {
	DbHost         string `yaml:"dbHost"`
	DbUser         string `yaml:"dbUser"`
	DbPwd          string `yaml:"dbPwd"`
	DbName         string `yaml:"dbName"`
	DbMaxIdleConns int    `yaml:"dbMaxIdleConns"`
	DbMaxOpenConns int    `yaml:"dbMaxOpenConns"`
}

func RegisterMysql(name string, cfg *MysqlConfig) (err error) {
	mysqlMu.Lock()
	defer mysqlMu.Unlock()
	md := new(MySqlDriver)
	md.cfg = cfg
	err = md.Register(cfg)
	if err != nil {
		return
	}
	MySqlDb[name] = md
	return
}
func GetMysql(name string) (engine *xorm.Engine) {
	md := MySqlDb[name]
	if md == nil {
		return
	}
	if !md.Connected {
		md.Connect()
	}
	engine = md.Engine
	return
}

func (p *MySqlDriver) Register(cfg *MysqlConfig) (err error) {
	//todo cfg的基础判断
	p.cfg = cfg
	return
}

func (p *MySqlDriver) Connect() (err error) {
	dsn := p.cfg.DbUser + ":" + p.cfg.DbPwd + "@tcp(" + p.cfg.DbHost + ")/" + p.cfg.DbName + "?parseTime=true&charset=utf8"
	p.Engine, err = xorm.NewEngine("mysql", dsn)
	p.Engine.SetMaxIdleConns(p.cfg.DbMaxIdleConns)
	p.Engine.SetMaxOpenConns(p.cfg.DbMaxIdleConns)
	if err != nil {
		return
	}
	p.Connected = true
	return
}

func (p *MySqlDriver) Close() {
	p.Engine.Close()
	p.Connected = false
}

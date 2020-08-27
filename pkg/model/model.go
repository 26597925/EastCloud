package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sapi/pkg/logger"
	"sync"
)

var (
	DB *gorm.DB
	m      sync.RWMutex
	isInit bool
)

const (
	HostUri = "%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&loc=Local"
)

type Option func(*Options)

type Callback struct {}

//用来写入日志
func (callback Callback) Print(values ...interface{}) {
	logger.Info(values)
}

func Init (option *Options) (err error) {
	m.Lock()
	defer m.Unlock()

	if isInit {
		logger.Warn("已经初始化过数据库")
		return
	}

	logger.Debug(option)
	DB, err = gorm.Open(option.Driver, fmt.Sprintf(HostUri,
		option.Username,
		option.Password,
		option.Host,
		option.Port,
		option.Database,
		option.Charset,
		))

	if err == nil {
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return option.Prefix + defaultTableName
		}

		DB.LogMode(true)
		DB.SetLogger(Callback{})

		isInit = true
		logger.Info("初始化数据库连接完成")
	} else {
		logger.Error(err)
	}

	return  err
}
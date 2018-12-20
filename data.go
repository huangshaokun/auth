package authhandle

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	//Debug 测试开关， 打开之后，会打印错误
	Debug = false
)

const (
	//超级用户
	superUser     = "admin"
	superPassword = "superman"

	//在数据库中增加
	userPrefix      = "user/"
	authorityPrefix = "authority/"

	logfile = "authlog"
)

var (
	//表示有这个用户
	errHaveUser = fmt.Errorf("have the user")
)

//Auth leveldb
type Auth struct {
	db      *leveldb.DB
	jsonlog *zerolog.Logger
}

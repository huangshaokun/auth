package authhandle

import (
	"loghandle"
	"path"

	auth "github.com/abbot/go-http-auth"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

//newAuth 创建 leveldb 数据库实例， 一个路径只能调用一次
func newAuth(userPath, logPath string, superAuthority []byte) (au *Auth, err error) {
	au = &Auth{}

	//创建日志
	au.jsonlog = loghandle.NewJSONLog(path.Join(logPath, logfile), false)

	//打开数据库
	au.db, err = leveldb.OpenFile(userPath, nil)
	if err != nil {
		if Debug {
			loghandle.Debug(au.jsonlog).Err(err).Msg("")
		}
		return
	}

	//查询是否有超级用户，没有则生成默认的超级用户
	passwordKey := createUserPasswordKey(superUser)
	password, err := au.db.Get(passwordKey, nil)
	if err != nil {
		if Debug {
			loghandle.Debug(au.jsonlog).Err(err).Msg("")
		}
	}

	if len(password) == 0 {
		//生成默认的超级用户
		err = au.db.Put(passwordKey, []byte(superPassword), nil)
		if err != nil {
			if Debug {
				loghandle.Debug(au.jsonlog).Err(err).Msg("")
			}
			return
		}
	}

	//创建超级用户的权限
	err = au.updateUserAuthority(superUser, superAuthority)
	if err != nil {
		if Debug {
			loghandle.Debug(au.jsonlog).Err(err).Msg("")
		}
		return
	}

	return
}

//close 关闭数据库
func (pa *Auth) close() (err error) {
	err = pa.db.Close()
	return
}

//secret 查询密码
func (pa *Auth) secret(user, realm string) string {
	passwordKey := createUserPasswordKey(user)

	password, err := pa.db.Get(passwordKey, nil)
	if err != nil {
		loghandle.Error(pa.jsonlog).Err(err).Msg("")
		return ""
	}

	//如果是超级用户，但是没有查到密码就使用初始密码
	if len(password) == 0 {
		if user == superUser {
			return superPassword
		}
		return ""
	}
	return string(auth.MD5Crypt(password, []byte("dlPL2MqE"), []byte("$1$")))
}

//addUser 增加用户
func (pa *Auth) addUser(user, password string, authority []byte) (err error) {
	passwordKey := createUserPasswordKey(user)

	//检查是否有这个用户
	_, err = pa.db.Get(passwordKey, nil)
	if err != nil {
		err = pa.db.Put(passwordKey, []byte(password), nil)
		if err != nil {
			if Debug {
				loghandle.Debug(pa.jsonlog).Err(err).Msg("")
			}
			return
		}
		authoritykey := createUserAuthorityKey(user)
		err = pa.db.Put(authoritykey, authority, nil)
		if err != nil {
			if Debug {
				loghandle.Debug(pa.jsonlog).Err(err).Msg("")
			}
			return
		}
		return
	}

	//如果密码不为空，表示有这个用户
	return errHaveUser

}

//deleteUser 删除用户
func (pa *Auth) deleteUser(user string) (err error) {
	//删除用户
	passwordKey := createUserPasswordKey(user)
	err = pa.db.Delete(passwordKey, nil)
	if err != nil {
		if Debug {
			loghandle.Debug(pa.jsonlog).Err(err).Msg("")
		}
		return
	}

	//删除权限
	authoritykey := createUserAuthorityKey(user)
	err = pa.db.Delete(authoritykey, nil)
	if err != nil {
		if Debug {
			loghandle.Debug(pa.jsonlog).Err(err).Msg("")
		}
		return
	}
	return
}

//findUser 查找某个用户信息
func (pa *Auth) findUser(user string) (password string, authority []byte, err error) {
	passwordKey := createUserPasswordKey(user)

	pabyte, err := pa.db.Get(passwordKey, nil)
	if err != nil {
		if Debug {
			loghandle.Debug(pa.jsonlog).Err(err).Msg("")
		}
		return
	}
	password = string(pabyte)

	authoritykey := createUserAuthorityKey(user)
	authority, err = pa.db.Get(authoritykey, nil)
	if err != nil {
		if Debug {
			loghandle.Debug(pa.jsonlog).Err(err).Msg("")
		}
		return
	}
	return
}

//loadAllUser 查找所有用户信息
func (pa *Auth) loadAllUser(fcallback func(user, password string, authority []byte) bool) (err error) {
	//遍历用户
	iter := pa.db.NewIterator(util.BytesPrefix([]byte(userPrefix)), nil)
	for iter.Next() {
		user := string(iter.Key()[len([]byte(userPrefix)):])
		password := string(iter.Value())
		//获取权限
		authoritykey := createUserAuthorityKey(user)
		authority, geterr := pa.db.Get(authoritykey, nil)
		if geterr != nil {
			if Debug {
				loghandle.Debug(pa.jsonlog).Err(err).Msg("")
			}
			if err == nil {
				err = geterr
			}
			continue
		}
		if !fcallback(user, password, authority) {
			break
		}
	}

	//释放资源
	iter.Release()

	//返回累计的错误
	err = iter.Error()
	if err != nil {
		if Debug {
			loghandle.Debug(pa.jsonlog).Err(err).Msg("")
		}
		return
	}
	return
}

//updateUserPassword 修改用户密码
func (pa *Auth) updateUserPassword(user, password string) (err error) {
	//检查是否有这个用户
	passwordKey := createUserPasswordKey(user)
	_, err = pa.db.Get(passwordKey, nil)
	if err != nil {
		if Debug {
			loghandle.Debug(pa.jsonlog).Err(err).Msg("")
		}
		return
	}

	err = pa.db.Put(passwordKey, []byte(password), nil)
	if err != nil {
		if Debug {
			loghandle.Debug(pa.jsonlog).Err(err).Msg("")
		}
		return
	}
	return
}

//updateUserAuthority 修改用户权限
func (pa *Auth) updateUserAuthority(user string, authority []byte) (err error) {
	//检查是否有这个用户
	passwordKey := createUserPasswordKey(user)
	_, err = pa.db.Get(passwordKey, nil)
	if err != nil {
		if Debug {
			loghandle.Debug(pa.jsonlog).Err(err).Msg("")
		}
		return
	}

	authoritykey := createUserAuthorityKey(user)
	err = pa.db.Put(authoritykey, authority, nil)
	if err != nil {
		if Debug {
			loghandle.Debug(pa.jsonlog).Err(err).Msg("")
		}
		return
	}
	return
}

//isOperablePassword 判断是否可操作 密码
func (pa *Auth) isOperablePassword(operatedUser string, beOperated string) (operated bool) {
	//如果是此用户则可以操作
	if operatedUser == beOperated {
		return true
	}

	//如果是超级用户则可以操作
	if operatedUser == superUser {
		return true
	}

	return false
}

//isOperableAuthority 判断是否可操做 权限
func (pa *Auth) isOperableAuthority(operatedUser string, beOperated string) (operated bool) {
	//如果是超级用户则可以操作
	if operatedUser == superUser {
		return true
	}
	return false
}

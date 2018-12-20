package authhandle

//NewAuth 创建 leveldb 数据库实例， 一个路径只能调用一次
func NewAuth(userPath, logPath string, superAuthority []byte) (au *Auth, err error) {
	return newAuth(userPath, logPath, superAuthority)

}

//Close 关闭数据库
func (pa *Auth) Close() (err error) {
	return pa.close()
}

//Secret 查询密码
func (pa *Auth) Secret(user, realm string) string {
	return pa.secret(user, realm)
}

//AddUser 增加用户
func (pa *Auth) AddUser(user, password string, authority []byte) (err error) {
	return pa.addUser(user, password, authority)
}

//DeleteUser 删除用户
func (pa *Auth) DeleteUser(user string) (err error) {
	return pa.deleteUser(user)
}

//FindUser 查找某个用户信息
func (pa *Auth) FindUser(user string) (password string, authority []byte, err error) {
	return pa.findUser(user)
}

//LoadAllUser 查找所有用户信息
func (pa *Auth) LoadAllUser(fcallback func(user, password string, authority []byte) bool) (err error) {
	return pa.loadAllUser(fcallback)
}

//UpdateUserPassword 修改用户密码
func (pa *Auth) UpdateUserPassword(user, password string) (err error) {
	return pa.updateUserPassword(user, password)
}

//UpdateUserAuthority 修改用户权限
func (pa *Auth) UpdateUserAuthority(user string, authority []byte) (err error) {
	return pa.updateUserAuthority(user, authority)
}

//IsOperablePassword 判断是否可操作 密码
func (pa *Auth) IsOperablePassword(operatedUser string, beOperated string) (operated bool) {
	return pa.isOperablePassword(operatedUser, beOperated)
}

//IsOperableAuthority 判断是否可操做 权限
func (pa *Auth) IsOperableAuthority(operatedUser string, beOperated string) (operated bool) {
	return pa.isOperableAuthority(operatedUser, beOperated)
}

package authhandle

//创建用户密码的key
func createUserPasswordKey(userName string) []byte {
	return []byte(userPrefix + userName)
}

//创建用户的权限的key
func createUserAuthorityKey(userName string) []byte {
	return []byte(authorityPrefix + userName)
}

##功能
* 依赖谷歌默认的用户和密码选项
* 默认最高权限用户 admin, 密码 superman 
* 存数据库存储的是
    user/{username}:{密码}
    authority/{username}:{权限}
* 提供增删查改用户信息等接口， 内部不判断权限等，是否有权限调用此函数由使用者编写
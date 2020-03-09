package models

type Server struct {
	Port     int64  `default:30000`
	Hostname string `defualt:http://localhost/`
}
type Service struct {
	Threads int64  `deault:200`
	Counter string `deault:redis` // inner 使用内部还是Redis进行计数
}
type Redis struct {
	Redishost string `deault:127.0.0.1`
	Redisport int64  `deault:6379`
	Status    bool   `deault:false` // 是否初始化Redis计数器
}

type Config struct {
	Server  Server
	Service Service
	Redis   Redis
}

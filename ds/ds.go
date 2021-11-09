package ds

import "time"

const MaxConnAge = 1000 * time.Second
const MaxConnCount = 5

type Configs map[string]Config

type Config struct {
	Engine   string
	Host     string
	Name     string
	Port     string
	User     string
	Password string
	Secure bool
}
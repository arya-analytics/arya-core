package ds

import "time"

const MaxConnAge = 1000 * time.Second
const MaxConnCount = 5

type Config map[string] ConnParams

type ConnParams struct {
	Engine   string
	Host     string
	Name     string
	Port     string
	User     string
	Password string
	Secure bool
}
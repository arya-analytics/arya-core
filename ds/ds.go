package ds

import "time"

// || CONSTANTS ||

const MaxConnAge = 1000 * time.Second
const MaxConnCount = 5

// || CONFIG ||

type Configs map[string] Config

type Config struct {
	Engine   Engine
	Host     string
	Name     string
	Port     string
	Auth AuthConfig
	Secure bool
}

// || AUTHENTICATION ||

type AuthMethod int

const (
	TLS AuthMethod = iota
	Credentials
	Token
	None
)

type AuthConfig struct {
	Mode     AuthMethod
	User     string
	Password string
	Token    string
}

// || ENGINE ||

type Engine string

const (
	Postgres Engine = "github.com/uptrace/bun/driver/pgdriver"
	GorillaWS Engine = "github.com/gorilla/websocket"
)

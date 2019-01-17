package dsncfg

type ConnectionConfig struct {
	MaxOpenConnections int `json:"max_open_connections"`
	MaxIdleConnections int `json:"max_idle_connections"`
	MaxLifeTime        int `json:"max_life_time"` // max connection life time seconds
}

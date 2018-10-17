package dsncfg;

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

const (
	MySql      = "mysql"
	Postgresql = "postgres"
	Sqlite     = "sqlite"
)

// 0 - dbTypeMySql
// 1 - dbTypePgSql
// 3 - dbTypeSqLite
const (
	dbTypeMySql  = iota
	dbTypePgSql
	dbTypeSqLite
)

const (
	defaultMysqlPort      = 3306
	defaultPgsqlPort      = 5432
	defaultHost           = "localhost"
	defaultSqliteDataFile = "/tmp/sqlite.db"
)

// Database configuration
type Database struct {
	// Database port, only for tcp connection
	// Default:
	//   mysql: 3306
	//   pgsql: 5432
	Port int `json:"port"`
	// Database host. Use this paramater to set host address or
	// path to the unix socket or database file path
	// Default:
	//   mysql/pgsql: localhost
	//   sqlite: /tmp/sqlite.db
	Host string `json:"host"`
	// Database name, only for the MySql, PgSql
	Name string `json:"name"`
	// Use tcp or unix connection type for the MySql and PgSql.
	// Default: tcp
	Protocol string `json:"protocol"`
	// Database type: sqlite, mysql, postgres
	Type string `json:"type"`
	// Database user
	User string `json:"user"`
	// Databse password
	Password string `json:"password"`
	// Additional connection options in URI query string format. options get separated by &
	Parameters map[string]string `json:"parameters"`

	dbType int
}

var (
	ErrorUnsupportedDB    = errors.New("Unsupported database")
	ErrorUnsupportedProto = errors.New("Unsupported database protocol")
	ErrorDBNameRequired   = errors.New("Database name is required")
	ErrorUserNameRequired = errors.New("Database user is required")
)

// Validate fields according to the databae type
// Must be called instantly configuration was parsed
func (this *Database) Init() (err error) {
	if err = this.setdbType(); err != nil {
		return
	}

	if err = this.setProtocol(); err != nil {
		return
	}

	// rewrite empty host and port for the MySql and Postgresql
	if this.dbType != dbTypeSqLite {
		this.Host = strDefault(this.Host, defaultHost)

		if this.dbType == dbTypeMySql {
			this.Port = intDefault(this.Port, defaultMysqlPort)
		} else {
			this.Port = intDefault(this.Port, defaultPgsqlPort)
		}

		if isEmpty(this.User) {
			return ErrorUserNameRequired
		}

		if isEmpty(this.Name) {
			return ErrorDBNameRequired
		}
	} else {
		this.Host = strDefault(this.Host, defaultSqliteDataFile)
	}

	return nil
}

// Create the data source name string to pass itd to the database driver
func (this *Database) DSN() (dsn string) {
	if this.dbType != dbTypeSqLite {
		dsn = this.authString()
		dsn += "@"
	}

	dsn += this.sourceString()

	if this.dbType != dbTypeSqLite {
		dsn += "/" + this.Name
	}

	dsn += this.uriString()

	return
}

// Set database type flag
func (this *Database) setdbType() error {
	this.Type = strings.ToLower(this.Type)

	switch this.Type {
	case MySql:
		this.dbType = dbTypeMySql;
	case Postgresql:
		this.dbType = dbTypePgSql;
	case Sqlite:
		this.dbType = dbTypeSqLite;
	default:
		return ErrorUnsupportedDB;
	}

	return nil
}

// Set valid connection type: tcp or unix
func (this *Database) setProtocol() error {
	this.Protocol = strings.ToLower(this.Protocol)

	if this.dbType == dbTypeMySql || this.dbType == dbTypePgSql {
		this.Protocol = strDefault(this.Protocol, "tcp")

		if this.Protocol != "tcp" && this.Protocol != "unix" {
			return ErrorUnsupportedProto
		}
	}

	return nil
}

// Mysql and Postgresql connection:
// tcp - tcp(host:port)
// unix - unix(host)
// Sqlite will use host value as database file path
func (this *Database) sourceString() (source string) {
	if this.dbType == dbTypeSqLite {
		source = this.Host
	} else {
		if this.Protocol == "tcp" {
			source = this.Host + ":" + int2Str(this.Port)
		} else {
			source = this.Host
		}

		source = this.Protocol + "(" + source + ")"
	}

	return source
}

// Mysql and Postgresql authorization string (login:password)
func (this *Database) authString() (auth string) {
	if this.dbType == dbTypeSqLite {
		return
	}

	auth = this.User

	if !isEmpty(this.Password) {
		auth += ":" + this.Password
	}

	return
}

// Encoded paramaeters
func (this *Database) uriString() string {
	var uri = url.Values{}

	if this.Parameters == nil || len(this.Parameters) == 0 {
		return ""
	}

	for idx, val := range this.Parameters {
		uri.Add(idx, val)
	}

	return "?" + uri.Encode()
}

// Check if string is empty or int is 0
func isEmpty(in interface{}) bool {
	switch in.(type) {
	case string:
		if str := strings.Trim(in.(string), " "); str != "" {
			return false
		}

	case int:
		if in.(int) > 0 {
			return false
		}
	}
	return true
}

// Return string value or default if value is empty
func strDefault(val, def string) string {
	if isEmpty(val) {
		return def
	}

	return val
}

// Return int value or default if value is empty
func intDefault(val, def int) int {
	if isEmpty(val) {
		return def
	}

	return val
}

// Convert integer value to the string
func int2Str(v interface{}) (s string) {
	switch v.(type) {
	case string:
		s = v.(string)

	case int:
		s = strconv.Itoa(v.(int))
	}

	return
}

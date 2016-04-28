package dsncfg;

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

// 0 - MySql
// 1 - PgSql
// 3 - SqLite
const (
	MySql = iota
	PgSql
	SqLite
)

// Database configuration
type Database struct {
	// Database port, only for tcp connection
	Port int
	// Database host. Use this paramater to set host address or
	// path to the unix socket or database file path
	Host,
	// Database name, only for the MySql, PgSql
	Name,
	// Use tcp or unix connection type for the MySql and PgSql
	Protocol,
	// Database type: sqlite, mysql, pgsql
	Type,
	// Database user
	User,
	// Databse password
	Password string
	// Additional connection options in URI query string format. options get separated by &
	Parameters map[string]string

	dbType int
}

var (
	ErrorUnsupportedDB = errors.New("Unsupported database")
	ErrorUnsupportedProto = errors.New("Unsupported database protocol")
	ErrorDBNameRequired = errors.New("Database name is required")
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

	// rewrite empty host and port for the MySql and PgSql
	if this.dbType != SqLite {
		this.Host = strDefault(this.Host, "localhost")

		if this.dbType == MySql {
			this.Port = intDefault(this.Port, 3306)
		} else {
			this.Port = intDefault(this.Port, 5555)
		}

		if isEmpty(this.User) {
			return ErrorUserNameRequired
		}

		if isEmpty(this.Name) {
			return ErrorDBNameRequired
		}
	} else {
		this.Host = strDefault(this.Host, "/tpm/unknown/db/path")
	}

	return nil
}

// Create the data source name string to pass itd to the database driver
func (this *Database) DSN() (dsn string) {
	if this.dbType != SqLite {
		dsn = this.authString()
		dsn += "@"
	}

	dsn += this.sourceString()

	if this.dbType != SqLite {
		dsn += "/" + this.Name
	}

	dsn += this.uriString()

	return
}

// Set database type flag
func (this *Database) setdbType() error {
	this.Type = strings.ToLower(this.Type)

	switch this.Type {
		case "mysql": this.dbType = MySql;
		case "pgsql": this.dbType = PgSql;
		case "sqlite": this.dbType = SqLite;
		default: return ErrorUnsupportedDB;
	}

	return nil
}

// Set valid connection type: tcp or unix
func (this *Database) setProtocol() error {
	this.Protocol = strings.ToLower(this.Protocol)

	if this.dbType == MySql || this.dbType == PgSql {
		this.Protocol = strDefault(this.Protocol, "tcp")

		if this.Protocol != "tcp" && this.Protocol != "unix" {
			return ErrorUnsupportedProto
		}
	}

	return nil
}

// Mysql and PgSql connection:
// tcp - tcp(host:port)
// unix - unix(host)
// SqLite will use host value as database file path
func (this *Database) sourceString() (source string) {
	if this.dbType == SqLite {
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

// MySql and PgSql authorization string (login:password)
func (this *Database) authString() (auth string) {
	if this.dbType == SqLite {
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

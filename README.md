## DSN (Data Source Name) helper

Save Your time develop in the database configuration

Good service should have configuration file especially with human readable format: toml, ini etc...
This micro library helps to convert database configuration section to the data source name string

### Configuration object

| Field          | Description                                              |
| :------------- | :------------------------------------------------------- |
| Pot            | Database port, only for tcp connection                   |
| Host           | Database host. Use this paramater to set host address or |
|                | path to the unix socket or database file path            |
| Name           | Database name, only for the MySql, PgSql                 |
| Protocol       | Use tcp or unix connection type for the MySql and PgSql  |
| Type           | Database type: sqlite, mysql, pgsql                      |
| User           | Database user                                            |
| Password       | Databse password                                         |
| Parameters     | Additional connection options in URI                     |
|                | query string format. options get separated by &          |

### Usage

Project code:

```go
pacakge main

import (
	"github.com/BurntSushi/toml"
	"dsncfg"
	"fmt"
)

type Config struct {
	Database *dsncfg.Database{} `toml:"database"`
}


func main() {
	var (
		cfg *Config
		err error
	)

	if cfg, err = ReadConfig(); err != nil {
		fmt.Error(err.Error())
	} else {
		fmt.Println(cfg.Database.DSN())
	}
}

func ReadConfig() (cfg *Config, err error) {
	cfg = &Config{
		Database: &dsncfg.Database{},
	}

	if _, err = toml.DecodeFile("/path/to/config.toml", cfg); err != nil {
		return
	}

	err = this.Database.Init()

	return
}
```

Config:

```toml
[database]
host = "my.host"
user = "any"
name = "dbname"
  [database.parameters]
  charset = "utf8"
```

Result:

```
any@tcp(my.host:3306)/dbname?charset=utf8
```

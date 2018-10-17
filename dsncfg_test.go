package dsncfg;

import (
	"testing"
)

type DatabaseMock struct {
	Database
	Valid string
}

func Test_DSNValid(t *testing.T) {
	var (
		cfg  *Database
		mock = []DatabaseMock{
			DatabaseMock{
				Database: Database{Type: "sqlite", Host: "/some/database/file.db", Name: "anydb", User: "bla", Port: 432},
				Valid:    "/some/database/file.db",
			},
			DatabaseMock{
				Database: Database{Type: "sqlite"},
				Valid:    defaultSqliteDataFile,
			},
			DatabaseMock{
				Database: Database{Type: "mysql", User: "User", Name: "DB"},
				Valid:    "User@tcp(localhost:3306)/DB",
			},
			DatabaseMock{
				Database: Database{Type: "postgres", User: "User", Name: "DB"},
				Valid:    "User@tcp(localhost:5432)/DB",
			},
			DatabaseMock{
				Database: Database{Type: "postgres", User: "User", Port: 65101, Name: "DB"},
				Valid:    "User@tcp(localhost:65101)/DB",
			},
			DatabaseMock{
				Database: Database{Type: "postgres", Protocol: "Unix", User: "User", Password: "$anyPw334!", Name: "DB"},
				Valid:    "User:$anyPw334!@unix(localhost)/DB",
			},
		}
	)

	for _, v := range mock {
		cfg = &Database{
			Type:     v.Type,
			Protocol: v.Protocol,
			Host:     v.Host,
			Port:     v.Port,
			User:     v.User,
			Password: v.Password,
			Name:     v.Name,
		}

		if err := cfg.Init(); err != nil {
			t.Errorf("Unexpected error %s", err.Error())
			continue
		}

		if str := cfg.DSN(); str != v.Valid {
			t.Errorf("expecting %s, but got %s", v.Valid, str)
		}
	}
}

func Test_getAuthStringValid(t *testing.T) {
	var (
		cfg  *Database
		mock = []DatabaseMock{
			DatabaseMock{
				Database: Database{Type: "mysql", User: "User", Name: "DB"},
				Valid:    "User",
			},
			DatabaseMock{
				Database: Database{Type: "postgres", User: "User", Password: "&dhg(0saq1", Name: "DB"},
				Valid:    "User:&dhg(0saq1",
			},
		}
	)

	for _, v := range mock {
		cfg = &Database{
			Type:     v.Type,
			User:     v.User,
			Password: v.Password,
			Name:     v.Name,
		}

		if err := cfg.Init(); err != nil {
			t.Errorf("Unexpected error %s", err.Error())
			continue
		}

		if str := cfg.authString(); str != v.Valid {
			t.Errorf("expecting %s, but got %s", v.Valid, str)
		}
	}
}

func Test_getSourceStringValid(t *testing.T) {
	var (
		cfg  *Database
		mock = []DatabaseMock{
			DatabaseMock{
				Database: Database{Type: "sqlite", Host: "/some/database/file.db"},
				Valid:    "/some/database/file.db",
			},
			DatabaseMock{
				Database: Database{Type: "sqlite"},
				Valid:    defaultSqliteDataFile,
			},
			DatabaseMock{
				Database: Database{Type: "mysql", User: "User", Name: "DB"},
				Valid:    "tcp(localhost:3306)",
			},
			DatabaseMock{
				Database: Database{Type: "postgres", User: "User", Name: "DB"},
				Valid:    "tcp(localhost:5432)",
			},
			DatabaseMock{
				Database: Database{Type: "postgres", User: "User", Port: 65101, Name: "DB"},
				Valid:    "tcp(localhost:65101)",
			},
			DatabaseMock{
				Database: Database{Type: "postgres", Protocol: "Unix", User: "User", Name: "DB"},
				Valid:    "unix(localhost)",
			},
		}
	)

	for _, v := range mock {
		cfg = &Database{
			Type:     v.Type,
			Protocol: v.Protocol,
			Host:     v.Host,
			Port:     v.Port,
			User:     v.User,
			Name:     v.Name,
		}

		if err := cfg.Init(); err != nil {
			t.Errorf("Unexpected error %s", err.Error())
			continue
		}

		if str := cfg.sourceString(); str != v.Valid {
			t.Errorf("expecting %s, but got %s", v.Valid, str)
		}
	}
}

func Test_setProtocolUnknown(t *testing.T) {
	var (
		cfg  *Database
		mock = []DatabaseMock{
			DatabaseMock{
				Database: Database{Type: "mysql", Protocol: "abc"},
			},
			DatabaseMock{
				Database: Database{Type: "postgres", Protocol: "ynix"},
			},
		}
	)

	for _, v := range mock {
		cfg = &Database{Type: v.Type, Protocol: v.Protocol}
		cfg.setdbType()

		if err := cfg.setProtocol(); err == nil || err != ErrorUnsupportedProto {
			t.Errorf("Expecting error %s, but got %v", ErrorUnsupportedProto.Error(), err)
		}
	}
}

func test_setProtocol(t *testing.T) {
	var (
		cfg  *Database
		mock = []DatabaseMock{
			DatabaseMock{
				Database: Database{Type: "mysql", Protocol: "unix"},
				Valid:    "unix",
			},
			DatabaseMock{
				Database: Database{Type: "postgres", Protocol: "tcp"},
				Valid:    "tcp",
			},
			DatabaseMock{
				Database: Database{Type: "mysql"},
				Valid:    "tcp",
			},
		}
	)

	for _, v := range mock {
		cfg = &Database{Type: v.Type, Protocol: v.Protocol}

		if err := cfg.setdbType(); err != nil {
			t.Errorf("Unexpected err %s", err.Error())
		}

		if err := cfg.setProtocol(); err == nil || err != ErrorUnsupportedProto {
			t.Errorf("Expecting error %s, but got %v", ErrorUnsupportedProto.Error(), err)
		}

		if cfg.Protocol != v.Valid {
			t.Errorf("Expecting valid protocol valud %s, but got %s", v.Valid, cfg.Protocol)
		}
	}
}

func Test_setdbTypeUnknownType(t *testing.T) {
	var cfg = &Database{Type: "UNKNOWN"}

	if err := cfg.setdbType(); err != ErrorUnsupportedDB {
		t.Errorf("Expected error %s, but got %v", ErrorUnsupportedDB.Error(), err);
	}
}

func Test_setdbTypeValid(t *testing.T) {
	type Mock struct {
		Type  string
		Valid int
	}

	var (
		cfg  *Database
		mock = []Mock{
			Mock{"mysql", dbTypeMySql},
			Mock{"postgres", dbTypePgSql},
			Mock{"sqlite", dbTypeSqLite},
		}
	)

	for _, v := range mock {
		cfg = &Database{Type: v.Type}

		if err := cfg.setdbType(); err != nil {
			t.Errorf("Unexpected error %s for the type %s", ErrorUnsupportedDB.Error(), v.Type);
		}

		if v.Valid != cfg.dbType {
			t.Errorf("Unexpected type at dbType on %s", v.Type)
		}
	}
}

func Test_isEmptyStringTrue(t *testing.T) {
	var (
		mock = []string{
			"",
			" ",
		}
	)

	for _, v := range mock {
		if isEmpty(v) != true {
			t.Errorf("Expecting false, but got true (%v)", v)
		}
	}
}

func Test_isEmptyIntTrue(t *testing.T) {
	var (
		mock = []int{
			0,
			-1,
		}
	)

	for _, v := range mock {
		if isEmpty(v) != true {
			t.Errorf("Expecting false, but got true (%v)", v)
		}
	}
}

func Test_isEmptyStringFalse(t *testing.T) {
	var (
		mock = []string{
			"Some data",
		}
	)

	for _, v := range mock {
		if isEmpty(v) != false {
			t.Error("Expecting true, but got false")
		}
	}
}

func Test_isEmptyIntFalse(t *testing.T) {
	var (
		mock = []int{
			123,
			3306,
		}
	)

	for _, v := range mock {
		if isEmpty(v) != false {
			t.Error("Expecting true, but got false")
		}
	}
}

func Test_strDefault(t *testing.T) {
	type Mock struct {
		origin, def, valid string
	}
	var (
		mock = []Mock{
			Mock{"", "Value A default", "Value A default"},
			Mock{"Value A", "Value B default", "Value A"},
		}
	)

	for _, v := range mock {
		if str := strDefault(v.origin, v.def); str != v.valid {
			t.Errorf("Expecting value %s, but got %s", v.valid, str)
		}
	}
}

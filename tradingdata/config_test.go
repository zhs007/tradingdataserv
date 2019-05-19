package tradingdata

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("../test/test.yaml")
	if err != nil {
		t.Fatalf("TestLoadConfig %v", err)

		return
	}

	if cfg.BindAddr != "127.0.0.1:7051" {
		t.Fatalf("TestLoadConfig BindAddr %v", cfg.BindAddr)

		return
	}

	if cfg.HTTPAddr != "127.0.0.1:7053" {
		t.Fatalf("TestLoadConfig HTTPAddr %v", cfg.HTTPAddr)

		return
	}

	if cfg.AnkaDB.DBPath != "./dat" {
		t.Fatalf("TestLoadConfig AnkaDB.DBPath %v", cfg.AnkaDB.DBPath)

		return
	}

	if cfg.AnkaDB.Engine != "leveldb" {
		t.Fatalf("TestLoadConfig AnkaDB.Engine %v", cfg.AnkaDB.Engine)

		return
	}

	t.Logf("TestLoadConfig is OK")
}

func TestCheckConfig(t *testing.T) {

	type testData struct {
		cfg Config
		err error
	}

	lst := []testData{
		testData{
			cfg: Config{
				BindAddr: "127.0.0.1:1234",
				HTTPAddr: "127.0.0.1:1235",
				AnkaDB: AnkaDBConfig{
					DBPath: "./data",
					Engine: "leveldb",
				},
			},
			err: nil,
		},
		testData{
			cfg: Config{
				BindAddr: "127.0.0.1:1234",
				HTTPAddr: "127.0.0.1:1235",
			},
			err: ErrNoAnkaDBConfig,
		},
		testData{
			cfg: Config{
				BindAddr: "127.0.0.1:1234",
				AnkaDB: AnkaDBConfig{
					DBPath: "./data",
					Engine: "leveldb",
				},
			},
			err: ErrNoHTTPServerAddr,
		},
		testData{
			cfg: Config{
				HTTPAddr: "127.0.0.1:1235",
				AnkaDB: AnkaDBConfig{
					DBPath: "./data",
					Engine: "leveldb",
				},
			},
			err: ErrNoBindAddress,
		},
	}

	for _, v := range lst {
		err := checkConfig(&v.cfg)
		if err != v.err {
			t.Fatalf("TestCheckConfig (%v != %v)", err, v.err)

			return
		}
	}

	t.Logf("TestCheckConfig is OK")
}

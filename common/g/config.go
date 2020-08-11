package g

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-transfer/utils"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
)

var (
	cfg            = flag.String("c", "./config/cfg.json", "specify config file")
	v              = flag.Bool("v", false, "show version")
	vv             = flag.Bool("vv", false, "show version detail")
	ConfigFile     string
	configFileLock = new(sync.RWMutex)
	config         *GlobalConfig
)

type LogConfig struct {
	LogPath      string `json:"log_path"`
	LogLevel     string `json:"log_level"`
	LogFileName  string `json:"log_file_name"`
	LogKeepHours int    `json:"log_keep_hours"`
}

type RpcConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type JudgeConfig struct {
	Enabled     bool              `json:"enabled"`
	Batch       int               `json:"batch"`
	ConnTimeout int               `json:"conn_timeout"`
	CallTimeout int               `json:"call_timeout"`
	MaxConn     int               `json:"max_conn"`
	MaxIdle     int               `json:"max_idle"`
	Replicas    int               `json:"replicas"`
	Cluster     map[string]string `json:"cluster"`
}

type GraphConfig struct {
	Enabled     bool              `json:"enabled"`
	Batch       int               `json:"batch"`
	ConnTimeout int               `json:"conn_timeout"`
	CallTimeout int               `json:"call_timeout"`
	MaxConn     int               `json:"max_conn"`
	MaxIdle     int               `json:"max_idle"`
	Replicas    int               `json:"replicas"`
	Cluster     map[string]string `json:"cluster"`
}

type GlobalConfig struct {
	Debug      bool         `json:"debug"`
	Log        *LogConfig   `json:"log"`
	Judge      *JudgeConfig `json:"judge"`
	Graph      *GraphConfig `json:"graph"`
	Rpc        *RpcConfig   `json:"rpc"`
	Http       *HttpConfig  `json:"http"`
	MaxCPURate float64      `json:"max_cpu_rate"`
	MaxMemRate float64      `json:"max_mem_rate"`
}

func Conf() *GlobalConfig {
	configFileLock.RLock()
	defer configFileLock.RUnlock()

	return config
}

func InitConfig() {
	var (
		cfgFile   string
		bs        []byte
		err       error
		maxMemMB  int
		maxCPUNum int
	)

	flag.Parse()

	if *v {
		fmt.Println(VersionInfo())
		os.Exit(0)
	}

	if *vv {
		fmt.Println(HbsInfo())
		os.Exit(0)
	}

	cfgFile = *cfg
	ConfigFile = cfgFile

	if cfgFile == "" {
		dlog.Fatal("config file not specified: use -c $filename")
	}

	if _, err = os.Stat(cfgFile); os.IsNotExist(err) {
		dlog.Fatalf("config file specified not found: %s", cfgFile)
	} else {
		dlog.Infof("use config file: %s", ConfigFile)
	}

	if bs, err = ioutil.ReadFile(cfgFile); err != nil {
		dlog.Fatalf("read config file failed: %s", err.Error())
	} else {
		if err = json.Unmarshal(bs, &config); err != nil {
			dlog.Fatalf("decode config file failed: %s", err.Error())
		} else {
			dlog.Infof("load config success from %s", cfgFile)
		}
	}

	if err = Validator(); err != nil {
		dlog.Errorf("validator config file fail: %s", err)
		os.Exit(127)
	}

	// 最大使用内存
	maxMemMB = utils.CalculateMemLimit(config.MaxMemRate)

	// 最大cpu核数
	maxCPUNum = utils.GetCPULimitNum(config.MaxCPURate)

	dlog.Infof("bind [%d] cpu core", maxCPUNum)
	runtime.GOMAXPROCS(maxCPUNum)

	dlog.Infof("memory limit: %d MB", maxMemMB)
}

func ReloadConfig() error {
	var (
		bs  []byte
		err error
	)

	if _, err = os.Stat(ConfigFile); os.IsNotExist(err) {
		dlog.Fatalf("config file specified not found: %s", ConfigFile)
		return err
	} else {
		dlog.Infof("reload config file: %s", ConfigFile)
	}

	if bs, err = ioutil.ReadFile(ConfigFile); err != nil {
		dlog.Fatalf("reload config file failed: %s", err)
		return err
	} else {
		configFileLock.RLock()
		defer configFileLock.RUnlock()

		if err = json.Unmarshal(bs, &config); err != nil {
			dlog.Fatalf("decode config file failed: %s", err)
			return err
		} else {
			dlog.Infof("reload config success from %s", ConfigFile)
		}
	}

	if err = Validator(); err != nil {
		dlog.Errorf("validator config file fail: %s", err)
		return err
	}

	return nil
}

func Validator() error {
	// 设置默认日志路径为 ./logs
	if config.Log.LogPath == "" {
		config.Log.LogPath = "./logs"
	}

	// 设置默认日志文件名称为 run.log
	if config.Log.LogFileName == "" {
		config.Log.LogFileName = "run.log"
	}

	// 设置默认日志级别为 LogLevel
	if config.Log.LogLevel == "" {
		config.Log.LogLevel = "INFO"
	}

	// 设置默认保留24小时的日志
	if config.Log.LogKeepHours == 0 {
		config.Log.LogKeepHours = 24
	}

	// judge 设置
	if config.Judge.Enabled {
		if config.Judge.Batch < 0 {
			config.Judge.Batch = 200
		}

		if config.Judge.MaxConn < 0 {
			config.Judge.MaxConn = 32
		}

		if config.Judge.MaxIdle < 0 {
			config.Judge.MaxIdle = 32
		}

		if config.Judge.CallTimeout < 0 {
			config.Judge.CallTimeout = 5000
		}

		if config.Judge.ConnTimeout < 0 {
			config.Judge.ConnTimeout = 1000
		}

		if config.Judge.Replicas < 0 {
			config.Judge.Replicas = 500
		}

		if len(config.Judge.Cluster) == 0 {
			return errors.New("judge cluster is empty")
		}
	}

	// graph 设置
	if config.Graph.Enabled {
		if config.Graph.Batch < 0 {
			config.Graph.Batch = 200
		}

		if config.Graph.MaxConn < 0 {
			config.Graph.MaxConn = 32
		}

		if config.Graph.MaxIdle < 0 {
			config.Graph.MaxIdle = 32
		}

		if config.Graph.CallTimeout < 0 {
			config.Graph.CallTimeout = 5000
		}

		if config.Graph.ConnTimeout < 0 {
			config.Graph.ConnTimeout = 1000
		}

		if config.Graph.Replicas < 0 {
			config.Graph.Replicas = 500
		}

		if len(config.Graph.Cluster) == 0 {
			return errors.New("graph cluster is empty")
		}
	}

	// rpc 设置
	if config.Rpc.Enabled {
		if config.Rpc.Listen == "" {
			return errors.New("rpc listen addr is empty")
		}
	}

	// http 设置
	if config.Http.Enabled {
		if config.Http.Listen == "" {
			return errors.New("http listen addr is empty")
		}
	}

	// MaxCPURate
	if config.MaxCPURate < 0 || config.MaxCPURate > 1 {
		return errors.New("max_cpu_rate is range 0 to 1")
	}

	// MaxMemRate
	if config.MaxMemRate < 0 || config.MaxMemRate > 1 {
		return errors.New("max_mem_rate is range 0 to 1")
	}

	return nil
}

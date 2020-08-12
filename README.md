# ops-transfer

基于 open-falcon 二次开发的 transfer 组件

## 特性
- 代码重构, 移除无用的 http 接口;
- 移除 influxdb 支持;
- 新增运行日志配置, 支持日志滚动;
- 新增 cpu 核心绑定、内存阈值配置; 当 agent 内存达到阈值的 50% 时, 打印告警信息；当内存达到阈值的 100%, 程序直接退出;


## 编译

it is a golang classic project

``` shell
cd $GOPATH/src/github.com/fanghongbo/ops-transfer/
./control build
./control start
```

## 配置
Refer to `cfg.example.json`, modify the file name to `cfg.json` :

```config
{
  "debug": true,
  "min_step": 30,
  "log": {
    "log_level": "INFO",
    "log_path": "./logs",
    "log_file_name": "run.log",
    "log_keep_hours": 3
  },
  "judge": {
    "enabled": true,
    "batch": 200,
    "conn_timeout": 1000,
    "call_timeout": 5000,
    "max_conn": 32,
    "max_idle": 32,
    "replicas": 500,
    "cluster": {
      "judge-00": "127.0.0.1:6080"
    }
  },
  "graph": {
    "enabled": true,
    "batch": 200,
    "conn_timeout": 1000,
    "call_timeout": 5000,
    "max_conn": 32,
    "max_idle": 32,
    "replicas": 500,
    "cluster": {
      "graph-00": "127.0.0.1:6070"
    }
  },
  "tsdb": {
    "enabled": false,
    "batch": 200,
    "conn_timeout": 1000,
    "call_timeout": 5000,
    "max_conn": 32,
    "max_idle": 32,
    "retry": 3,
    "address": "127.0.0.1:8088"
  },
  "http": {
    "enabled": true,
    "listen": ":6060"
  },
  "rpc": {
    "enabled": true,
    "listen": ":8433"
  },
  "max_cpu_rate": 0.2,
  "max_mem_rate": 0.3
}
```

## License

This software is licensed under the Apache License. See the LICENSE file in the top distribution directory for the full license text.

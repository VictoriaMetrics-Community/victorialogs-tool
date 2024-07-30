# victorialogs-tool
[English](https://github.com/VictoriaMetrics-Community/victorialogs-tool) | 简体中文

![image_01](image_01.jpg)

一个用于查询 victoria-logs 的工具(你会爱上它)

基于toml配置文件查询，将查询结果集输出至终端

## 特性

- 简单易用的命令行界面
- 多种组合查询方式
- 支持大范围的时间查询(几天内的日志数据也可手到擒来，甚至更多)
- 基于toml配置文件, 可多个配置文件切换查询，类似于session存储
- 结果集输出至终端，可发挥你的想象力，配合`grep`、`awk`、`>`等工具使用

## 安装

确保你的电脑已经安装了Go环境

选项一:

```bash
go install github.com/VictoriaMetrics-Community/victorialogs-tool@latest && mv $GOPATH/bin/victorialogs-tool $GOPATH/bin/vtool
```

选项二:

```bash
git clone git@github.com:VictoriaMetrics-Community/victorialogs-tool.git
make build
```

## 用法

```bash
$ ./vtool --help
A wonderful query tool for Victorialogs

Usage:
  victorialogs-tool [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  query       query logs from victoriametrics
  setcfg      Set up the configuration file for query logs

Flags:
  -h, --help     help for victorialogs-tool
  -t, --toggle   Help message for toggle

Use "victorialogs-tool [command] --help" for more information about a command.
```

1. 首先使用`vtool setcfg`命令设置配置文件

```bash
vtool setcfg xxx/i-love-coding.toml
```

该命令会在家目录下生成一个`vtool.json`的文件, 用于存储配置文件的路径

配置文件参考: https://github.com/VictoriaMetrics-Community/victorialogs-tool/blob/master/cfgs/example.toml

2. 运行`vtool query`即可

## 贡献

欢迎提交PR

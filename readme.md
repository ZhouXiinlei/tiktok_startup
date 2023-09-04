# Tiktok-Startup

本项目为第六届字节青训营后端基础班结营大项目，采用Go语言开发，主要技术栈为`gozero`+`gorm`。

团队汇报文档： [Tik-Start Docs](https://nsb6rfg8gy.feishu.cn/docx/RgSndkUG7oDisJxwKIacs34Xn2e?from=from_copylink "抖音启动团队飞书")


Apifox文档：[Tik-Start API](https://apifox.com/apidoc/shared-8f2e87c0-89ba-4c9b-9d15-e07237edba6d "抖音启动团队API FOX")

1024 Code：[Tik-Start Homepage](https://1024code.com/~rmulj3k "抖音启动团队主页")

## 目录

- [背景](#背景)
- [安装](#安装)
- [功能特点](#功能特点)
- [贡献](#贡献)
- [许可证](#许可证)

## 背景

`TikTok_Startup`项目是第六届字节青训营后端基础班的结营大项目，由来自杭州电子科技大学以及西安交通大学的学生团队`抖音，启动`开发。

## 启动服务

在开始之前，请确保你已经安装了Go以及Python，并且已经安装了 `subprocess` 和 `os` 库。虽然 `subprocess` 和 `os` 是 Python 的内置库，但仍然建议检查其可用性。

1. 首先，确保Go以及Python 安装正确：

```bash
python --version
```
```bash
go version
```

2.检查 `subprocess` 和 `os` 库是否可以在你的 Python 环境中导入：

```bash
python -c "import subprocess, os"
```

如果没有错误信息，那么你就可以正常使用了。

3.运行脚本
```bash
python start_services.py
```

4.检测状态

如果弹出四个Windows Shell，证明脚本启动成功。四个Shell分别对应：`contact.go`,`user.go`,`video.go`以及`tik_start.go`

不改变默认配置的情况下，服务将启动于`0.0.0.0:8888`
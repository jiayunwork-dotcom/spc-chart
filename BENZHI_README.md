Go 实现的制造业统计过程控制（SPC）分析服务，默认启动 HTTP 在 :8080 提供控制图与过程能力计算接口，也可通过子命令 ichart/xbar-r/cusum/ewma/capability/rules 在命令行对测量数据做离线分析。

## 构建与启动

```bash
go build -o spc-chart .
./spc-chart            # 启动 HTTP 服务 :8080
./spc-chart ichart -in example/measurements.txt -sigma 3
```

## 评测镜像

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh spc-chart
```

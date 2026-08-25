# spc-chart 示例

对一组过程测量值做最基础的「单值控制图（I-chart）」分析：给出中心线（均值）、
总体标准差、σ 控制限，并标出越界（失控）点。

本目录 `measurements.txt` 是 25 个测量值，绝大多数在 10 附近，最后一个是明显离群值（25.40）。

## 运行（默认 3σ 控制限）

```bash
go run . -in example/measurements.txt
```

预期：均值约 10.3，最后一点 25.40 被标为 out-of-control。

## 调整控制限宽度

```bash
go run . -in example/measurements.txt -sigma 2
```

> 缺 `-in`、文件不存在、或数据为空/无法解析 → 受控报错（exit 1），不 panic。
> 数值可用换行 / 空格 / 制表符 / 逗号分隔（支持单列 CSV）。

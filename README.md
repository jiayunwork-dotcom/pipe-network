# pipe-network

供水管网稳态水力求解器。提交节点-管段拓扑、水箱水头和节点需水量，后端解出各管流量与各点水头。水力内核，不是 SCADA 监控台。

## 摩阻模型

全程使用 **Hazen-Williams**（SI 单位）：

```
hf = 10.67 · L · Q^1.852 / (C^1.852 · D^4.87)
```

- `hf`：管段水头损失（m），与流量同向
- `L`：长度（m），`D`：管径（m），`C`：Hazen-Williams 粗糙系数
- `Q`：流量（m³/s）

所有摩阻计算集中在 `internal/hydraulics` 的 `Resistance` / `HeadLoss` / `FlowFromHead`，求解与漏损试算共用同一来源。

## 求解方法

主求解为 **全局梯度法（Global Gradient Algorithm, Todini & Pilati 1988）**，即 `hydraulics.Solve`：Newton 型迭代，每步同时以管流 Q（m 条）与非固定节点水头 H（n 条）为未知量，联立能量方程与节点连续性方程组成 (m+n) 线性系统求解。相比只以水头为未知量的节点牛顿法，GGA 的 (m+n) 耦合系统在环网里不会奇异、数值更稳健；对合法拓扑（每个需水节点都连通到固定水头源）稳定收敛。另提供 **Hardy-Cross 环法**（`hydraulics.HardyCross`）作为能量闭合对照，二者调用同一摩阻函数，收敛到同一物理解。

- 容差与迭代上限可调（默认 `tolerance=1e-8`，`max_iter=300`）
- 不收敛返回 `no_convergence` 错误，不会死循环
- 环网满足环路能量闭合：`Σ 绕环损失 ≈ 0`（见 `energy_test.go`）

## 校验

非法拓扑返回明确错误：悬空节点、未知端点、零管径、零长度、无水箱/固定水头点、需求节点不可达源（`network/validate.go`）。

## 漏损试算

`internal/leak` 支持在某节点加漏点：

- `demand` 模式：漏点作为额外需水量并入节点
- `orifice` 模式：孔口出流 `Q = Cd·A·√(2g·H)`，与节点水头耦合迭代

漏点接入后该点压力应下降。

## 运行

```bash
go run . solve example/loop4.json          # CLI 求解
go run . -http :8080                        # 启动 Web + API
```

API：

- `POST /api/solve`：拓扑 + 水箱水头 + 需水量 + 迭代上限/容差 → 各管流量、各点水头、迭代次数、是否收敛
- `GET /api/meta`：示例列表与摩阻模型信息
- `GET /example/<name>`：示例网

## 包结构

- `internal/network`：节点/管段/解析/校验/拓扑/找环
- `internal/hydraulics`：摩阻、全局梯度法(GGA)、Hardy-Cross、残差、能量闭合
- `internal/leak`：漏损试算
- `internal/api`：HTTP 服务
- `example/loop4.json`：含环小网，质量守恒可手算抽查

## 手算抽查（loop4）

对称环网：源 R(100m) 供 A、B(各 0.01)、C(0.02)。由对称性 A-B 流量为 0，R-A=R-B=0.02，A-C=B-C=0.01，C 水头约 99.52m。求解结果与之一致。

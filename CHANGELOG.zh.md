# Changelog

## v0.3.1
- venus-sector-manager：
  - 支持用于调节 PoSt 环节消息发送策略的 `MaxPartitionsPerPoStMessage` 和 `MaxPartitionsPerRecoveryMessage` 配置项

## v0.3.0

- venus-sector-manager：
  - 适配和支持 nv16
  - 对于一些特定类型的异常，返回特殊结果，方便 venus-worker 处理：
    - c2 消息上链成功，但扇区未上链 [#88](https://github.com/ipfs-force-community/venus-cluster/issues/88)
    - 对于能够确定感知到 ticket expired 的场景，直接终止当前扇区 [#143](https://github.com/ipfs-force-community/venus-cluster/issues/143)
    - 对于通过 `venus-sector-manager` 终止的扇区，或因其他原因确实扇区状态信息的情况，直接终止当前扇区 [#89](https://github.com/ipfs-force-community/venus-cluster/issues/89)
  - 升级 `go-jsonrpc` 依赖，使之可以支持部分网络异常下的重连 [#97](https://github.com/ipfs-force-community/venus-cluster/issues/97)
  - 支持新的可配置策略：
    - 各阶段是否随消息发送 funding [#122](https://github.com/ipfs-force-community/venus-cluster/issues/122)
    - SnapUp 提交重试次数 [#123](https://github.com/ipfs-force-community/venus-cluster/issues/123)
  - 支持可配置的 `WindowPoSt Chanllenge Confidential` [#163](https://github.com/ipfs-force-community/venus-cluster/issues/163)
  - 迁入更多管理命令
  - 配置调整：
    - 新增 [Miners.Commitment.Terminate] 配置块
    - 新增 [Miners.SnapUp.Retry] 配置块
    - 新增 [Miners.SnapUp] 中的 SendFund 配置项
    - 新增 [Miners.Commitment.Pre] 中的 SendFund 配置项
    - 新增 [Miners.Commitment.Prove] 中的 SendFund 配置项
    - 新增 [Miners.PoSt] 中的 ChallengeConfidence  配置项
- venus-worker：
  - 适配 venus-market 对于 oss piece store 的支持
  - 支持指定阶段批次启动 [#144](https://github.com/ipfs-force-community/venus-cluster/issues/144)
  - 支持外部处理器根据权重分配任务 [#145](https://github.com/ipfs-force-community/venus-cluster/issues/145)
  - 支持新的订单填充逻辑：
    - 禁止cc扇区 [#161](https://github.com/ipfs-force-community/venus-cluster/issues/161)
    - `min_used_space` [#183](https://github.com/ipfs-force-community/venus-cluster/pull/183)
  - 日志输出当前时区时间 [#87](https://github.com/ipfs-force-community/venus-cluster/issues/87)
  - 配置调整：
    - 废弃 [processors.limit] 配置块，替换为 [processors.limitation.concurrent] 配置块
    - 新增 [processors.limitation.staggered] 配置块
    - 新增 [[processors.{stage name}]] 中的 weight 配置项
    - 新增 [sealing] 中的 min_deal_space 配置项
    - 新增 [sealing] 中的 disable_cc 配置项
- 工具链：
  - 支持 cuda 版本编译
- 文档：
  - 更多 QA 问答
  - [10.venus-worker任务管理](./docs/zh/10.venus-worker任务管理.md)
  - [11.任务状态流转.md](./docs/zh/11.任务状态流转.md)
- 其他改善和修复



## v0.2.0
- 支持 snapup 批量生产模式
  - venus-worker 支持配置 `snapup` 类型任务
  - venus-sector-manager 支持配置 `snapup` 类型任务
  - venus-sector-manager 新增 `snapup` 相关的命令行工具：
    - `util sealer snap fetch` 用于按 `deadline` 将可用于升级的候选扇区添加到本地
	- `util sealer snap candidates` 用于按 `deadline` 展示可用于升级的本地候选扇区数量
  - 参考文档：[08.snapdeal的支持](https://github.com/ipfs-force-community/venus-cluster/blob/9be393761645f5fbd3a415b5ff1f50ec9254943c/docs/zh/08.snapdeal%E7%9A%84%E6%94%AF%E6%8C%81.md)

- 增强 venus-sector-manager 管理 venus-worker 实例的能力：
  - 新增 venus-worker 定期向 venus-secotr-manager 上报一些统计数据的机制
  - 新增 venus-sector-manager 的 `util worker` 工具集

- 增强 venus-sector-manager 根据功能拆分实例的能力：
  - 新增数据代理模式
  - venus-sector-manager 的 `util daemon run` 新增 `--conf-dir` 参数，可以指定配置目录
  - 新增外部证明计算器 (external prover) 的支持
  - 参考文档：[09.独立运行的poster节点](https://github.com/ipfs-force-community/venus-cluster/blob/9be393761645f5fbd3a415b5ff1f50ec9254943c/docs/zh/09.%E7%8B%AC%E7%AB%8B%E8%BF%90%E8%A1%8C%E7%9A%84poster%E8%8A%82%E7%82%B9.md)

- 修复 PreCommit/Prove 的 Batch Commit 未使用相应的费用配置的问题

- 其他调整
  - venus-worker 的配置调整
    - 新增 [[sealing_thread.plan]] 项
	- 新增 [attched_selection] 块，提供 `enable_space_weighted` 项，用于启用以剩余空间为权重选择持久化存储的策略，默认不启用
  - venus-secotr-manager 的配置调整
    - 废弃原 [Miners.Deal] 块，调整为 [Miners.Sector.EnableDeals] 项
    - 新增扇区生命周期项 [Miners.Sector.LifetimeDays]
    - 新增 [Miners.SnapUp] 块
	- 新增 [Miners.Sector.Verbose] 项，用于控制封装模块中的部分日志详尽程度
  - venus-sector-manager 的 `util storage attach` 现在默认同时检查 `sealed_file` 与 `cache_dir` 中目标文件的存在性
  - 其他改善和修复

## v0.1.2
- 一些为 SnapUp 支持提供准备的设计和实现变更
- `util storage attach` 新增 `--allow-splitted` 参数， 支持 `sealed_file` 与 `cache_dir` 不处于同一个持久化存储实例中的场景
  参考文档 [06.导入已存在的扇区数据.md#sealed_file-与-cache_dir-分离](https://github.com/ipfs-force-community/venus-cluster/blob/main/docs/zh/06.%E5%AF%BC%E5%85%A5%E5%B7%B2%E5%AD%98%E5%9C%A8%E7%9A%84%E6%89%87%E5%8C%BA%E6%95%B0%E6%8D%AE.md#sealed_file-%E4%B8%8E-cache_dir-%E5%88%86%E7%A6%BB)
- 开始整理 `Q&A` 文档
- 添加针对本项目内的组件的统一版本升级工具

## v0.1.1
- 外部处理器支持更灵活的工作模式，包含 **多任务并发**、**自定义锁**，使用方式参考：
  - [`processors.{stage_name} concurrent`](https://github.com/ipfs-force-community/venus-cluster/blob/main/docs/zh/03.venus-worker%E7%9A%84%E9%85%8D%E7%BD%AE%E8%A7%A3%E6%9E%90.md#processorsstage_name)
  - [`processors.ext_locks`](https://github.com/ipfs-force-community/venus-cluster/blob/main/docs/zh/03.venus-worker%E7%9A%84%E9%85%8D%E7%BD%AE%E8%A7%A3%E6%9E%90.md#processorsext_locks)
  - [07.venus-worker外部执行器的配置范例](https://github.com/ipfs-force-community/venus-cluster/blob/main/docs/zh/07.venus-worker%E5%A4%96%E9%83%A8%E6%89%A7%E8%A1%8C%E5%99%A8%E7%9A%84%E9%85%8D%E7%BD%AE%E8%8C%83%E4%BE%8B.md)

- 调整封装过程中获取 `ticket` 的时机，避免出现 `ticket` 过期的情况
- 简单的根据剩余空间量选择持久化存储的机制
- 跟进 `venus-market` 的更新
- 调整和统一日志输出格式
- 其他改善和修复

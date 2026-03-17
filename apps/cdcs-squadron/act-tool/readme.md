# act-tool 设计与实现计划

本项目基于 `gobox/tools/act` 包实现一个命令行工具，用于对目标目录中的 ACT XML 文档进行批量备份、重发、修改和回滚操作，并支持基于条件（患者 ID / 时间）列出文件概览。

---

## ✅ 一、需求概述

工具需支持以下功能：

1. **备份（Archive）**：将指定时间之前的文件归档到备份目录，按年月子目录组织。
2. **重发（Resend）**：将指定文件从目标目录移动到重发目录，并生成新的 XML 文件（保留原内容）。
3. **修改（Modify）**：将指定文件中的 `PatientID` 修改为指定值并生成新的 XML 文件做重发。
4. **回滚（Rollback）**：根据某次操作的结果清单，将生成的新文件删除，并把原文件还原回目标目录。
5. **列出（List）**：根据筛选条件（PatientID / 时间区间 / 时间点）列出符合条件的文件概览。

---

## 🧩 二、配置与运行方式

### 1) 配置文件（config.yaml 或 config.json）
将“目标路径”和“备份路径”配置在文件中，例如：

```yaml
targetPath: "D:/data/act/target"
backupPath: "D:/data/act/backup"
```

### 2) 命令行参数（可覆盖配置文件）
- `--config <path>`：指定配置文件路径
- `--target <path>`：覆盖目标路径
- `--backup <path>`：覆盖备份路径
- `--dry-run`：仅显示将要处理的文件和操作，不执行实际移动/写入
- `--interactive`：在执行实际操作前，展示列表并要求用户确认（适用于非交互式命令模式）
- `--mode <interactive|batch>`：选择运行模式，默认 `interactive`。

---
## 🧭 运行模式（模式选择）

程序支持两种运行方式：

- **交互式模式（默认）**：启动后显示菜单，用户通过输入数字选择操作，并根据提示输入参数。
- **批处理模式（batch）**：通过命令行参数直接指定要执行的操作（例如 `--cmd archive --before 2026-01-01`），适用于脚本自动化。

交互式模式也是批处理模式的超集：在交互式模式下仍可支持 `--dry-run` 和 `--interactive` 来增强安全性。

### 运行流程：交互式菜单（示例）

交互式模式启动后，程序显示主菜单：

```
act-tool (interactive mode)
1) 备份 (archive)
2) 列出 (list)
3) 重发 (resend)
4) 修改 (modify)
5) 回滚 (rollback)
0) 退出
请选择操作：
```

用户输入菜单编号后，程序依次提示所需参数（例如时间范围、PatientID、新 PatientID 等），并展示预览结果：

1. 计算待处理的文件列表（与 `list` 命令相同）。
2. 显示摘要：文件数量、时间范围、示例文件名等。
3. 若启动 `--interactive`，则要求输入 `Y` 或 `N`：
   - 输入 `Y`：执行操作（或根据 `--dry-run` 只输出预览）。
   - 输入 `N`：取消操作并返回主菜单。
4. 操作结束后返回主菜单，或按照用户选择退出。

交互式菜单还应支持：
- 输入 `b` 或 `back` 取消当前输入并返回上一级。
- 输入 `q` 或 `exit` 立即退出程序。

---

## 🚀 三、功能设计（命令与流程）

### 3.1 Dry Run 模式（dry-run）
- 所有变更型命令（archive/resend/modify/rollback）在 `--dry-run` 下仅输出将要执行的操作列表（源文件、目标路径、新文件名等），不实际移动或写入文件。
- 适用于审阅操作范围、验证筛选条件是否正确。

### 3.2 交互式模式（interactive）
- 当启用 `--interactive` 时，命令在执行前先列出将处理的文件，并提示用户输入确认（Y/n）。
- 交互流程：
  1. 计算出待处理文件列表
  2. 打印摘要（文件数量、示例文件名、操作类型）
  3. 等待用户输入，确认继续或取消

### 3.3 备份（archive）
- 根据 `--before` 或 `--until` 参数指定时间点，将目标目录中时间早于该时间的 XML 文件移动到备份目录。
- 备份目录结构: `<BackupPath>/<YYYY>/<MM>/...`
- 输出：处理文件数量、每月统计。
- 若启用 `--dry-run`：仅输出将归档的文件列表与对应目的路径。
- 若启用 `--interactive`：在执行归档前提示用户确认。

### 3.4 列出（list）
- 支持按 `PatientID`、按时间（等于/之前/之后/区间）过滤。
- 列出文件名、患者 ID、测试时间。
- 可与 `--dry-run` 组合使用以验证筛选条件（等价于 list 命令本身）。

### 3.5 重发（resend）
- 根据过滤条件筛选文件（同 list），将它们移动到 `backup/resend/<n>/` 目录作为操作快照。
- 生成新文件（同内容）写回目标目录，文件名加前缀 `R<timestamp>_<oldtesttime>.xml`。
- 保存操作结果：`result.txt` 记录 `{newFilename}|{originalFilename}`。
- `--dry-run` 模式下仅输出从哪个文件到哪个目标文件的映射关系，不写文件。
- `--interactive` 模式下在写文件前要求确认。

### 3.6 修改（modify）
- 同重发流程，但在生成新文件时将 `PatientID` 替换为命令参数指定的新值。
- 结果文件仍为 `result.txt`。
- `--dry-run` 模式下仅输出替换前后的 PatientID 与目标新文件路径。
- `--interactive` 模式下在写文件前要求确认。

### 3.7 回滚（rollback）
- 通过指定操作类型（resend/modify）和操作索引（目录编号），读取对应 `result.txt`。
- 删除生成的新文件（目标目录），将原文件从快照目录恢复回目标目录。
- `--dry-run` 模式下仅输出将要删除/恢复的文件对，不执行实际删除/移动。
- `--interactive` 模式下在执行回滚前要求确认。

---

## 🛠 四、实现细节（设计要点）

- 使用 `gobox/tools/act` 包中的核心接口（如 `ActManager`）实现核心操作。
- `ActManager` 负责：读取目录、解析 XML 中的 `PatientID` 与 `TestDateTime`、筛选文件、移动/复制文件、生成结果文件。
- 对 IO 操作放在可复用函数（`moveFile`、`readFile`、`writeFile`）中，保证跨平台兼容性。
- 增加 `dryRun` 与 `interactive` 选项层，用于决定是否执行文件系统操作与是否询问用户确认。
- 操作目录命名规则：
  - 备份（archive）直接按年月归档
  - 重发/修改：`<BackupPath>/resend/<seq>/`、`<BackupPath>/modify/<seq>/`
- 并发扫描：文件筛选可以使用 goroutine 并发读取，加速大目录处理。

---

## 📌 五、代码结构设计

为了让程序清晰、易维护，并方便后续扩展，建议将代码分层并按职责拆分成以下模块：

### 5.1 目录与文件结构（示例）

```
apps/cdcs-squadron/act-tool/
  ├── cmd/
  │   └── act-tool/
  │       └── main.go            # 程序入口，负责 flag 解析与启动模式选择
  ├── internal/
  │   ├── app/
  │   │   ├── app.go             # 运行入口（Interactive / Batch）
  │   │   ├── menu.go            # 交互式菜单与用户输入逻辑
  │   │   ├── commands.go        # 各功能命令的调度与执行
  │   │   └── confirm.go         # 交互确认与提示辅助函数
  │   ├── config/
  │   │   ├── config.go          # 配置结构与加载（yaml/json）
  │   │   └── defaults.go        # 默认值与验证逻辑
  │   └── act/
  │       └── manager.go         # 基于 gobox/tools/act 的文件操作逻辑
  └── go.mod
```

> 注：目录命名可根据团队习惯调整，关键是将“CLI 入口/交互”、“配置处理”和“业务逻辑（ActManager）”分离。

### 5.2 主要组件职责

#### 5.2.1 `cmd/act-tool/main.go`
- 解析命令行参数（`--mode`、`--dry-run`、`--interactive`、`--config`、`--cmd` 等）
- 加载配置
- 创建 `ActManager` 实例
- 根据 `--mode` 选择 `RunInteractive()` 或 `RunBatch()`

#### 5.2.2 `internal/app/app.go`
- 实现 `RunInteractive(ctx, config)`：
  - 显示主菜单（备份/列出/重发/修改/回滚/退出）
  - 通过 `menu.go` 读取用户选择与参数
  - 调用 `commands.go` 执行对应命令
- 实现 `RunBatch(ctx, config)`：
  - 根据 `--cmd`、`--target`、`--patient` 等参数直接调度命令
  - 支持 `--dry-run` 与 `--interactive`（命令前确认）

#### 5.2.3 `internal/app/menu.go`
- 负责交互式菜单展示与输入校验
- 提供通用输入函数（如：`PromptString`, `PromptInt`, `PromptTimeRange`）
- 支持快捷命令：`b/back` 返回、`q/exit` 退出

#### 5.2.4 `internal/app/commands.go`
- 定义每个操作的参数结构体（如 `ArchiveParams`、`ResendParams`）
- 将参数转换为 `act.Manager` 可用的过滤条件与行为开关
- 统一处理 `dryRun`/`interactive`：
  - 如果 `dryRun`：调用 `ActManager` 提供的预览方法（或在命令层打印预览）
  - 如果 `interactive`：在执行前调用确认函数

#### 5.2.5 `internal/act/manager.go`
- 封装和扩展 `gobox/tools/act` 的业务逻辑：
  - `FilterFiles(filter)` 返回文件列表
  - `Archive(before, dryRun)`
  - `Resend(params)`
  - `Modify(params)`
  - `Rollback(operation, index, dryRun)`
- 提供“预览”方法（仅列出将执行的文件/路径）
- 在内部统一处理文件移动/复制、结果写入、错误捕获

#### 5.2.6 `internal/config/config.go`
- 定义 `Config` 结构：`TargetPath`, `BackupPath`, `Mode`, `DryRun`, `Interactive` 等
- 支持从 YAML/JSON 加载，并支持环境变量覆盖（可选）

### 5.3 交互式/批处理的状态传递

- 通过 `context.Context` + 参数对象（如 `RunOptions{ DryRun bool, Interactive bool }`）在各层传递运行模式
- `ActManager` 只关注业务逻辑：当 `dryRun=true` 时仅返回“拟执行”列表，不进行写入

### 5.4 测试策略（与结构设计配套）

- 对 `internal/act` 进行单元测试：模拟目录与 XML 内容，验证筛选、备份、重发、修改、回滚行为
- 对 `internal/app` 进行集成测试：使用临时目录跑一轮交互式/批处理流程，验证流程控制与参数传递

---

## 📌 未来扩展（可选）

- 支持并发控制（限制 goroutine 数量，避免文件系统 I/O 饱和）。
- 将 `PatientID` / 时间筛选条件支持正则或模糊匹配。

- 支持并发控制（限制 goroutine 数量，避免文件系统 I/O 饱和）。
- 将 `PatientID` / 时间筛选条件支持正则或模糊匹配。

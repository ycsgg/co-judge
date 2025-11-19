
# co-judge

这是一个用于 MIPS 指令处理、汇编/反汇编、仿真与自动评测的小型工具集合。代码用 Go 语言实现：

- `mips2hex`：将 MIPS 汇编代码（.s）汇编为每行 32-bit 十六进制指令文件（无 `0x` 前缀）。
- `hex2mips`：把十六进制/二进制/十进制的机器码行反汇编为可读的汇编指令。
- `mipsim`：一个简单的逐条指令仿真器，会打印每步的寄存器写入与存储器写入信息
- `judger`：针对 Logisim 的评测封装（`-mode logisim`），用于用 Logisim 的仿真结果与本地仿真比较并输出差异。

以下说明基于仓库 `src` 目录下的源码（已包含 `go.work`，建议在 `src` 下执行构建/运行命令）。

## 目录结构

- `src/mips2hex`：汇编器，入口 `main.go`，内部包含 `parser`, `assembler`, `emitter`。
- `src/hex2mips`：反汇编器，入口 `main.go`，使用 `hex2mips/disassembler` 进行指令解码。
- `src/mipsim`：仿真器，入口 `main.go`，核心位于 `cpu` 包（`cpu.Run` 将指令装入内存并逐步执行）。
- `src/judger`：评测器（Logisim 集成），入口 `main.go`。


## 构建与运行

建议在 `src` 目录下运行下面的命令（仓库根目录有 `src/go.work`）：

# 进入源码目录
```
cd .\src
```
# 1) 直接运行汇编器

```
go run .\mips2hex -input .\mips2hex\test\code0.s -output .\out_instr.txt -base 0x3000
```

# 上述命令说明：
- `-input`：输入 MIPS 汇编源文件（支持 .text/.word、标签、li、常见指令等）
- `-output`：输出每行 8 字节（32-bit）大写十六进制字符串（无 0x 前缀），由 `emitter.WriteHexLines` 生成
- `-base`：指定 `.text` 段基址（十进制或 `0x..`）。分支与跳转会按该基址进行 PC 与目标地址计算。

# 2) 反汇编单条或文件
使用 hex2mips 可以对一个 hex 行或文件进行反汇编。
```
go run .\hex2mips -input .\out_instr.txt
```
或者对单个 hex 字符串：
```
go run .\hex2mips -ih 0x012a4020
```
# 3) 运行仿真器（mipsim）
默认最大执行步数 10000（防止死循环），可用 `-limit` 覆盖。
```
go run .\mipsim -f .\out_instr.txt
```
限制步数示例：
```
go run .\mipsim -f .\out_instr.txt -limit 500
```
# 仿真器行为：
 - 从 PC 基址 0x3000 开始将指令装入内存（见 cpu.New 初始化）
 - 按顺序执行，每步打印反汇编文本、PC、寄存器写入与内存写入信息

# 4) 评测（Logisim）
# 用法：
```
go run .\judger -mode logisim <logisim_jar> <circuit.circ> <hex_path> <output_path>
```
其中 <hex_path> 为被评测的 hex 文件，评测结果会写入 <output_path>（不一致时会在输出目录写 detail.log）

## 主要实现细节与约定

- 汇编器：`mips2hex` 使用两遍解析（`parser.ParseLines`）：第一遍收集标签地址，第二遍生成指令/数据项。`.text` 段起始地址被假定为 0（汇编时以偏移计数），指令大小通常为 4 字节，`li` 会被扩展为两条指令（8 字节）。
- 指令编码在 `mips2hex/assembler` 中实现，支持常见的 R/I/J 类型指令、移位、分支等。特殊伪指令 `li` 与 `nop` 被单独处理。
- 输出格式：`emitter.WriteHexLines` 会将每个 uint32 按字节写为 8 位十六进制小写字符串（每行一条指令）。
- 反汇编：`hex2mips` 可接受二进制串（32 位）、十六进制（包含/不包含 `0x` 前缀）或十进制，且支持从 stdin、文件或单个输入字符串反汇编。
- 仿真器：`mipsim` 的 `cpu` 从 PC 基址 0x3000 装载指令并执行（参见 `cpu.New()`），打印每步状态，适用于步进观察指令效果。`cpu` 中包含 `signExtend16`、通用寄存器数组、Hi/Lo 寄存器和内存映射（map）。


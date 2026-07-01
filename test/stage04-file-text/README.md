# Stage 04 测试：文件与文本处理（日志采集雏形）

这一关直接对着你的目标项目 **日志采集工具** 去。  
考察标准库 `os`、`bufio`、`strings`、`strconv`，以及你前面学的 `map`、错误处理。

本目录已经给你准备好一份日志文件 **`access.log`**（20 行，模拟真实应用日志）。  
在本目录新建 `main.go`，写一个**日志分析工具**。

## 背景场景（运维向）

你天天 `grep ERROR app.log | wc -l`、`awk '{print $3}' | sort | uniq -c` 分析日志。  
现在用 Go 把这套"读文件 → 逐行解析 → 统计 → 输出"的活儿写成程序。  
这就是日志采集/分析工具的内核。

---

## 📚 知识铺垫（先读这个！）

### 读文件的标准姿势

```go
file, err := os.Open("access.log")   // 打开文件（只读）
if err != nil {
    fmt.Println("打开文件失败:", err)
    return
}
defer file.Close()                   // 别忘了关闭（类比 trap EXIT）

scanner := bufio.NewScanner(file)    // 按行读取的利器
for scanner.Scan() {                 // 每次读一行，读到 EOF 返回 false
    line := scanner.Text()           // 拿到这一行的字符串（不含换行符）
    // 处理 line ...
}
```

> **类比运维**：`bufio.Scanner` 就是 Go 版的 `while read line; do ... done < file`。  
> `os.Open` 只读，类比 `cat`；后面任务要写文件才用 `os.Create`。

### 日志每行的结构

```
2026-06-30 10:23:05 ERROR db connection timeout host=10.0.0.21
└──日期──┘ └─时间─┘ └级别┘ └────────── 消息 ──────────┘
  字段0     字段1   字段2         字段3...
```

用 `strings.Fields(line)` 按空白切分成 `[]string`，`fields[2]` 就是日志级别。  
（`strings.Fields` 类比 `awk` 的默认按空格分列，比 `strings.Split` 更省心，多个空格也能正确处理。）

### 常用字符串工具

| 函数 | 作用 | 类比 |
|------|------|------|
| `strings.Fields(s)` | 按空白切成切片 | `awk '{print $1,$2}'` |
| `strings.Split(s, "=")` | 按指定分隔符切 | `cut -d= -f2` |
| `strings.Contains(s, "x")` | 是否包含 | `grep x` |
| `strings.HasPrefix(s, "x")` | 是否以...开头 | `grep '^x'` |
| `strconv.Atoi(s)` | 字符串转 int | 数字比较前必做 |

---

## 需求（分 4 个任务，逐步加功能）

### 任务 1：读文件 + 数行数

打开 `access.log`，逐行读取，统计**总行数**，打印出来。

**要求**：
- 用 `os.Open` + `bufio.Scanner`
- `defer file.Close()`
- 处理打开失败的错误（文件不存在时别 panic，优雅报错）
- 输出类似：`总日志行数: 20`

---

### 任务 2：按日志级别统计（核心）

遍历每一行，取出**日志级别**（`INFO` / `WARN` / `ERROR` / `DEBUG`），  
用 `map[string]int` 统计每种级别出现了多少次。

**要求**：
- 用 `strings.Fields(line)` 切分，取第 3 个字段（下标 2）作为级别
- 用 map 计数
- 输出类似（顺序不强求）：
  ```
  级别统计:
    INFO  : 8
    WARN  : 4
    ERROR : 6
    DEBUG : 2
  ```
- **健壮性**：万一某行字段不够（空行、格式乱），不能数组越界 panic。
  提示：先判断 `len(fields) >= 3` 再取 `fields[2]`。

---

### 任务 3：过滤 + 提取错误详情（grep 的活儿）

只挑出 `ERROR` 级别的行，打印出来，并统计 ERROR 总数。

**要求**：
- 遍历时判断级别是不是 `ERROR`，是才处理
- 打印每条 ERROR 的**时间 + 消息**（消息 = 第 4 个字段往后的内容）
  提示：`strings.Join(fields[3:], " ")` 把消息部分拼回来
- 最后输出：`共发现 X 条 ERROR`
- 输出类似：
  ```
  === ERROR 日志 ===
  [10:23:05] db connection timeout host=10.0.0.21
  [10:23:09] panic recovered handler=UserHandler err=nil_pointer
  ...
  共发现 6 条 ERROR
  ```

---

### 任务 4：把统计结果写入文件（采集工具的输出）

把任务 2 的级别统计结果，**写入一个新文件** `report.txt`。

**要求**：
- 用 `os.Create("report.txt")` 创建文件
- `defer f.Close()`
- 用 `fmt.Fprintf(f, ...)` 往文件里写（注意是 `Fprintf`，第一个参数是文件）
- 写完后在终端提示：`报告已生成: report.txt`
- report.txt 内容类似：
  ```
  日志分析报告
  生成时间: 2026-06-30
  ----------------
  INFO  : 8
  WARN  : 4
  ERROR : 6
  DEBUG : 2
  总计  : 20 行
  ```

---

## 🎯 最终要求

`main.go` 实现上面 4 个任务（可以分成 `task1()`~`task4()` 依次调用，也可以合并成一次遍历完成全部统计——后者更高效，是加分项）。

**必须能跑通并输出：**
- 总行数
- 各级别计数
- ERROR 明细
- 生成 report.txt 文件

**通过标准：**
- 文件正确打开并 `defer Close()`
- 打开失败时优雅报错，不 panic
- 字段不足时不越界 panic
- map 统计正确
- report.txt 成功生成且内容正确

---

## 💡 提示（卡住了看这里）

### 提示 1：一次遍历搞定所有统计（推荐）
不用读 4 遍文件，读一遍，在循环里同时干所有事：
```go
levelCount := make(map[string]int)
var errorLines []string
total := 0

for scanner.Scan() {
    line := scanner.Text()
    total++
    fields := strings.Fields(line)
    if len(fields) < 3 {
        continue   // 跳过格式不对的行
    }
    level := fields[2]
    levelCount[level]++
    if level == "ERROR" {
        errorLines = append(errorLines, line)
    }
}
```

### 提示 2：Scanner 读完要检查错误
```go
if err := scanner.Err(); err != nil {
    fmt.Println("读取出错:", err)
}
```

### 提示 3：写文件的两种写法
```go
f, _ := os.Create("report.txt")
defer f.Close()
fmt.Fprintf(f, "INFO: %d\n", levelCount["INFO"])   // 格式化写
// 或
f.WriteString("一行文本\n")                          // 直接写字符串
```

---

写完对我说「**检查**」。这一关比并发轻松，而且写完你就有一个能用的日志分析小工具了，很有成就感。加油！💪

# Stage 04 总结：文件与文本处理（日志采集雏形）

## 这一关做了什么
写了一个**日志分析工具**：读日志文件，逐行解析，按级别统计，抓 ERROR 明细，生成 report.txt。这是目标项目"日志采集工具"的内核。

## 核心知识点

### 1. 读文件标准姿势
```go
file, err := os.Open("access.log")   // 只读打开
if err != nil {
    fmt.Println("打开失败:", err)
    return
}
defer file.Close()                   // 别忘关闭

scanner := bufio.NewScanner(file)
for scanner.Scan() {                 // 逐行读，EOF 返回 false
    line := scanner.Text()           // 当前行（不含换行符）
}
if err := scanner.Err(); err != nil {// 循环后检查真错误
    ...
}
```
> 类比：`while read line; do ...; done < file`

### 2. 字符串处理工具箱
| 函数 | 作用 | 类比 |
|------|------|------|
| `strings.Fields(s)` | 按空白切片 | `awk '{print $1}'` |
| `strings.Split(s,"=")` | 按分隔符切 | `cut -d=` |
| `strings.Join(sl," ")` | 切片拼回字符串 | 反向操作 |
| `strings.Contains` | 是否包含 | `grep` |
| `strconv.Atoi` | 字符串转数字 | 数字比较前必做 |

### 3. 写文件
```go
f, err := os.Create("report.txt")    // 创建（会覆盖）
defer f.Close()
fmt.Fprintf(f, "INFO: %d\n", n)       // 格式化写入（第一个参数是文件）
f.WriteString("一行\n")               // 直接写字符串
```
> 记住：`Fprintf` 的 `F` = File，第一个参数是写入目标。

### 4. time 格式化（反直觉但重要）
```go
time.Now().Format("2006-01-02 15:04:05")
```
`2006-01-02 15:04:05` 是**记忆口诀**（不是 Go 生日）：
```
01/02 03:04:05PM '06 -0700
月  日  时 分 秒    年   时区
1   2   3  4  5     6    7
```
想要什么格式，就把这个基准时间摆成什么样。

## 关键领悟

### 领悟 1：Scanner 是"一次性流"，读完到底不倒带
```go
for scanner.Scan() { ... }   // 第一遍读到 EOF
for scanner.Scan() { ... }   // 第二遍：一次都进不去！文件已读完
```
**这不是编译错误，是运行逻辑错误**——第二个循环啥也读不到。
> 类比：`cat file | while read` 读过的数据像水流走了，回不来。
> 想重读：`file.Seek(0,0)` 倒带，或重新 `os.Open`，或——**一次遍历干完所有事**（首选）。

### 领悟 2：map 遍历顺序是随机的
```go
for level, count := range levelCount { ... }  // 每次顺序可能不同！
```
给用户看的输出要顺序稳定 → 用一个固定切片控制顺序：
```go
standardLevels := []string{"INFO","WARN","ERROR","DEBUG"}
for _, lvl := range standardLevels {
    fmt.Printf("%-6s: %d\n", lvl, levelCount[lvl])
}
```

### 领悟 3：一次遍历，多项统计（真实工具的写法）
真实日志可能几个 G，读多遍是灾难。一个循环里同时数行数、统计级别、收集 ERROR：
```go
for scanner.Scan() {
    lineCount++
    fields := strings.Fields(scanner.Text())
    if len(fields) < 3 { continue }   // 防越界
    level := fields[2]
    levelCount[level]++
    if level == "ERROR" {
        errorDetails = append(errorDetails, ...)
    }
}
```

## 一句话记忆
- 读文件三件套：`os.Open` + `bufio.Scanner` + `defer Close()`
- Scanner **单向流**，想多用一次遍历全算完
- map 顺序随机，输出要顺序就用**固定切片**引导
- 取字段前先 `len(fields) >= N` **防越界**

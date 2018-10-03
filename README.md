
# CLI 命令行实用程序开发基础

## Golang开发selpg

### 应用概述
该实用程序从标准输入或从作为命令行参数给出的文件名读取文本输入。它允许用户指定来自该输入并随后将被输出的页面范围。例如，如果输入含有 100 页，则用户可指定只打印第 35 至 65 页。这种特性有实际价值，因为在打印机上打印选定的页面避免了浪费纸张。另一个示例是，原始文件很大而且以前已打印过，但某些页面由于打印机卡住或其它原因而没有被正确打印。在这样的情况下，则可用该工具来只打印需要打印的页面。

### 命令格式
selpg -s start_page -e end_page [-f | -l lines_per_page] [-d dest] [in_filename]

### 命令参数解释
-s : -s后接正数，表示读取标准输入或文件的开始页数

-e : -e后接整数，表示读取标准输入或文件的结尾页数，注意: start_page>0,end_page>0且end_page>= start_page

-f : 表示以\f换页符作为标志一页

-l : 后接正数，表示一页有多少行。注意：-f和-l参数指令不能同时出现，默认为-l模式，且默认一页有72行

-d : 将选定的页直接发送至打印机。若不存在-d参数，则默认为标准输出

in_filename : 输入文件。不存在时，默认为标准输入

### 函数分析
Init : 初始化flag

usage : 打印selpg的用法

process_args : 处理命令行参数，利用flag或pflag对参数进行操作

process_input ：根据参数进行文本的输入与输出

outputWithModeL : 以-l的模式进行输出

outputWithModeP : 以-f的模式进行输出

main : 调用Init(),process_args(),process_input()函数


### 应用使用与测试

1、selpg -s start_page -e end_page
```
selpg -s 1 -e 1 
```


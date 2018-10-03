
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

从标准输入选定特定页写至标准输出（即屏幕）

```
selpg -s 1 -e 1 
```
![image](https://github.com/Howlyao/selpg/blob/master/image/1.png)

2、selpg -s start_page -e end_page input_file

命令将把“input_file”特定页写至标准输出（也就是屏幕）

```
selpg -s 1 -e 1 test.md
```
![image](https://github.com/Howlyao/selpg/blob/master/image/2.png)
![image](https://github.com/Howlyao/selpg/blob/master/image/3.png)

3、selpg -s start_page -e end_page < input_file

selpg 读取标准输入，而标准输入已被 shell／内核重定向为来自“input_file”而不是显式命名的文件名参数。输入的特定页被写至屏幕

```
selpg -s 1 -e 1 < test.md
```
![image](https://github.com/Howlyao/selpg/blob/master/image/4.png)

4、command | selpg -s start_page -e end_page

“other_command”的标准输出被 shell／内核重定向至 selpg 的标准输入。将特定页写至 selpg 的标准输出（屏幕）

```
cat test.md | selpg -s 1 -e 1
```
![image](https://github.com/Howlyao/selpg/blob/master/image/5.png)

5、selpg -s start_page -e end_page input_file > output_file

selpg 将特定页写至标准输出；标准输出被 shell／内核重定向至“output_file”
```
selpg -s 1 -e 1 test.md > output_file
```
![image](https://github.com/Howlyao/selpg/blob/master/image/6.png)

6、selpg -s start_page -e end_page input_file 2>error_file

selpg 将特定页写至标准输出（屏幕）；所有的错误消息被 shell／内核重定向至“error_file”
```
selpg -s 1 -e 1 test.md 2>error_file
```
![image](https://github.com/Howlyao/selpg/blob/master/image/7.png)

7、selpg -s start_page -e end_page input_file | command

selpg 的标准输出透明地被 shell／内核重定向，成为“other_command”的标准输入，特定页被写至该标准输入

```
selpg -s 1 -e 1 test.md | selpg -s 2 -e 2 -l 2
```
![image](https://github.com/Howlyao/selpg/blob/master/image/8.png)

8、selpg -s start_page -e end_page -l lines_per_page input_file

该命令将页长设置为lines_page_page行，这样 selpg 就可以把输入当作被定界为该长度的页那样处理。特定页被写至 selpg 的标准输出（屏幕）。
```
selpg -s 1 -e 1 -l 1 test.md 
```
![image](https://github.com/Howlyao/selpg/blob/master/image/9.png)

### 参考文献

[开发 Linux 命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)

[Linux命令行程序设计](https://wenku.baidu.com/view/c7cf91ee5ef7ba0d4a733b58.html)

[Golang之使用Flag和Pflag](https://o-my-chenjian.com/2017/09/20/Using-Flag-And-Pflag-With-Golang/)

[Command](https://godoc.org/os/exec#example-Command)

[Pipe](https://godoc.org/io#Pipe)

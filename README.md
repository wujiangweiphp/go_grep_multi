# go_grep_multi

## 使用
<pre>
NAME:
scan - A new cli application

USAGE:
scan [global options] command [command options]

COMMANDS:
help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
--path value, -p value                                   需要扫描的目录路径
--lineCount, -l                                          是否需要统计扫描行数 (default: false)
--exclude value, -e value                                需要跳过的路径标识
--include value, -i value                                需要包含的路径或文件标识
--regexp, -r                                             是否使用正则匹配 (default: 默认使用字符串)
--content value, -c value [ --content value, -c value ]  需要匹配的内容, 多个代表同一个文件里这几个都出现过
--onlyFile, -o                                           只扫描文件名 (default: false)
--ignoreCase, -u                                         只扫描文件名 (default: false)
--help, -h                                               show help
</pre>



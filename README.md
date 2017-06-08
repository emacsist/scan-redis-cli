# 说明

这是一个用于 redis scan 的小工具

即按条件来scan并打印出它们的 key 的名字

# 用法

直接执行该命令，它会自动打印出用法了，如:

```bash
[15:04:06] emacsist:src git:(master*) $ ./main
Usage of ./main:
  -a string
    	-a=密码，默认为空
  -c int
    	-c=最大key个数，默认为0，即不限
  -h string
    	-h=IP地址:端口，默认为 127.0.0.1:6379 (default "127.0.0.1:6379")
  -help
    	-help 显示该帮助 (default true)
  -l int
    	-l=固定key的长度，默认值为0，即不限
  -p string
    	-p=匹配符，默认为* (default "*")
```

# 下载

如果不想自己编译，可以下载好我已经编译的版本. 仓库里的 

scan-redis-cli-amd64 => Linux 64 位

scan-redis-cli-macos64 => MacOS 64位

其他版本，请自行编译使用。
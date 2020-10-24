# scanPort

默认扫描conf/ip.txt文件中的ip或ip段，将结果转化为json写入log中。

### 帮助信息
```
scanPort -h 
Options:
  -h    帮助信息
  -n int
        进程数 例如:-n=10 (default 100)
  -p string
        端口号范围 例如:-p=80,81,88-1000 (default "80")
  -t int
        超时时长(毫秒) 例如:-t=200 (default 200)

```

#### 例1：指定端口号扫描，使用1000个协程进行
```
scanport -p 80,81,88-3306 -n 1000 
```

# scanPort

### 更新记录
1. 修改为每个端口扫描两次，应对防火墙代应答1次的情况。
2. 增加支持扫描IP段、IP范围"-"。如 192.168.1.0/24，192.168.1.0-24
3. 修改为从文件读取目标进行扫描。默认扫描conf/ip.txt文件中的目标
4. 增加扫描前使用ping探测主机存活
5. 修改为将结果转化成json写入log文件。

### TODO

- [ ] 增加取消ping选项
- [ ] 增加识别web服务功能
- [ ] 增加不存在日志文件提醒（需要创建log/result.log，不然报错invalid argument）
- [ ] 增加随机IP和随机端口扫描探测功能
### Usage
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

#### 例1：扫描conf/ip.txt，指定端口号扫描，使用1000个协程进行
```
scanport -p 80,81,88-3306 -n 1000 
scanport -p 1-65535 -n 1000 
```

### 性能测试

1核1G服务器扫描一个ip的1-65535端口，随着参数变化，记录所用时间。

```
-n 进程数，默认为10
-t 超时时长(毫秒)，默认为200
```

参数 | 时间 | 备注
---|---|---
-n 100 | 268.24s | 扫描结果：80 443
-n 200 | 135.73s | 扫描结果：80 443
-n 200 -t 100 | 69.96s | 没扫出来端口
-n 100 -t 100 | 138.48s | 没扫出来端口
-n 200 | 135.54s | 测试是否被封，可以扫出来，证明是-t的问题
-n 200 -t 150 | 102.39s | 没扫出来端口
-n 250 -t 200 | 110.34s | 扫描结果：443

粗略的测试，大概最佳参数是-n 200 -t 200吧:)

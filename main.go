package main

import (
	"bufio"
	"flag"
	"fmt"
	"scanPort/lib"
	"scanPort/scan"
	"os"
	"time"
	"encoding/json"
	"strings"
)

var (
	startTime = time.Now()
	port      = flag.String("p", "80", "端口号范围 例如:-p=80,81,88-1000")
	ping 	  = flag.Bool("ping", true, "是否探测主机存活")
	timeout   = flag.Int("t", 200, "超时时长(毫秒) 例如:-t=200")
	process   = flag.Int("n", 100, "进程数 例如:-n=10")
	h         = flag.Bool("h", false, "帮助信息")

)

type IpInfo struct{
	IP string `json:"ipaddr"`
	Ports  []int `json:"ports"`
}

//go run main.go -h
func main() {
	flag.Parse()
	//帮助信息
	if *h == true {
		lib.Usage("scanPort version: scanPort/1.10.0\n Usage: scanPort [-h] [-ip ip地址] [-n 进程数] [-p 端口号范围] [-t 超时时长] [-path 日志保存路径]\n\nOptions:\n")
		return
	}

	fmt.Printf("========== Start %v  ==================== \n", time.Now().Format("2006-01-02 15:04:05"))

	//初始化
	scanIP:=scan.NewScanIp(*timeout,*process,true)
	fileName := "log/result.log"				//结果文件
	file, err := os.Open("conf/ip.txt")	//扫描ip文件
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	//遍历扫描文件ip.txt
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		cidr := sc.Text()
		ips := []string{}
		//处理cidr格式，如192.168.1.0/24
		if strings.Contains(cidr, "/"){
			hosts, _ := lib.Hosts(cidr)
			for _, ip := range hosts {
				ips = append(ips, ip)
			}
		}else if strings.Contains(cidr, "-"){
			//处理 192.168.1.0-24 格式
			ips, err = scanIP.GetAllIp(cidr)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		} else{
			//处理单个ip格式
			ips = append(ips, cidr) //单个ip
		}

		//遍历ip数组
		for i := 0; i < len(ips); i++ {

			//ping判断是否存活
			if *ping && !lib.Ping(ips[i]){
				//fmt.Println(ips[i]," 不存活")
				continue
			}
			fmt.Println(ips[i]," 存活")

			//获取扫描结果
			ports := scanIP.GetIpOpenPort(ips[i], *port)
			if len(ports) > 0 {
				f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					if err := f.Close(); err != nil {
						fmt.Println(err)
					}
					continue
				}
				//构造扫描结果结构体
				ipinfo := IpInfo{
					IP: ips[i],
					Ports: ports,
				}
				//结构体转json格式
				jsonBytes, err := json.Marshal(ipinfo)
				if err != nil {
					fmt.Println(err)
				}

				//写入结果文件
				var str = fmt.Sprintf("%v\n", string(jsonBytes))
				if _, err := f.WriteString(str); err != nil {
					if err := f.Close(); err != nil {
						fmt.Println(err)
					}
					continue
				}
			}
		}
	}
	fmt.Printf("========== End %v 总执行时长：%.2fs ================ \n", time.Now().Format("2006-01-02 15:04:05"), time.Since(startTime).Seconds())
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"scanPort/lib"
	"scanPort/scan"
	"strings"
	"time"
)

var (
	startTime = time.Now()
	help      = flag.Bool("h", false, "帮助信息")
	port      = flag.String("p", "80", "端口号范围 例如:-p=80,81,88-1000")
	nping 	  = flag.Bool("nping", false, "是否禁用ping探测主机存活")
	timeout   = flag.Int("t", 200, "超时时长(毫秒) 例如:-t=200")
	process   = flag.Int("n", 100, "进程数 例如:-n=100")
	format    = flag.String("f", "json", "输出格式支持json, txt。默认为json")
)

//go run main.go -h
func main() {
	flag.Parse()
	//帮助信息
	if *help == true {
		lib.Usage("scanPort version: scanPort/1.10.0\n Usage: scanPort [-h] [-p 端口号范围] [-n 进程数] [-t 超时时长] [-nping 是否禁用ping] [-f 输出格式]\n\nOptions:\n")
		return
	}

	fmt.Printf("========== Start %v  ==================== \n", time.Now().Format("2006-01-02 15:04:05"))

	//初始化
	scanIP:=scan.NewScanIp(*timeout,*process,true)
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
			if !*nping && !lib.Ping(ips[i]){
				//fmt.Println(ips[i]," 不存活")
				continue
			}
			fmt.Println(ips[i]," 存活")

			//获取扫描结果
			ports := scanIP.GetIpOpenPort(ips[i], *port)
			if *format == "json" {
				lib.WriteJsonResults(ips[i], ports)
			}else if(*format == "txt"){
				lib.WriteTxtResults(ips[i], ports)
			}else {
				lib.WriteJsonResults(ips[i], ports)
			}

		}
	}
	fmt.Printf("========== End %v 总执行时长：%.2fs ================ \n", time.Now().Format("2006-01-02 15:04:05"), time.Since(startTime).Seconds())
}

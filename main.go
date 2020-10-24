package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/xinhaiyucheng/scanPort/lib"
	"github.com/xinhaiyucheng/scanPort/scan"
	"os"
	"time"
	"encoding/json"
)

var (
	startTime = time.Now()
	//ip        = flag.String("ip", "127.0.0.1", "ip地址 例如:-ip=192.168.0.1-255 或直接输入域名 xs25.cn")
	port      = flag.String("p", "80", "端口号范围 例如:-p=80,81,88-1000")
	//path      = flag.String("path", "log", "日志地址 例如:-path=log")
	timeout   = flag.Int("t", 200, "超时时长(毫秒) 例如:-t=200")
	process   = flag.Int("n", 100, "进程数 例如:-n=10")
	h         = flag.Bool("h", false, "帮助信息")
	//ipfile      = flag.Bool("file", false, "ip文件")
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

	//创建目录
	lib.Mkdir(*path)

	//初始化
	scanIP:=scan.NewScanIp(*timeout,*process,true)

	file, err := os.Open("conf/ip.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		cidr := sc.Text()
		hosts, _ := lib.Hosts(cidr)
		ips := []string{}
		for _, ip := range hosts {
			ips = append(ips, ip) //得到ip list 切片
		}

		fileName := "log/result.log"
		for i := 0; i < len(ips); i++ {
			ports := scanIP.GetIpOpenPort(ips[i], *port)
			if len(ports) > 0 {
				f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					if err := f.Close(); err != nil {
						fmt.Println(err)
					}
					continue
				}
				ipinfo := IpInfo{
					IP: ips[i],
					Ports: ports,
				}
				fmt.Println("struct: ", ipinfo)
				jsonBytes, err := json.Marshal(ipinfo)
				if err != nil {
					fmt.Println(err)
				}
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

	/*
	ips, err := scanIP.GetAllIp(*ip)
	if err != nil {
		fmt.Println(err.Error())
		return
	}*/
	//log

	fmt.Printf("========== End %v 总执行时长：%.2fs ================ \n", time.Now().Format("2006-01-02 15:04:05"), time.Since(startTime).Seconds())

}

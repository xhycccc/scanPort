package lib

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	fileName = "result.log"
)

func Usage(str string) {
	fmt.Fprintf(os.Stderr, str)
	flag.PrintDefaults()
}

func Mkdir(path string){
	f, err := os.Stat(path)
	if err != nil || f.IsDir() == false {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			fmt.Println("创建目录失败！", err)
			return
		}
	}
}

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func WriteJsonResults(ip string, ports []int) (bool){

	type IpInfo struct{
		IP string `json:"ipaddr"`
		Ports  []int `json:"ports"`
	}

	if len(ports) > 0 {
		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
			return false
		}
		//构造扫描结果结构体
		ipinfo := IpInfo{
			IP: ip,
			Ports: ports,
		}
		//结构体转json格式
		jsonBytes, err := json.Marshal(ipinfo)
		if err != nil {
			fmt.Println(err)
			return false
		}

		//写入结果文件
		var str = fmt.Sprintf("%v\n", string(jsonBytes))
		if _, err := f.WriteString(str); err != nil {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
			return false
		}
	}
	return true
}

func WriteTxtResults(ip string, ports []int) (bool){

	if len(ports) <= 0 {
		return false
	}
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
		return false
	}
	for _,port := range ports{
		var str = fmt.Sprintf("%s	%d\n", ip, port)
		if _, err := f.WriteString(str); err != nil {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
			return false
		}
	}
	return true
}
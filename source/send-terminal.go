/**
 * anthor:liuchao@gomeplus.com
 */
package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"flag"
	"time"
)
const VIDEO_LOG = "/app/logs/applog/api/finish_log/"
const HELP_FILE = `    ***Send File Terminal Sever***
  -host string 
		IP address of receiving server,required
  -port string
		PORT address of receiving server,default value is 8888
  -path string
		The path to which you want to send a fileAbsolute path,
		default value is:/app/logs/applog/api/finish_log/2019/04/2019_04_03.txt
  -source string
		Send Host File Identification,eg: 1_1,1_2
  eg:
	./send-terminal -host 10.115.1.17 -port 8888 -path /tmp/send.txt -source 10.115.1.18_1_1
`
func main(){
	var h bool
	flag.BoolVar(&h, "h", false, "this help")
	var host = flag.String("host","", "IP address of receiving server,required")
	var port = flag.String("port", "8888","")
	var path = flag.String("path","","")
	var source = flag.String("source","","")
	flag.Usage = usage
	flag.Parse()
	if h{
		flag.Usage()
		os.Exit(1)
	}
	if *host==""{
		fmt.Println("IP address of receiving server,required")
		return
	}
	fmt.Println("receiving host：",*host)
	fmt.Println("receiving port：",*port)
	default_path := VIDEO_LOG +time.Now().Format("2006/01")+"/"+time.Now().Format("2006_01_02")+".txt"
	if *path == ""{
		*path = default_path
	}
	file_info,err := os.Stat(*path)
	if err != nil{
		fmt.Println("please check local path, ",err)
		return
	}
	conn,err1 := net.Dial("tcp",*host+":"+*port)
	if err1 != nil{
		fmt.Println("please check receiving server, ",err1)
		return
	}
	defer conn.Close()
	//send name
	send_file_name := "From_"+getLocalIp()+*source+"_"+file_info.Name()
	_,err2 := conn.Write([]byte(send_file_name))
	if err2!= nil{
		fmt.Println("send local filename expect, ",err2)
		return
	}
	var n int
	buf := make([]byte,1024)
	n,err = conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err = ", err)
		return
	}
	if "ok"== string(buf[:n]){
		SenderFile(*path,conn)
	}

}
//send file
func SenderFile(path string,conn net.Conn){
	var err error
	var n int
	//打开本地文件
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("send file open fail",err)
		return
	}
	defer f.Close()
	buf := make([]byte,1024)
	for{
		n,err = f.Read(buf)
		if err!=nil{
			if err == io.EOF{
				fmt.Println("success to file ",path," time ",time.Now().Format("2006/01/01 10:10:10"))
			}else{
				fmt.Println("Source file reading failed ",err)
			}
			return
		}
		conn.Write(buf[:n])
	}
}
func usage(){
	fmt.Fprintf(os.Stderr,HELP_FILE)
}
//get local host
func getLocalIp() string{
	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		return"localhost"
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

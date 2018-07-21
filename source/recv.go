package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"flag"
)

func main(){
	var ip = flag.String("h","", "please input ip:")
	var port = flag.String("p", "8888","please input port,default 8888")
	flag.Parse()
	if *ip==""{
		fmt.Println("please input ip,eg:-h 127.0.0.1")
		return
	}
	fmt.Println("接收端IP地址：",*ip)	
	fmt.Println("接收端PORT地址：",*port)	
	fmt.Println("文件接收服务开始....")
	listenner,err := net.Listen("tcp",*ip+":"+*port)
	if err != nil{
		fmt.Println(" net.Listen err=",err)
	}
	defer listenner.Close()
	//开始接收信息
	conn,err1 := listenner.Accept()
	if err1 != nil{
		fmt.Println("listenner.Accept err=",err1)
	}
	buf := make([]byte, 1024)
	var n int
	n, err = conn.Read(buf) //读取对方发送的文件名
	if err!=nil{
		fmt.Println("conn.Read err=",err)
		return
	}
	fileName := string(buf[:n])
	conn.Write([]byte("ok"))
	RecvFile(fileName, conn)	
}

func RecvFile(fileName string,conn net.Conn){
	f,err := os.Create(fileName)
	if err !=nil{
		fmt.Println("os.Create err=",err)
		return
	}
	buf := make([]byte,4*1024)
	for{
		n,err1 := conn.Read(buf)
		if err1!= nil{
			if err1==io.EOF{
				fmt.Println("file receive finish")
			}else{
				fmt.Println("file receive conn.Read err=",err)	
			}
		return
		}
		if n == 0 {
			fmt.Println("n == 0 文件接收完毕")
			break
		}
		f.Write(buf[:n])
	}
}

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
		fmt.Println("发送端IP地址：",*ip)	
		fmt.Println("发送端PORT地址：",*port)	
        fmt.Println("请输入需要传输的文件：")
        var path string
        fmt.Scan(&path)
        file_info,err := os.Stat(path)
        if err != nil{
                fmt.Println("文件输入错误，请重新谁人",err)
                return
        }
        conn,err1 := net.Dial("tcp",*ip+":"+*port)
        if err1 != nil{
                fmt.Println("接收服务器连接异常，请启动接收端",err1)
                return
        }
        defer conn.Close()
        //发送文件名
        _,err2 := conn.Write([]byte(file_info.Name()))
        if err2!= nil{
                fmt.Println("文件名发送异常",err2)
                return
        }
        var n int
        buf := make([]byte,2*1024)
        n,err = conn.Read(buf)
                if err != nil {
                fmt.Println("conn.Read err = ", err)
                return
        }
        if "ok"== string(buf[:n]){
                SenderFile(path,conn)
        }

}
//发送文件
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
                                fmt.Println("恭喜你文件发送完毕")
                        }else{
                                fmt.Println("源文件读取出错",err)
                        }
                        return
                }
                conn.Write(buf[:n])
        }



}

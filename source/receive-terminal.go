/**
 * anthor:liuchao@gomeplus.com
 */
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	MAX_CONN_NUM = 5
	HELP_FILE = `		***Receive File Terminal Sever***
  -port string
		PORT address of Listen server,default value is 8888
		note:send file terminal sever to use this port
  -path string
		receive file save address,default value is /tmp.
   eg:
		./listen-terminal -path /app/log/ -port 8888
`
)

//initial listener and run
func main() {
	var h bool
	flag.BoolVar(&h,"h",false,"this help")
	var port = flag.String("port","8888","input save port")
	var path = flag.String("path","/tmp","input save path")
	flag.Usage = usage
	flag.Parse()
	if h{
		flag.Usage()
		os.Exit(1)
	}
	path_dir := *path
	if path_dir[len(*path)-1:] != "/"{
		path_dir = path_dir+"/"
	}
	listener, err := net.Listen("tcp", "0.0.0.0:"+*port)
	if err != nil {
		fmt.Println("error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("running ...\n")

	conn_chan := make(chan net.Conn)

	for i := 0; i < MAX_CONN_NUM; i++ {
		go func() {
			for conn := range conn_chan {
				RecvFile(path_dir,conn)
			}
		}()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accept:", err.Error())
			return
		}
		conn_chan <- conn
	}
}

func RecvFile(savePath string,conn net.Conn){
	defer conn.Close()
	buf := make([]byte, 1024)
	var n int
	n, err := conn.Read(buf)
	if err!=nil{
		fmt.Println("conn.Read err=",err)
		return
	}
	fileName := string(buf[:n])
	conn.Write([]byte("ok"))
	fileName = savePath+fileName
	f,err := os.Create(fileName)
	if err !=nil{
		fmt.Println("os.Create err=",err)
		return
	}

	for{
		n,err1 := conn.Read(buf)
		if err1!= nil{
			if err1==io.EOF{
				fmt.Println("success receive file,",fileName," time ",time.Now().Format("2006/01/01 10:10:10"))
			}else{
				fmt.Println("file receive conn.Read err=",err)
			}
			return
		}
		if n == 0 {
			fmt.Println("success receive file,",fileName," time ",time.Now().Format("2006/01/01 10:10:10"))
			break
		}
		f.Write(buf[:n])
	}
}

func usage(){
	fmt.Fprintf(os.Stderr,HELP_FILE)
}

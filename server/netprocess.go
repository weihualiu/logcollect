package server

import (
	"github.com/weihualiu/logcollect/store"
	"github.com/weihualiu/logcollect/util"
	"log"
	"net"
	"syscall"
)

func Start() {
	//listen
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()
	//accept
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatal(err)
		}
		//每个请求一个goroutine
		go receive(conn)
	}

}

func receive(conn net.Conn) {
	log.Println("new conenction come in.")
	// one package buffer
	packageData := make([]byte, 0)
	// loop receive
	for {
		readBuf := make([]byte, 1024)
		readLen, err := conn.Read(readBuf)

		switch err {
		case nil:
			packageData = append(packageData, readBuf[:readLen]...)
			//数据拆分
			flag := true
			for flag {
				if packageData == nil {
					flag = false
				}
				pdLen := len(packageData)
				if pdLen == 0 || pdLen < 16 {
					flag = false
				} else if packageData[0] == byte(0xF0) {
					packageLen := util.BytesToUInt32(packageData[1:5])
					if packageLen == 0 {
						packageData = make([]byte, 0, 0)
						flag = false
					} else if uint32(len(packageData)) >= packageLen {
						//如果数据满足一个完整包则进入下一步处理
						parseData := make([]byte, packageLen)
						copy(parseData, packageData[0:packageLen])
						store.Parse(parseData)
						//减去完整包
						if packageLen == uint32(len(packageData)) {
							packageData = make([]byte, 0, 0)
							flag = false
						} else {
							newSize := len(packageData) - int(packageLen)
							tmp := make([]byte, newSize, newSize+1)
							tmp = packageData[packageLen:len(packageData)]
							packageData = make([]byte, newSize, newSize+1)
							copy(packageData, tmp)
						}
					} else {
						flag = false
					}
				} else {
					//错误数据，抛弃
					packageData = make([]byte, 0, 0)
					//readBuf = make([]byte, 0)
					flag = false
				}
			}
		case syscall.EAGAIN:
			continue
		default:
			log.Println(err)
			goto DISCONNECT
		}
	}

DISCONNECT:
	packageData = nil
	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Close connection: ", conn.RemoteAddr().String())

}

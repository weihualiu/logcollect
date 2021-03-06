package store

import (
	"github.com/weihualiu/logcollect/model"
	"github.com/weihualiu/logcollect/util"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Parse(data []byte) {
	if data == nil {
		log.Fatal("data is nil")
	}
	if len(data) < 5 {
		log.Fatal(data)
	}
	if data[5] == byte(0x00) {
		//收到心跳包，回复ACK
		//log.Println("heartbeat data")
	} else {
		pack, err := model.NewPackCommon(data)
		if err != nil {
			log.Fatal(err)
		}

		if pack.Type == byte(0x01) {
			//api
			//tag1:项目名称,tag2:项目环境,tag3:接口名称
			if int(pack.TagNum) < 3 {
				log.Fatal("data tag number error!")
			}
			tag1 := util.BytesToString(pack.TagList[0].Name)
			tag2 := util.BytesToString(pack.TagList[1].Name)
			tag3 := util.BytesToString(pack.TagList[2].Name)
			basepath := filepath.Join("data/api", tag1, tag2, util.BytesToString(pack.Date))
			//log.Println(basepath)
			strtime := "------------------" + time.Now().Truncate(time.Second).String() + "----------------------"

			// memory save
			Monitors.Write(filepath.Join(tag1, tag2, util.BytesToString(pack.Date), tag3), []byte(strtime))
			Monitors.Write(filepath.Join(tag1, tag2, util.BytesToString(pack.Date), tag3), pack.Body)
			// local stroage
			err = os.MkdirAll(basepath, os.ModePerm)
			//创建存储目录
			if err != nil {
				log.Println(err)
			}
			filePath := filepath.Join(basepath, tag3)
			//写入文件
			err = util.AppendToFile(filePath, []byte(strtime))
			if err != nil {
				log.Println(err)
			}

			err = util.AppendToFile(filePath, pack.Body)
			if err != nil {
				log.Println(err)
			}

		}
	}

}

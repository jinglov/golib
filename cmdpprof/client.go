package cmdpprof

import (
	"fmt"
	"github.com/jinglov/golib/cmdservice"
	"io"
	"os"
)

type clientHandler func()

var idMap = make(map[uint8]clientHandler)

func NewCmdClient(net, addr string) error {
	cmdservice.NewClient(net, addr)
	defer cmdservice.Close()
	for _, p := range pprofHandlers {
		cmdservice.ClientHandler(p.id, p.name)
		idMap[p.id] = p.clientHandler
		//idMap[p.name] = p.clientHandler
	}
	showList()
	/*	b, err := cmdservice.Send(*status, make([]byte, 0))
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Debug(string(b))*/
	return nil
}

func showList() {
	var cmd uint8
	for {
		fmt.Println("请选择")
		for _, p := range pprofHandlers {
			fmt.Printf("%d. %s\n", p.id, p.description)
		}
		fmt.Println("0.退出")
		fmt.Print("请输入前的数字：")
		fmt.Scanln(&cmd)
		if handler, ok := idMap[cmd]; ok {
			if handler != nil {
				handler()
			}
		} else {
			fmt.Println("退出")
			return
		}
	}
}

func cpuClient() {
	s := inputSecond()
	fp := inputFile("cpu.pprof")
	send("cpu", []byte(s), fp)
}

func memClient() {
	fp := inputFile("mem.pprof")
	send("mem", make([]byte, 0), fp)
	fmt.Println("文件保存成功。")
}

func inputFile(defFile string) io.Writer {
	var inputFmt, fileName string
	if defFile == "" {
		inputFmt = "请输入保存的文件名："
	} else {
		inputFmt = fmt.Sprintf("请输入保存的文件名 (%s)：", defFile)
	}
	for {
		fmt.Print(inputFmt)
		fmt.Scanln(&fileName)
		if fileName == "" {
			fileName = defFile
		}
		if fileName != "" {
			fp, err := os.Create(fileName)
			if err == nil {
				return fp
			}
		}
	}
}

func inputSecond() string {
	var s string
	fmt.Print("请输入取样时间（秒）：")
	fmt.Scanln(&s)
	return s
}

func send(name string, param []byte, w io.Writer) error {
	b, err := cmdservice.Send(name, param)
	if err != nil {
		fmt.Println(err)
		return err
	}
	w.Write(b)
	return nil
}

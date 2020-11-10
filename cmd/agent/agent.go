package main

import "github.com/26597925/EastCloud/internal/agent"

//C:\\Windows\\System32\\notepad.exe D:\\1.txt
func main()  {
	//svr := agent.NewServer()
	//svr.Start()
	//
	//cli := agent.NewClient()
	//cli.Start()
	//
	//t := time.NewTimer(5 * time.Second)
	//for {
	//	select {
	//	case <-t.C:
	//		cli.Test()
	//	}
	//}
	//

	//info := &core.Info{
	//	ProcInfo: process.ProcInfo{
	//		Pid:     0,
	//		//Cmdline: "./ttorrent test.torrent",
	//		Cmdline: "C:\\Windows\\System32\\notepad.exe\nD:\\1.txt",
	//		Env: "",
	//	},
	//	Name:     "ttorrent",
	//	Status:   0,
	//	SavePath: "",
	//}
	//
	//process := core.NewProcess("E:\\GoWork\\src\\douyu\\bsw\\cmd\\agent\\work")
	//process.AddProcess(info)
	//
	//pid, err := process.Start("ttorrent")
	//fmt.Println(err)
	//fmt.Println(pid)
	//
	//fmt.Println(info)

	cli := agent.NewClient()
	cli.Start()

	for  {
		select {

		}
	}

}

package agent

import (
	"github.com/26597925/EastCloud/internal/agent/core"
	"github.com/26597925/EastCloud/internal/agent/msg"
	"github.com/26597925/EastCloud/internal/agent/router/tcpclient"
	"github.com/26597925/EastCloud/pkg/network/tcp"
	"github.com/26597925/EastCloud/pkg/process"
	"time"
)

type Client struct {
	client *tcp.Client
	manager *core.Manager
	process *core.Process

	quit chan bool
}

func NewClient() *Client {
	cli := &Client{
		client: tcp.NewClient(),
		manager: core.NewManager("E://GoWork//src//cnas//"),
		process: core.NewProcess("E:\\GoWork\\src\\douyu\\bsw\\cmd\\agent\\work"),
		quit: make(chan bool, 1),
	}

	return cli
}

func (c *Client) Start() {
	info := &core.Info{
		ProcInfo: process.ProcInfo{
			Pid:     0,
			//Cmdline: "./ttorrent test.torrent",
			Cmdline: "C:\\Windows\\System32\\notepad.exe\nD:\\1.txt",
			Env: "",
		},
		Name:     "notepad.exe",
		Status:   0,
		SavePath: "",
	}
	c.process.AddProcess(info)

	c.client.AddRouter(msg.Heartbeat, &tcpclient.HeartbeatRouter{})
	c.client.AddRouter(msg.Command, &tcpclient.CommandRouter{
		Process: c.process,
	})
	c.client.AddRouter(msg.File, tcpclient.NewFileRouter(c.manager))

	go c.client.Connect()

	th := time.NewTicker(5 * time.Second)
	go func(t *time.Ticker) {
		for {
			select {
			case <-t.C:
				err := c.client.SendBuffMsg(msg.Heartbeat, nil)
				if err != nil {
					c.client.Close()
				}
			case <-c.quit:
				c.client.Close()
				return
			}
		}
	}(th)

	tm := time.NewTicker(500 * time.Millisecond)
	go func(t *time.Ticker) {
		for {
			select {
			case <-t.C:
				c.process.Monitor()
			}
		}
	}(tm)
}

func (c *Client) Stop() {
	close(c.quit)
}
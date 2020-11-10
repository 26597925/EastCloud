package websocket

import (
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/network/common"
	"github.com/26597925/EastCloud/pkg/util/atomic"
	"github.com/gorilla/websocket"
	"html/template"
	"net"
	"net/http"
)

type Server struct {
	cid *atomic.Uint32
	msgHandler *common.MsgHandle
	ConnMgr *common.ConnManager
	OnConnStart func(conn *common.Connection)
	OnConnStop func(conn *common.Connection)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1500000,
	WriteBufferSize: 1500000,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebsocket() *Server {
	s := &Server{
		msgHandler: common.NewMsgHandle(10, 1024),
		ConnMgr:    common.NewConnManager(),
	}
	return s
}

func (s *Server) Start() {
	logger.InfoF("[START] Server name: %s,listenner at addr %s is starting\n", "aa", ":3001")
	s.cid = &atomic.Uint32{}
	s.cid.Set(0)

	go func() {
		s.msgHandler.StartWorkerPool()

		httpServeMux := http.NewServeMux()
		httpServeMux.HandleFunc("/", home)
		httpServeMux.HandleFunc("/ws", s.ServeWebSocket)
		addr, err := net.ResolveTCPAddr("tcp4", ":3001")
		if err != nil {
			logger.Error("resolve tcp addr err: ", err)
			return
		}

		listener, err := net.ListenTCP("tcp4", addr)
		server := &http.Server{Handler: httpServeMux}

		go func(host string) {
			if err = server.Serve(listener); err != nil {
				logger.Error("server.Serve(\"%s\") error(%v)", host, err)
				panic(err)
			}
		}(":3001")
	}()
}

func (s *Server) ServeWebSocket(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		logger.Error("Websocket Upgrade error(%v), userAgent(%s)", err, req.UserAgent())
		return
	}
	defer ws.Close()

	logger.Info("Get conn remote addr = ", ws.RemoteAddr().String())

	if s.ConnMgr.Len() >= 12000 {
		ws.Close()
	}

	dealConn := common.NewConnection(common.Websocket, ws.UnderlyingConn(), s.cid.Get(), 1024, 4096)
	dealConn.SetWebsocketConn(ws)
	dealConn.SetOnClose(s.OnConnStop)

	s.ConnMgr.Add(dealConn)
	s.cid.Incr()

	go s.startWriter(dealConn)
	if s.OnConnStart != nil {
		s.OnConnStart(dealConn)
	}

	s.startReader(dealConn)
}

func (s *Server) Stop() {
	logger.Info("[STOP] server , name ", "aa")

	s.ConnMgr.ClearConn()
}

func (s *Server) AddRouter(msgId byte, router common.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnMgr() *common.ConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hookFunc func(*common.Connection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(*common.Connection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) startReader(c *common.Connection) {
	logger.Info("[Reader Goroutine is running]")
	defer logger.Info(c.RemoteAddr().String(), "[conn Reader exit!]")
	defer c.Close()

	for {
		dp := newDataPack()

		if types, data, err := c.WebsocketConn.ReadMessage(); err != nil {
			logger.Error("read msg error ", err)
			break
		}else {
			if types == websocket.TextMessage {
				msg, err := dp.Unpack(data)
				if err != nil {
					logger.Error("unpack error ", err)
					break
				}

				req := common.Request{
					Conn: c,
					Msg:  msg,
				}

				go s.msgHandler.DoReceive(&req)
			}
		}
	}
}

func (s *Server) startWriter(c *common.Connection) {
	logger.Info("[Writer Goroutine is running]")
	defer logger.Info(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case msg := <-c.GetMsgChan():
			dp := newDataPack()
			data, err := dp.Pack(msg)
			if err != nil {
				logger.Error("Pack error msg id = ", msg.Id)
				return
			}
			if err := c.WebsocketConn.WriteMessage(websocket.TextMessage, data); err != nil {
				logger.Error("Send Data error:, ", err, " Conn Writer exit")
				return
			}
			//fmt.Printf("Send data succ! data = %+v\n", data)
		case msg, ok := <-c.GetMsgBuffChan():
			if ok {
				dp := newDataPack()
				data, err := dp.Pack(msg)
				if err != nil {
					logger.Error("Pack error msg id = ", msg.Id)
					return
				}
				if err := c.WebsocketConn.WriteMessage(websocket.BinaryMessage, data); err != nil {
					logger.Error("Server Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				logger.Info("msgBuffChan is Closed")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value='{"Id":3,"Data":"{\"cmd\":\"adb shell\"}","Len":19}'>
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
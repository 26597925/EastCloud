package common

type IRouter interface {
	PreHandle(req *Request)
	Handle(req *Request)
	PostHandle(req *Request)
}

type BaseRouter struct{}
func (br *BaseRouter) PreHandle(req *Request)  {}
func (br *BaseRouter) Handle(req *Request)     {}
func (br *BaseRouter) PostHandle(req *Request) {}

type Request struct {
	Conn *Connection
	Msg  *Message
}

func (r *Request) GetConnection() *Connection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.Msg.GetData()
}

func (r *Request) GetMsgID() byte {
	return r.Msg.GetMsgId()
}
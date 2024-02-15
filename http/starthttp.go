package httpServer

import (
	"KVStore/raft"
)

type StartNode struct {
	id       int
	endpoint string
	peerIds  map[int]string
	hs       *HTTPServer
	s        *raft.Server
}

func NewStartNode(id int, port int, peerIds map[int]string, filePath string) *StartNode {
	storage, _, err := raft.NewStorage(filePath)
	if err != nil {
		return nil
	}
	sd := new(StartNode)
	sd.id = id
	sd.endpoint = GetEndpoint(port + 10)
	sd.peerIds = peerIds
	sd.s = raft.NewServer(id, GetEndpoint(port), peerIds, storage)
	sd.hs = NewHTTPServer(id, sd.endpoint, sd.s)
	return sd
}

func (sd *StartNode) Report() (int, int, bool) {
	return sd.s.Report()
}

func (sd *StartNode) getState() raft.CMState {
	return sd.s.GetState()
}

func (sd *StartNode) Start() {
	sd.s.Serve()
	go sd.hs.Start()
	// 等待程序连接
}
func (sd *StartNode) GetForTest(key string) (reply raft.CommandReply) {
	sd.s.Get(key, &reply)
	return
}
func (sd *StartNode) PutForTest(key, value string) (reply raft.CommandReply) {
	sd.s.Put(key, value, &reply)
	return
}
func (sd *StartNode) Shutdown() {
	sd.s.DisconnectAll()
	sd.s.Shutdown()
	sd.hs.Shutdown()
}

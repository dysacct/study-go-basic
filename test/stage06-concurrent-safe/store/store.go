package store

import (
	"project-001/model"
	"sync"
)

type ServerStore struct {
	mu sync.RWMutex

	servers []model.Server
}

// New 是 ServerStore 的构造函数：返回一个带初始种子数据的 Store 指针。
// 为什么返回指针 *ServerStore？因为里面有 sync.RWMutex——
// Mutex 绝对不能被复制（复制会得到两把独立的锁，锁不住同一份数据），
// 所以永远用指针传递带锁的结构体。
func New() *ServerStore {
	return &ServerStore{
		servers: []model.Server{
			{Name: "web-01", IP: "10.0.0.5", Status: "running"},
			{Name: "db-01", IP: "10.0.0.21", Status: "running"},
			{Name: "cache-01", IP: "10.0.0.30", Status: "stopped"},
		},
	}
}

func (s *ServerStore) List() []model.Server {

	s.mu.RLock()

	defer s.mu.RUnlock()

	return append(
		[]model.Server(nil),
		s.servers...,
	)

}

func (s *ServerStore) Add(server model.Server) {

	s.mu.Lock()

	defer s.mu.Unlock()

	s.servers = append(
		s.servers,
		server,
	)

}

func (s *ServerStore) Find(name string) (model.Server, bool) {

	s.mu.RLock()

	defer s.mu.RUnlock()

	for _, server := range s.servers {

		if server.Name == name {

			return server, true

		}

	}

	return model.Server{}, false

}

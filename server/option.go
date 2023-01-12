package server

import (
	"fmt"
	"net"
	"time"
)

// Option 修改 server 相关参数
type Option func(*Server)

// Port 修改端口
func Port(v string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", v)
	}
}

// ShutdownTimeout 停止服务超时时间
func ShutdownTimeout(v time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = v
	}
}

func ReadTimeout(v time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = v
	}
}

func WriteTimeout(v time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = v
	}
}

// DefaultPrintln 默认输出信息
func DefaultPrintln() Option {
	return func(s *Server) {
		fmt.Printf("server start : addr(%s)\n", s.server.Addr)
	}
}

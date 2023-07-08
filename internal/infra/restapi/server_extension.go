package restapi

import "net"

func (s *Server) SetHttpListener(l net.Listener) {
	s.httpsServerL = l
}

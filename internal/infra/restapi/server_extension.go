package restapi

import "net"

func (s *Server) SetHTTPListener(l net.Listener) {
	s.httpsServerL = l
}

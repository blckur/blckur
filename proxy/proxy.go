// Wrapper for elb proxy protocol server.
package proxy

import (
	"github.com/armon/go-proxyproto"
	"github.com/blckur/blckur/errortypes"
	"github.com/dropbox/godropbox/errors"
	"net"
	"net/http"
	"time"
)

type keepAliveListener struct {
	*net.TCPListener
}

func (l keepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := l.AcceptTCP()
	if err != nil {
		return
	}

	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	c = tc

	return
}

func ListenAndServe(svr *http.Server) (err error) {
	addr := svr.Addr
	if addr == "" {
		addr = ":http"
	}

	list, err := net.Listen("tcp", addr)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "proxy: Listen error"),
		}
		return
	}

	keepAliveList := keepAliveListener{list.(*net.TCPListener)}
	proxyList := &proxyproto.Listener{keepAliveList}

	err = svr.Serve(proxyList)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "proxy: Serve error"),
		}
		return
	}

	return
}

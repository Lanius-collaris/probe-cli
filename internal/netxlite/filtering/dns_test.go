package filtering

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/apex/log"
	"github.com/miekg/dns"
	"github.com/ooni/probe-cli/v3/internal/netxlite"
	"github.com/ooni/probe-cli/v3/internal/netxlite/mocks"
)

func TestDNSProxy(t *testing.T) {
	newproxy := func(action DNSAction) (DNSListener, <-chan interface{}, error) {
		p := &DNSProxy{
			OnQuery: func(domain string) DNSAction {
				return action
			},
		}
		return p.start("127.0.0.1:0")
	}

	newresolver := func(listener DNSListener) netxlite.Resolver {
		dlr := netxlite.NewDialerWithoutResolver(log.Log)
		r := netxlite.NewResolverUDP(log.Log, dlr, listener.LocalAddr().String())
		return r
	}

	t.Run("DNSActionProxy with default proxy", func(t *testing.T) {
		ctx := context.Background()
		listener, done, err := newproxy(DNSActionProxy)
		if err != nil {
			t.Fatal(err)
		}
		r := newresolver(listener)
		addrs, err := r.LookupHost(ctx, "dns.google")
		if err != nil {
			t.Fatal(err)
		}
		if addrs == nil {
			t.Fatal("unexpected empty addrs")
		}
		var foundQuad8 bool
		for _, addr := range addrs {
			foundQuad8 = foundQuad8 || addr == "8.8.8.8"
		}
		if !foundQuad8 {
			t.Fatal("did not find 8.8.8.8")
		}
		listener.Close()
		<-done // wait for background goroutine to exit
	})

	t.Run("DNSActionNXDOMAIN", func(t *testing.T) {
		ctx := context.Background()
		listener, done, err := newproxy(DNSActionNXDOMAIN)
		if err != nil {
			t.Fatal(err)
		}
		r := newresolver(listener)
		addrs, err := r.LookupHost(ctx, "dns.google")
		if err == nil || err.Error() != netxlite.FailureDNSNXDOMAINError {
			t.Fatal("unexpected err", err)
		}
		if addrs != nil {
			t.Fatal("expected empty addrs")
		}
		listener.Close()
		<-done // wait for background goroutine to exit
	})

	t.Run("DNSActionRefused", func(t *testing.T) {
		ctx := context.Background()
		listener, done, err := newproxy(DNSActionRefused)
		if err != nil {
			t.Fatal(err)
		}
		r := newresolver(listener)
		addrs, err := r.LookupHost(ctx, "dns.google")
		if err == nil || err.Error() != netxlite.FailureDNSRefusedError {
			t.Fatal("unexpected err", err)
		}
		if addrs != nil {
			t.Fatal("expected empty addrs")
		}
		listener.Close()
		<-done // wait for background goroutine to exit
	})

	t.Run("DNSActionLocalHost", func(t *testing.T) {
		ctx := context.Background()
		listener, done, err := newproxy(DNSActionLocalHost)
		if err != nil {
			t.Fatal(err)
		}
		r := newresolver(listener)
		addrs, err := r.LookupHost(ctx, "dns.google")
		if err != nil {
			t.Fatal(err)
		}
		if addrs == nil {
			t.Fatal("expected non-empty addrs")
		}
		var found127001 bool
		for _, addr := range addrs {
			found127001 = found127001 || addr == "127.0.0.1"
		}
		if !found127001 {
			t.Fatal("did not find 127.0.0.1")
		}
		listener.Close()
		<-done // wait for background goroutine to exit
	})

	t.Run("DNSActionEmpty", func(t *testing.T) {
		ctx := context.Background()
		listener, done, err := newproxy(DNSActionEmpty)
		if err != nil {
			t.Fatal(err)
		}
		r := newresolver(listener)
		addrs, err := r.LookupHost(ctx, "dns.google")
		if err == nil || err.Error() != netxlite.FailureDNSNoAnswer {
			t.Fatal(err)
		}
		if addrs != nil {
			t.Fatal("expected empty addrs")
		}
		listener.Close()
		<-done // wait for background goroutine to exit
	})

	t.Run("DNSActionTimeout", func(t *testing.T) {
		// Implementation note: if you see this test running for more
		// than one second, then it means we're not checking the context
		// immediately. We should be improving there but we need to be
		// careful because lots of legacy code uses SerialResolver.
		const timeout = time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		listener, done, err := newproxy(DNSActionTimeout)
		defer cancel()
		if err != nil {
			t.Fatal(err)
		}
		r := newresolver(listener)
		addrs, err := r.LookupHost(ctx, "dns.google")
		if err == nil || err.Error() != netxlite.FailureGenericTimeoutError {
			t.Fatal(err)
		}
		if addrs != nil {
			t.Fatal("expected empty addrs")
		}
		listener.Close()
		<-done // wait for background goroutine to exit
	})

	t.Run("Start with invalid address", func(t *testing.T) {
		p := &DNSProxy{}
		listener, err := p.Start("127.0.0.1")
		if err == nil {
			t.Fatal("expected an error")
		}
		if listener != nil {
			t.Fatal("expected nil listener")
		}
	})

	t.Run("oneloop", func(t *testing.T) {
		t.Run("ReadFrom failure after which we should continue", func(t *testing.T) {
			expected := errors.New("mocked error")
			p := &DNSProxy{}
			conn := &mocks.QUICUDPLikeConn{
				MockReadFrom: func(p []byte) (n int, addr net.Addr, err error) {
					return 0, nil, expected
				},
			}
			okay := p.oneloop(conn)
			if !okay {
				t.Fatal("we should be okay after this error")
			}
		})

		t.Run("ReadFrom the connection is closed", func(t *testing.T) {
			expected := errors.New("use of closed network connection")
			p := &DNSProxy{}
			conn := &mocks.QUICUDPLikeConn{
				MockReadFrom: func(p []byte) (n int, addr net.Addr, err error) {
					return 0, nil, expected
				},
			}
			okay := p.oneloop(conn)
			if okay {
				t.Fatal("we should not be okay after this error")
			}
		})

		t.Run("Unpack fails", func(t *testing.T) {
			p := &DNSProxy{}
			conn := &mocks.QUICUDPLikeConn{
				MockReadFrom: func(p []byte) (n int, addr net.Addr, err error) {
					if len(p) < 4 {
						panic("buffer too small")
					}
					p[0] = 7
					return 1, &net.UDPAddr{}, nil
				},
			}
			okay := p.oneloop(conn)
			if !okay {
				t.Fatal("we should be okay after this error")
			}
		})

		t.Run("reply fails", func(t *testing.T) {
			p := &DNSProxy{}
			conn := &mocks.QUICUDPLikeConn{
				MockReadFrom: func(p []byte) (n int, addr net.Addr, err error) {
					query := &dns.Msg{}
					query.Question = append(query.Question, dns.Question{})
					query.Question = append(query.Question, dns.Question{})
					data, err := query.Pack()
					if err != nil {
						panic(err)
					}
					if len(p) < len(data) {
						panic("buffer too small")
					}
					for i := 0; i < len(data); i++ {
						p[i] = data[i]
					}
					return len(data), &net.UDPAddr{}, nil
				},
			}
			okay := p.oneloop(conn)
			if !okay {
				t.Fatal("we should be okay after this error")
			}
		})

		t.Run("pack fails", func(t *testing.T) {
			p := &DNSProxy{
				mockableReply: func(query *dns.Msg) (*dns.Msg, error) {
					reply := &dns.Msg{}
					reply.MsgHdr.Rcode = -1 // causes pack to fail
					return reply, nil
				},
			}
			conn := &mocks.QUICUDPLikeConn{
				MockReadFrom: func(p []byte) (n int, addr net.Addr, err error) {
					query := &dns.Msg{}
					query.Question = append(query.Question, dns.Question{})
					data, err := query.Pack()
					if err != nil {
						panic(err)
					}
					if len(p) < len(data) {
						panic("buffer too small")
					}
					for i := 0; i < len(data); i++ {
						p[i] = data[i]
					}
					return len(data), &net.UDPAddr{}, nil
				},
			}
			okay := p.oneloop(conn)
			if !okay {
				t.Fatal("we should be okay after this error")
			}
		})
	})

	t.Run("proxy", func(t *testing.T) {
		t.Run("pack fails", func(t *testing.T) {
			p := &DNSProxy{}
			query := &dns.Msg{}
			query.Rcode = -1 // causes Pack to fail
			reply, err := p.proxy(query)
			if err == nil {
				t.Fatal("expected error here")
			}
			if reply != nil {
				t.Fatal("expected nil reply")
			}
		})

		t.Run("round trip fails", func(t *testing.T) {
			expected := errors.New("mocked error")
			p := &DNSProxy{
				Upstream: &mocks.DNSTransport{
					MockRoundTrip: func(ctx context.Context, query []byte) (reply []byte, err error) {
						return nil, expected
					},
					MockCloseIdleConnections: func() {},
				},
			}
			reply, err := p.proxy(&dns.Msg{})
			if !errors.Is(err, expected) {
				t.Fatal("unexpected err", err)
			}
			if reply != nil {
				t.Fatal("expected nil reply here")
			}
		})

		t.Run("unpack fails", func(t *testing.T) {
			p := &DNSProxy{
				Upstream: &mocks.DNSTransport{
					MockRoundTrip: func(ctx context.Context, query []byte) (reply []byte, err error) {
						return make([]byte, 1), nil
					},
					MockCloseIdleConnections: func() {},
				},
			}
			reply, err := p.proxy(&dns.Msg{})
			if err == nil {
				t.Fatal("expected error")
			}
			if reply != nil {
				t.Fatal("expected nil reply here")
			}
		})
	})
}
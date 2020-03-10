package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/mit-scripts/scripts/server/common/oursrc/scripts-proxy/ldap"
	"inet.af/tcpproxy"
)

var (
	httpAddrs   = flag.String("http_addrs", ":80", "comma-separated addresses to listen for HTTP traffic on")
	sniAddrs    = flag.String("sni_addrs", ":443,:444", "comma-separated addresses to listen for SNI traffic on")
	ldapServers = flag.String("ldap_servers", "scripts-ldap.mit.edu:389", "comma-spearated LDAP servers to query")
	defaultHost = flag.String("default_host", "scripts.mit.edu", "default host to route traffic to if SNI/Host header cannot be parsed or cannot be found in LDAP")
	baseDn      = flag.String("base_dn", "ou=VirtualHosts,dc=scripts,dc=mit,dc=edu", "base DN to query for hosts")
	localRange  = flag.String("local_range", "18.4.86.0/24", "IP block for client IP spoofing")
)

const ldapRetries = 3

func always(context.Context, string) bool {
	return true
}

type ldapTarget struct {
	localPoolRange *net.IPNet
	ldap           *ldap.Pool
}

func (l *ldapTarget) HandleConn(netConn net.Conn) {
	var pool string
	var err error
	if conn, ok := netConn.(*tcpproxy.Conn); ok {
		pool, err = l.ldap.ResolvePool(conn.HostName)
		if err != nil {
			log.Printf("resolving %q: %v", conn.HostName, err)
		}
	}
	if pool == "" {
		pool, err = l.ldap.ResolvePool(*defaultHost)
		if err != nil {
			log.Printf("resolving default pool: %v", err)
		}
	}
	if pool == "" {
		netConn.Close()
		return
	}
	laddr := netConn.LocalAddr().(*net.TCPAddr)
	destAddrStr := net.JoinHostPort(pool, fmt.Sprintf("%d", laddr.Port))
	destAddr, err := net.ResolveTCPAddr("tcp", destAddrStr)
	if err != nil {
		netConn.Close()
		log.Printf("parsing pool address %q: %v", pool, err)
		return
	}
	dp := &tcpproxy.DialProxy{
		Addr: destAddrStr,
	}
	raddr := netConn.RemoteAddr().(*net.TCPAddr)
	if l.localPoolRange.Contains(destAddr.IP) {
		sourceAddr := &net.TCPAddr{
			IP: raddr.IP,
		}
		dp.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.DialTCP(network, sourceAddr, destAddr)
		}
	}
	dp.HandleConn(netConn)
}

func main() {
	flag.Parse()

	_, ipnet, err := net.ParseCIDR(*localRange)
	if err != nil {
		log.Fatal(err)
	}

	ldapPool := ldap.NewPool(strings.Split(*ldapServers, ","), *baseDn, ldapRetries)

	var p tcpproxy.Proxy
	t := &ldapTarget{
		localPoolRange: ipnet,
		ldap:           ldapPool,
	}
	for _, addr := range strings.Split(*httpAddrs, ",") {
		p.AddHTTPHostMatchRoute(addr, always, t)
	}
	for _, addr := range strings.Split(*sniAddrs, ",") {
		p.AddStopACMESearch(addr)
		p.AddSNIMatchRoute(addr, always, t)
	}
	log.Fatal(p.Run())
}

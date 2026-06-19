package util

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

func HttpGet(feedURL string) (*http.Response, error) {
	parsedURL, err := url.Parse(feedURL)
	if err != nil {
		return nil, err
	}

	if err := validateURL(parsedURL); err != nil {
		return nil, err
	}

	dialer := &net.Dialer{}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			ips, err := net.LookupIP(host)
			if err != nil {
				return nil, err
			}

			if len(ips) == 0 {
				return nil, errors.New("no ip found")
			}

			for _, ip := range ips {
				if isPrivateIP(ip) {
					return nil, fmt.Errorf("private ip detected: %s", ip)
				}
			}

			return dialer.DialContext(
				ctx,
				network,
				net.JoinHostPort(ips[0].String(), port),
			)
		},
	}
	defer transport.CloseIdleConnections()

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= MaxRedirect {
				return errors.New("too many redirects")
			}

			return validateURL(req.URL)
		},
	}

	resp, err := client.Get(feedURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp, nil
}

func validateURL(u *url.URL) error {
	if u == nil {
		return errors.New("nil url")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("invalid scheme")
	}

	host := u.Hostname()
	if host == "" {
		return errors.New("missing host")
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return err
	}

	if len(ips) == 0 {
		return errors.New("no ip found")
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return fmt.Errorf("private address: %s", ip)
		}
	}

	port := u.Port()
	if port != "" && port != "80" && port != "443" {
		return fmt.Errorf("invalid port: %s", port)
	}

	return nil
}

func isPrivateIP(ip net.IP) bool {
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}

	privateCIDRs := []string{
		"0.0.0.0/8",
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",

		"100.64.0.0/10",
		"198.18.0.0/15",
		"224.0.0.0/4",
		"240.0.0.0/4",

		"::/128",
		"::1/128",
		"fc00::/7",
		"ff00::/8",
		"fe80::/10",
	}

	for _, cidr := range privateCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}

		if network.Contains(ip) {
			return true
		}
	}

	return false
}

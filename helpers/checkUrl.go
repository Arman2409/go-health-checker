package helpers

import (
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	Index int
	Msg   string
}

func CheckURL(url string, i int, c chan Result) {
	client := &http.Client{Timeout: 5 * time.Second}
	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)
	if err != nil {
		c <- Result{i, fmt.Sprintf("❌ %s is down: %v", url, err)}
		return
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	contentType := resp.Header.Get("Content-Type")
	contentLen := resp.ContentLength

	tlsExpiry := ""
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		tlsExpiry = resp.TLS.PeerCertificates[0].NotAfter.Format(time.RFC3339)
	}

	msg := fmt.Sprintf("✅ %s is up! status=%d time=%dms type=%s len=%d", url, status, elapsed.Milliseconds(), contentType, contentLen)
	if tlsExpiry != "" {
		msg = msg + " tls_expiry=" + tlsExpiry
	}

	msg = msg + "\n"

	c <- Result{i, msg}
}
package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

// SSLCheck represents SSL certificate check results
type SSLCheck struct {
	URL                string        `json:"url"`
	Host               string        `json:"host"`
	Port               string        `json:"port"`
	Valid              bool          `json:"valid"`
	ExpiresAt          time.Time     `json:"expires_at"`
	DaysUntilExpiry    int           `json:"days_until_expiry"`
	Issuer             string        `json:"issuer"`
	Subject            string        `json:"subject"`
	SerialNumber       string        `json:"serial_number"`
	SignatureAlgorithm string        `json:"signature_algorithm"`
	Version            int           `json:"version"`
	KeyUsage           []string      `json:"key_usage"`
	ExtKeyUsage        []string      `json:"ext_key_usage"`
	DNSNames           []string      `json:"dns_names"`
	IPAddresses        []string      `json:"ip_addresses"`
	Chain              []Certificate `json:"chain"`
	Error              string        `json:"error,omitempty"`
	CheckedAt          time.Time     `json:"checked_at"`
	ResponseTime       time.Duration `json:"response_time"`
}

// Certificate represents a certificate in the chain
type Certificate struct {
	Subject      string    `json:"subject"`
	Issuer       string    `json:"issuer"`
	SerialNumber string    `json:"serial_number"`
	NotBefore    time.Time `json:"not_before"`
	NotAfter     time.Time `json:"not_after"`
	IsCA         bool      `json:"is_ca"`
	KeyUsage     []string  `json:"key_usage"`
	ExtKeyUsage  []string  `json:"ext_key_usage"`
	Fingerprint  string    `json:"fingerprint"`
}

// SSLChecker performs SSL certificate checks
type SSLChecker struct {
	Timeout time.Duration
}

// NewSSLChecker creates a new SSL checker
func NewSSLChecker(timeout time.Duration) *SSLChecker {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &SSLChecker{
		Timeout: timeout,
	}
}

// CheckSSL performs SSL certificate validation for a URL
func (c *SSLChecker) CheckSSL(rawURL string) SSLCheck {
	start := time.Now()

	check := SSLCheck{
		URL:       rawURL,
		CheckedAt: start,
	}

	// Parse URL to extract host and port
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		check.Error = fmt.Sprintf("Invalid URL: %v", err)
		check.ResponseTime = time.Since(start)
		return check
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()

	// Default ports
	if port == "" {
		switch parsedURL.Scheme {
		case "https":
			port = "443"
		case "http":
			check.Error = "HTTP URL provided, HTTPS required for SSL check"
			check.ResponseTime = time.Since(start)
			return check
		default:
			port = "443"
		}
	}

	check.Host = host
	check.Port = port

	// Connect with TLS
	dialer := &net.Dialer{Timeout: c.Timeout}

	conn, err := tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(host, port), &tls.Config{
		ServerName: host,
		// Don't verify for now, we'll do our own validation
		InsecureSkipVerify: true,
	})

	if err != nil {
		check.Error = fmt.Sprintf("TLS connection failed: %v", err)
		check.ResponseTime = time.Since(start)
		return check
	}
	defer conn.Close()

	check.ResponseTime = time.Since(start)

	// Get certificate chain
	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		check.Error = "No certificates found"
		return check
	}

	// Analyze the leaf certificate
	cert := state.PeerCertificates[0]
	check.Valid = c.validateCertificate(cert, host)
	check.ExpiresAt = cert.NotAfter
	check.DaysUntilExpiry = int(time.Until(cert.NotAfter).Hours() / 24)
	check.Issuer = cert.Issuer.String()
	check.Subject = cert.Subject.String()
	check.SerialNumber = cert.SerialNumber.String()
	check.SignatureAlgorithm = cert.SignatureAlgorithm.String()
	check.Version = cert.Version
	check.KeyUsage = c.parseKeyUsage(cert.KeyUsage)
	check.ExtKeyUsage = c.parseExtKeyUsage(cert.ExtKeyUsage)
	check.DNSNames = cert.DNSNames

	// Convert IP addresses to strings
	for _, ip := range cert.IPAddresses {
		check.IPAddresses = append(check.IPAddresses, ip.String())
	}

	// Build certificate chain
	for _, chainCert := range state.PeerCertificates {
		certInfo := Certificate{
			Subject:      chainCert.Subject.String(),
			Issuer:       chainCert.Issuer.String(),
			SerialNumber: chainCert.SerialNumber.String(),
			NotBefore:    chainCert.NotBefore,
			NotAfter:     chainCert.NotAfter,
			IsCA:         chainCert.IsCA,
			KeyUsage:     c.parseKeyUsage(chainCert.KeyUsage),
			ExtKeyUsage:  c.parseExtKeyUsage(chainCert.ExtKeyUsage),
			Fingerprint:  c.calculateFingerprint(chainCert),
		}
		check.Chain = append(check.Chain, certInfo)
	}

	return check
}

// validateCertificate validates the certificate against the hostname
func (c *SSLChecker) validateCertificate(cert *x509.Certificate, hostname string) bool {
	now := time.Now()

	// Check time validity
	if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
		return false
	}

	// Check hostname
	if err := cert.VerifyHostname(hostname); err != nil {
		return false
	}

	return true
}

// parseKeyUsage converts KeyUsage flags to string slice
func (c *SSLChecker) parseKeyUsage(usage x509.KeyUsage) []string {
	var usages []string

	if usage&x509.KeyUsageDigitalSignature != 0 {
		usages = append(usages, "Digital Signature")
	}
	if usage&x509.KeyUsageContentCommitment != 0 {
		usages = append(usages, "Content Commitment")
	}
	if usage&x509.KeyUsageKeyEncipherment != 0 {
		usages = append(usages, "Key Encipherment")
	}
	if usage&x509.KeyUsageDataEncipherment != 0 {
		usages = append(usages, "Data Encipherment")
	}
	if usage&x509.KeyUsageKeyAgreement != 0 {
		usages = append(usages, "Key Agreement")
	}
	if usage&x509.KeyUsageCertSign != 0 {
		usages = append(usages, "Certificate Sign")
	}
	if usage&x509.KeyUsageCRLSign != 0 {
		usages = append(usages, "CRL Sign")
	}
	if usage&x509.KeyUsageEncipherOnly != 0 {
		usages = append(usages, "Encipher Only")
	}
	if usage&x509.KeyUsageDecipherOnly != 0 {
		usages = append(usages, "Decipher Only")
	}

	return usages
}

// parseExtKeyUsage converts ExtKeyUsage slice to string slice
func (c *SSLChecker) parseExtKeyUsage(usage []x509.ExtKeyUsage) []string {
	var usages []string

	for _, u := range usage {
		switch u {
		case x509.ExtKeyUsageServerAuth:
			usages = append(usages, "Server Authentication")
		case x509.ExtKeyUsageClientAuth:
			usages = append(usages, "Client Authentication")
		case x509.ExtKeyUsageCodeSigning:
			usages = append(usages, "Code Signing")
		case x509.ExtKeyUsageEmailProtection:
			usages = append(usages, "Email Protection")
		case x509.ExtKeyUsageTimeStamping:
			usages = append(usages, "Time Stamping")
		case x509.ExtKeyUsageOCSPSigning:
			usages = append(usages, "OCSP Signing")
		default:
			usages = append(usages, "Unknown")
		}
	}

	return usages
}

// calculateFingerprint calculates SHA-256 fingerprint of certificate
func (c *SSLChecker) calculateFingerprint(cert *x509.Certificate) string {
	fingerprint := fmt.Sprintf("%x", cert.Raw)
	// Format as XX:XX:XX...
	var formatted strings.Builder
	for i, char := range fingerprint {
		if i > 0 && i%2 == 0 {
			formatted.WriteString(":")
		}
		formatted.WriteRune(char)
	}
	return strings.ToUpper(formatted.String())
}

// GetExpiryStatus returns a human-readable expiry status
func (check *SSLCheck) GetExpiryStatus() string {
	if check.Error != "" {
		return "Error"
	}

	days := check.DaysUntilExpiry

	if days < 0 {
		return "Expired"
	} else if days == 0 {
		return "Expires Today"
	} else if days == 1 {
		return "Expires Tomorrow"
	} else if days <= 7 {
		return fmt.Sprintf("Expires in %d days", days)
	} else if days <= 30 {
		return fmt.Sprintf("Expires in %d days", days)
	} else {
		return fmt.Sprintf("Valid for %d days", days)
	}
}

// IsExpiringSoon checks if certificate expires within the given days
func (check *SSLCheck) IsExpiringSoon(warningDays int) bool {
	return check.Valid && check.DaysUntilExpiry <= warningDays
}

// IsExpired checks if certificate is already expired
func (check *SSLCheck) IsExpired() bool {
	return check.DaysUntilExpiry < 0
}

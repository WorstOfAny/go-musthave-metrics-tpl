package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(7119),
		Subject: pkix.Name{
			Organization: []string{"Metrics.Organization"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 1),
		SubjectKeyId: []byte{1, 3, 1, 4, 7},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 8192)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate private key")
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate cert")
	}

	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to pem encode cert")
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to pem encode private key")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Msg("failed to get user dir")
	}

	if err = os.WriteFile(filepath.Join(homeDir, "cert.pem"), certPEM.Bytes(), 0644); err != nil {
		log.Error().Err(err).Msg("failed to write cerf to file")
	}

	if err = os.WriteFile(filepath.Join(homeDir, "private.pem"), privateKeyPEM.Bytes(), 0644); err != nil {
		log.Error().Err(err).Msg("failed to write cerf to file")
	}
}

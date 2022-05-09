package tunnel

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
	"time"
	"tunn-hub/common/config"
)

//
// GenerateTlsCert
// @Description:
// @return cert
// @return privateKey
// @return err
//
func GenerateTlsCert(ipAddresses []net.IP, dnsNames []string) (cert []byte, privateKey []byte, err error) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization:       []string{"jackrabbit872568318/tunnel"},
		OrganizationalUnit: []string{"jackrabbit872568318/tunnel"},
		CommonName:         "network tunnel",
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(3650 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  ipAddresses,
		DNSNames:     dnsNames,
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	certBuf := bytes.NewBuffer([]byte{})
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	err = pem.Encode(certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return nil, nil, err
	}
	keyBuf := bytes.NewBuffer([]byte{})
	err = pem.Encode(keyBuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	if err != nil {
		return nil, nil, err
	}
	return certBuf.Bytes(), keyBuf.Bytes(), nil
}

//
// GenerateAndSaveTlsCert
// @Description:
// @param path
//
func GenerateAndSaveTlsCert(path string) error {
	certPath := path + "cert.pem"
	keyPath := path + "key.pem"
	ip := net.ParseIP(config.Current.Global.Address)
	var ipAddresses []net.IP
	var dnsNames []string
	if ip == nil {
		//域名
		dnsNames = append(dnsNames, config.Current.Global.Address)
	} else {
		//IP
		ipAddresses = append(ipAddresses, ip)
	}
	cert, key, err := GenerateTlsCert(ipAddresses, dnsNames)
	if err != nil {
		return err
	}
	certFile, err := os.OpenFile(certPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, err = certFile.Write(cert)
	if err != nil {
		return err
	}
	keyFile, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	_, err = keyFile.Write(key)
	if err != nil {
		return err
	}
	return nil
}

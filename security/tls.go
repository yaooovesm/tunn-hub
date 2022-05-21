package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

//
// NewTunnX509Certification
// @Description:
// @param addrs
// @param names
// @param before
// @return *TunnX509Certification
//
func NewTunnX509Certification(addrs []string, names []string, before time.Time) *TunnX509Certification {
	addresses := make([]net.IP, 0)
	for i := range addrs {
		ip := net.ParseIP(addrs[i])
		if ip != nil {
			addresses = append(addresses, ip)
		}
	}
	return &TunnX509Certification{
		IpAddresses: addresses,
		DNSNames:    names,
		NotAfter:    before,
	}
}

//
// TunnX509Certification
// @Description:
//
type TunnX509Certification struct {
	IpAddresses []net.IP  //allowed ip addresses
	DNSNames    []string  //allowed dns names
	NotAfter    time.Time //expire
}

//
// Create
// @Description:
// @receiver c
// @return cert
// @return privateKey
// @return err
//
func (c *TunnX509Certification) Create() (cert []byte, privateKey []byte, err error) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:       []string{"TunnHub"},
			OrganizationalUnit: []string{"Client"},
			SerialNumber:       serialNumber.String(),
		},
		NotBefore:   time.Now(),
		NotAfter:    c.NotAfter,
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: c.IpAddresses,
		DNSNames:    c.DNSNames,
		Issuer: pkix.Name{
			Organization:       []string{"TunnHub"},
			OrganizationalUnit: []string{"Hub"},
			SerialNumber:       serialNumber.String(),
		},
	}
	return createCertAndKey(template)
}

//
// CreateAndWriteWithRandomName
// @Description:
// @receiver c
// @param path
//
func (c *TunnX509Certification) CreateAndWriteWithRandomName(path string) (string, error) {
	stat, err := os.Stat(path)
	if err != nil || !stat.IsDir() {
		err := os.MkdirAll(path, 0600)
		if err != nil {
			return "", err
		}
	}
	cert, key, err := c.Create()
	if err != nil {
		return "", err
	}
	uid, _ := uuid.NewV4()
	name := strings.ReplaceAll(uid.String(), "-", "")
	certPath := path + name + ".cert"
	keyPath := path + name + ".key"
	err = ioutil.WriteFile(certPath, cert, 0644)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(keyPath, key, 0644)
	if err != nil {
		return "", err
	}
	return name, nil
}

//
// createCertAndKey
// @Description:
// @param template
// @return cert
// @return privateKey
// @return err
//
func createCertAndKey(template x509.Certificate) (cert []byte, privateKey []byte, err error) {
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

package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

func HmacSha1(keyStr, text string) string {
	//hmac ,use sha1
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(text))
	res := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return res
}

func HmacSha256(keyStr, text string)  string {
	key := []byte(keyStr)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(text))
	res := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return res
}

func ReadPEM(pem string) (*x509.CertPool, error) {
	if pem == "" {
		return nil, errors.New("pem not exist")
	}

	certBytes, err := ioutil.ReadFile(pem)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(certBytes)

	if ok {
		return caCertPool, nil
	}

	return nil, errors.New("read pem error")
}

func ReadTls(certFile, keyFile string) ([]tls.Certificate, error){
	if	certFile == "" || keyFile == "" {
		return nil, errors.New("tls not exist")
	}

	tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return []tls.Certificate{tlsCert}, nil
}

func Md5Read(read io.Reader) string {
	md5h := md5.New()
	io.Copy(md5h, read)
	return hex.EncodeToString(md5h.Sum([]byte(nil)))
}

func Md5File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		file.Close()
		return "", err
	}
	md5h := md5.New()
	io.Copy(md5h, file)
	file.Sync()
	file.Close()

	return hex.EncodeToString(md5h.Sum([]byte(nil))), nil
}
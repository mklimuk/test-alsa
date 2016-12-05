package soap

import (
	"io/ioutil"
	"testing"

	"crypto/tls"
	"crypto/x509"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SOAPTestSuite struct {
	suite.Suite
}

func (suite *SOAPTestSuite) TestSepeSOAPCall() {
	cert, err := tls.LoadX509KeyPair("../cert/cert.pem", "../cert/key.pem")
	assert.NoError(suite.T(), err, "Keypair could not be loaded.")
	caCert, err := ioutil.ReadFile("../cert/cacerts.pem")
	assert.NoError(suite.T(), err, "Could not load CACerts.")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	conf := tls.Config{
		ClientAuth:         tls.RequestClientCert,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS12,
	}
	conf.BuildNameToCertificate()

	rozklad := NewIRozkladJazdy("https://sdip.plk-sa.pl/v1.1/RozkladJazdy.svc", &conf, nil)
	req := PobierzRozkladPlanowyTygodniowyDlaStacji{
		IdStacji: 48355,
	}
	resp, err := rozklad.PobierzRozkladPlanowyTygodniowyDlaStacji(&req)
	assert.NoError(suite.T(), err, "Call should not fail.")
	assert.NotNil(suite.T(), resp, "Response should not be empty")
}

func TestSOAPTestSuite(t *testing.T) {
	suite.Run(t, new(SOAPTestSuite))
}

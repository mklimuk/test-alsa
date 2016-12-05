package soap

import (
	"fmt"
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

var req = `POST https://sdip.plk-sa.pl/v1.2/RozkladJazdy.svc HTTP/1.1
	Accept-Encoding: gzip,deflate
	Content-Type: text/xml;charset=UTF-8
	SOAPAction: "http://sdip.plk-sa.pl/v1.2/IRozkladJazdy/PobierzRozkladPlanowy"
	Content-Length: 750
	Host: sdip.plk-sa.pl
	Connection: Keep-Alive
	User-Agent: Apache-HttpClient/4.1.1 (java 1.5)

	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:v1="http://sdip.plk-sa.pl/v1.2">
	   <soapenv:Header/>
	   <soapenv:Body>
	      <v1:PobierzRozkladPlanowy>
	         <!--Optional:-->
	         <v1:kryteria>
	            <v1:ZakresRozkladu>Tygodniowy</v1:ZakresRozkladu>
	            <!--Optional:-->
	            <v1:DataOd>2016-08-02</v1:DataOd>
	            <!--Optional:-->
	            <v1:DataDo>2016-08-06</v1:DataDo>
	            <!--Optional:-->
	            <v1:ListaStacji>
	               <!--Zero or more repetitions:-->
	               <v1:IdentyfikatorStacji>48355</v1:IdentyfikatorStacji>
	            </v1:ListaStacji>
	         </v1:kryteria>
	      </v1:PobierzRozkladPlanowy>
	   </soapenv:Body>
	</soapenv:Envelope>`

func (suite *SOAPTestSuite) TestSepeSOAPCall() {
	cert, err := tls.LoadX509KeyPair("../cert/cert.pem", "../cert/key.pem")
	assert.NoError(suite.T(), err, "Keypair could not be loaded.")
	caCert, err := ioutil.ReadFile("../cert/cacerts.pem")
	assert.NoError(suite.T(), err, "Could not load CACerts.")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	conf := tls.Config{
		ClientAuth:   tls.RequestClientCert,
		Certificates: []tls.Certificate{cert},
		//RootCAs:            caCertPool,
		ClientCAs:          caCertPool,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
	}
	conf.BuildNameToCertificate()

	/*t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &conf,
	}*/

	//client := &http.Client{Transport: t}
	conn, err := tls.Dial("tcp", "sdip.plk-sa.pl:443", &conf)
	assert.NoError(suite.T(), err, "Dial should not fail.")
	client := tls.Client(conn, &conf)
	err = client.Handshake()
	assert.NoError(suite.T(), err, "Handshake should not fail.")
	fmt.Printf("%+v", client.ConnectionState())
	conn.Close()

	//req, err := http.NewRequest("POST", "https://sdip.plk-sa.pl/v1.1/RozkladJazdy.svc", nil)
	//	_, err = client.Do(req)
	//assert.NoError(suite.T(), err, "Call should not fail.")

	/*rozklad := NewIRozkladJazdy("https://sdip.plk-sa.pl/v1.1/RozkladJazdy.svc", &conf, nil)
	req := PobierzRozkladPlanowyTygodniowyDlaStacji{
		IdStacji: 48355,
	}
	resp, err := rozklad.PobierzRozkladPlanowyTygodniowyDlaStacji(&req)
	assert.NoError(suite.T(), err, "Call should not fail.")
	assert.NotNil(suite.T(), resp, "Response should not be empty")*/
}

func TestSOAPTestSuite(t *testing.T) {
	suite.Run(t, new(SOAPTestSuite))
}

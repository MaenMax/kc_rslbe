package test_hawk

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.kaiostech.com/cloud/common/config"
	"git.kaiostech.com/cloud/common/model/oauth2"
	"git.kaiostech.com/cloud/common/test/chttp"
	l4g "git.kaiostech.com/cloud/thirdparty/code.google.com/p/log4go"
)

var (
	_conf         *config.FEConfig
	_ca_cert_pool x509.CertPool
)

const (
	FILE_SIGN = 0x40

	ISRG_CA_ROOT = "conf/root_ca/root.pem" // ISRG Root X1 Certificate of Let's encrypt
	INTER_CA_X1  = "conf/root_ca/x3.pem"   // Primary Let's encrypt intermediate certificate
	INTER_CA_X2  = "conf/root_ca/x4.pem"   // Fallback Let's encrypt intermediate certificate
)

func init() {
	l4g.LoadConfiguration("conf/test-hawk-log.xml")
}

func Do(url_str string, http_verb string, body []byte, token *oauth2.T_AccessToken) (*http.Response, error) {
	var err error
	var data string
	var data_reader io.Reader
	var ext string
	var resp *http.Response
	var http_cli chttp.T_CHTTP
	var _compress bool

	// Allow l4g to flush its content into the log file before terminating.
	defer time.Sleep(time.Second * 1)

	// Making sure to have Debug log output.
	_conf = config.GetFEConfig()
	_conf.Common.Debug = true

	http_cli = *chttp.New()

	//multipart/form-data POST

	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to read file '%s': %s.\n", data[1:], err))
		return nil, err
	}

	data_reader = bytes.NewReader(body)

	if len(http_verb) == 0 {
		http_verb = "GET"
	}

	//Loading HAWK context if specified

	http_cli.SetToken(token)
	http_cli.SetHawkOn(true)

	if err = http_cli.AppendCert(ISRG_CA_ROOT); err != nil {
		fmt.Println(fmt.Sprintf("Failed to append cert: %s. Cause: %s", ISRG_CA_ROOT, err.Error()))
		return nil, err
	}

	if err = http_cli.AppendCert(INTER_CA_X1); err != nil {
		fmt.Println(fmt.Sprintf("Failed to append cert: %s. Cause: %s", INTER_CA_X1, err.Error()))
		return nil, err
	}

	if err = http_cli.AppendCert(INTER_CA_X2); err != nil {
		fmt.Println(fmt.Sprintf("Failed to append cert: %s. Cause: %s", INTER_CA_X2, err.Error()))
		return nil, err
	}
	header_spec := make(map[string]string)
	header_spec["Content-Type"] = "application/json"
	resp, err = http_cli.Request(http_verb, url_str, ext, data_reader, header_spec, _compress)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

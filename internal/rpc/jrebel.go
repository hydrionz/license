package rpc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// jrebelPrivateKeyPEM is the static signing key used by the legacy JRebel rpc
// dispatch path. The active JRebel handlers live in internal/jrebel; these are
// kept for parity with the original RpcController dispatch.
const jrebelPrivateKeyPEM = `
-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBALecq3BwAI4YJZwhJ+snnDFj3lF3DMqNPorV6y5ZKXCiCMqj8OeO
mxk4YZW9aaV9ckl/zlAOI0mpB3pDT+Xlj2sCAwEAAQJAW6/aVD05qbsZHMvZuS2A
a5FpNNj0BDlf38hOtkhDzz/hkYb+EBYLLvldhgsD0OvRNy8yhz7EjaUqLCB0juIN
4QIhAOeCQp+NXxfBmfdG/S+XbRUAdv8iHBl+F6O2wr5fA2jzAiEAywlDfGIl6acn
akPrmJE0IL8qvuO3FtsHBrpkUuOnXakCIQCqdr+XvADI/UThTuQepuErFayJMBSA
sNe3NFsw0cUxAQIgGA5n7ZPfdBi3BdM4VeJWb87WrLlkVxPqeDSbcGrCyMkCIFSs
5JyXvFTreWt7IQjDssrKDRIPmALdNjvfETwlNJyY
-----END RSA PRIVATE KEY-----
`

var (
	jrebelKeyOnce sync.Once
	jrebelKey     *rsa.PrivateKey
)

func getJrebelKey() *rsa.PrivateKey {
	jrebelKeyOnce.Do(func() {
		block, _ := pem.Decode([]byte(jrebelPrivateKeyPEM))
		if block == nil {
			return
		}
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return
		}
		jrebelKey = key
	})
	return jrebelKey
}

func jrebelSign(content string) string {
	key := getJrebelKey()
	if key == nil {
		return ""
	}
	hashed := sha256.Sum256([]byte(content))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(signature)
}

func jrebelPing(ctx *gin.Context, salt string) {
	xmlContent := fmt.Sprintf("<PingResponse><message></message><responseCode>OK</responseCode><salt>%s</salt></PingResponse>", salt)
	ctx.String(http.StatusOK, fmt.Sprintf("<!-- %s -->\n%s", jrebelSign(xmlContent), xmlContent))
}

func jrebelObtainTicket(ctx *gin.Context, username, salt string) {
	prolongationPeriod := "607875500"
	xmlContent := fmt.Sprintf(`<ObtainTicketResponse><message></message><prolongationPeriod>%s</prolongationPeriod><responseCode>OK</responseCode><salt>%s</salt><ticketId>1</ticketId><ticketProperties>licensee=%s\tlicenseType=0\t</ticketProperties></ObtainTicketResponse>`, prolongationPeriod, salt, username)
	ctx.String(http.StatusOK, fmt.Sprintf("<!-- %s -->\n%s", jrebelSign(xmlContent), xmlContent))
}

func jrebelReleaseTicket(ctx *gin.Context, salt string) {
	xmlContent := fmt.Sprintf("<ReleaseTicketResponse><message></message><responseCode>OK</responseCode><salt>%s</salt></ReleaseTicketResponse>", salt)
	ctx.String(http.StatusOK, fmt.Sprintf("<!-- %s -->\n%s", jrebelSign(xmlContent), xmlContent))
}

package provider

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// GenerateSelfSignedCert Generate a self-signed TLS Certificate, to be used as fallback
func GenerateSelfSignedCert() (string, string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	keyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", "", err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"authentik"},
			CommonName:   "authentik default certificate",
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	template.DNSNames = []string{"*"}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatal(err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatal(err)
	}
	privPemByes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	return string(pemBytes), string(privPemByes), nil
}

func TestAccResourceCertificateKeyPair(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	cert, key, err := GenerateSelfSignedCert()
	if err != nil {
		t.Fatal(err)
	}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCertificateKeyPairSimple(rName, cert, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_certificate_key_pair.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_certificate_key_pair.name", "certificate_data", cert),
					resource.TestCheckResourceAttr("authentik_certificate_key_pair.name", "key_data", key),
				),
			},
			{
				Config: testAccResourceCertificateKeyPairSimple(rName+"test", cert, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_certificate_key_pair.name", "name", rName+"test"),
				),
			},
		},
	})
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerTestFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceCertificateKeyPairSimple(rName, cert, key),
				ExpectError: regexp.MustCompile("mock-failed-request"),
			},
			{
				Config:      testAccResourceCertificateKeyPairSimple(rName+"test", cert, key),
				ExpectError: regexp.MustCompile("mock-failed-request"),
			},
		},
	})
}

func testAccResourceCertificateKeyPairSimple(name string, cert string, key string) string {
	return fmt.Sprintf(`
resource "authentik_certificate_key_pair" "name" {
  name              = "%[1]s"
  certificate_data = <<EOT
%[2]sEOT
  key_data = <<EOT
%[3]sEOT
}
`, name, cert, key)
}

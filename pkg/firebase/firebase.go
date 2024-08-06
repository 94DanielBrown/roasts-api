package firebase

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var (
	firebaseCertURL = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
	certCache       = &certs{}
)

type certs struct {
	mu    sync.Mutex
	certs map[string]*rsa.PublicKey
}

func (c *certs) fetchCertificates() error {
	resp, err := http.Get(firebaseCertURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unable to fetch certificates")
	}

	var certMap map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&certMap); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.certs = make(map[string]*rsa.PublicKey)
	for key, cert := range certMap {
		block, _ := pem.Decode([]byte(cert))
		if block == nil {
			return errors.New("failed to parse certificate PEM")
		}

		parsedCert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}

		if parsedCert.PublicKeyAlgorithm != x509.RSA {
			return errors.New("certificate key algorithm is not RSA")
		}

		c.certs[key] = parsedCert.PublicKey.(*rsa.PublicKey)
	}

	return nil
}

func (c *certs) getKey(token *jwt.Token) (interface{}, error) {
	if kid, ok := token.Header["kid"].(string); ok {
		c.mu.Lock()
		defer c.mu.Unlock()
		if pubKey, found := c.certs[kid]; found {
			return pubKey, nil
		}
	}
	return nil, errors.New("unable to find appropriate key")
}

func FirebaseJWTMiddleware() echo.MiddlewareFunc {
	// Fetch certificates at startup
	if err := certCache.fetchCertificates(); err != nil {
		panic(fmt.Sprintf("failed to fetch Firebase certificates: %v", err))
	}

	// Periodically refresh the certificates
	go func() {
		for {
			time.Sleep(24 * time.Hour)
			if err := certCache.fetchCertificates(); err != nil {
				fmt.Printf("failed to refresh Firebase certificates: %v\n", err)
			}
		}
	}()

	return echojwt.WithConfig(echojwt.Config{
		KeyFunc: certCache.getKey,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.MapClaims)
		},
		SuccessHandler: func(c echo.Context) {
			claims := c.Get("user").(*jwt.Token).Claims.(*jwt.MapClaims)
			userID := (*claims)["user_id"].(string)
			c.Set("userID", userID)
			fmt.Println("JWT validated successfully and userID set in context:", userID)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			token := c.Request().Header.Get("Authorization")
			fmt.Printf("Invalid JWT: %s\n", token)
			fmt.Printf("Error: %v\n", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
		},
	})
}

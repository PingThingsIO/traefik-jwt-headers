package traefik_jwt_headers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	ValueRewrite      map[string]map[string]string `json:"valueRewrite,omitempty"`
	Headers           map[string]string            `json:"headers,omitempty"`
	ClaimsPrefix      string                       `json:"claimsPrefix"`
	UnboxFirstElement bool                         `json:"unboxFirstElement"`
}

func CreateConfig() *Config {
	return &Config{
		Headers:           make(map[string]string),
		UnboxFirstElement: true,
	}
}

type JwtHeaders struct {
	next http.Handler
	name string

	headers           map[string]string
	claimsPrefix      string
	unboxFirstElement bool
	valueRewrite      map[string]map[string]string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &JwtHeaders{
		next:              next,
		name:              name,
		headers:           config.Headers,
		claimsPrefix:      config.ClaimsPrefix,
		unboxFirstElement: config.UnboxFirstElement,
		valueRewrite:      config.ValueRewrite,
	}, nil
}

func (a *JwtHeaders) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("token")
	if err != nil {
		log.Printf("ERROR DECODING COOKIE 'token': %+v", err)
		a.next.ServeHTTP(rw, req)
		return
	}

	body := c.Value

	p := jwt.NewParser()
	token, _, err := p.ParseUnverified(body, jwt.MapClaims{})
	if err != nil {
		log.Printf("CLAIMS ERROR: %+v", err)
	}

	log.Printf("DECODED TOKEN CLAIMS AS: %+v", token.Claims)

	if m, ok := token.Claims.(jwt.MapClaims); ok {
		var claims map[string]interface{} = m

		if a.claimsPrefix != "" {
			if c, ok := claims[a.claimsPrefix].(map[string]interface{}); ok {
				claims = c
			}
		}

		a.setHeaders(claims, req)
	}

	a.next.ServeHTTP(rw, req)
}

func (a *JwtHeaders) setHeaders(claims map[string]interface{}, req *http.Request) {
	if len(a.headers) == 0 {
		return
	}

	for key, value := range a.headers {
		if v, ok := claims[key]; ok {
			if a.unboxFirstElement {
				switch s := v.(type) {
				case []interface{}:
					if len(s) > 0 {
						v = s[0]
					}
				}
			}

			v := fmt.Sprintf("%v", v)

			if rewrite, ok := a.valueRewrite[key]; ok {
				if newValue, ok := rewrite[v]; ok {
					v = newValue
				}
			}

			log.Printf(`Setting Header "%s" to "%s"`, value, v)

			req.Header.Set(value, v)
		}
	}
}

package main

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

// JWK is json data struct for JSON Web Key
type JWK struct {
	Keys []JWKKey
}
type JWKKey struct {
	Alg string
	E   string
	Kid string
	Kty string
	N   string
	Use string
}

func getJWK(jwkURL string) map[string]JWKKey {

	jwk := &JWK{}

	getJSON(jwkURL, jwk)

	jwkMap := make(map[string]JWKKey, 0)
	for _, jwk := range jwk.Keys {
		jwkMap[jwk.Kid] = jwk
	}
	return jwkMap
}

func getJSON(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func convertKey(rawE, rawN string) *rsa.PublicKey {
	decodedE, err := base64.RawURLEncoding.DecodeString(rawE)
	if err != nil {
		panic(err)
	}
	if len(decodedE) < 4 {
		ndata := make([]byte, 4)
		copy(ndata[4-len(decodedE):], decodedE)
		decodedE = ndata
	}
	pubKey := &rsa.PublicKey{
		N: &big.Int{},
		E: int(binary.BigEndian.Uint32(decodedE[:])),
	}
	decodedN, err := base64.RawURLEncoding.DecodeString(rawN)
	if err != nil {
		panic(err)
	}
	pubKey.N.SetBytes(decodedN)
	// fmt.Println(decodedN)
	// fmt.Println(decodedE)
	// fmt.Printf("%#v\n", *pubKey)
	return pubKey
}

func VerifyJWTToken(tokenStr string, jwkURL string) (*jwt.Token, error) {
	jwk := getJWK(jwkURL)
	//log.Printf("Public jwks: %v",jwk)

	token, err := validateToken(tokenStr, jwk)
	if err != nil || !token.Valid {
		log.Println("Verify token err:", err)
		//fmt.Printf("token is not valid %v\n", err)
		return token, err
	}
	return token, nil
}
func validateToken(tokenStr string, jwk map[string]JWKKey) (*jwt.Token, error) {
	// 2. Decode the token string into JWT format.
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// signing method : RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// 5. Get the kid from the JWT token header and retrieve the corresponding JSON Web Key that was stored
		if kid, ok := token.Header["kid"]; ok {
			if kidStr, ok := kid.(string); ok {
				key := jwk[kidStr]
				// 6. Verify the signature of the decoded JWT token.
				rsaPublicKey := convertKey(key.E, key.N)
				return rsaPublicKey, nil
			}
		}

		return "", nil
	})
	if err != nil {
		return token, err
	}

	//claims := token.Claims.(jwt.MapClaims)

	/*
		iss, ok := claims["iss"]
		if !ok {
			return token, fmt.Errorf("token does not contain issuer")
		}
			issStr := iss.(string)
			if strings.Contains(issStr, "cognito-idp") {
				err = validateAWSJwtClaims(claims, region, userPoolID)
				if err != nil {
					return token, err
				}
			}
	*/
	if token.Valid {
		return token, nil
	}
	return token, err
}

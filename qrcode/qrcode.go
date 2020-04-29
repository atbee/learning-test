package qrcode

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/codeforpublic/morchana-static-qr-code-api/internal/jsonw"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	qrcode "github.com/skip2/go-qrcode"
	"gopkg.in/square/go-jose.v2/json"
)

type Data struct {
	AnonymousID string `json:"anonymousId"`
	Code        string `json:"code"`
}

type Qr struct {
	Type   string `json:"type"`
	Base64 string `json:"base64"`
}

func Generate(signature string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data Data

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			jsonw.BadRequest(w, err)
			return
		}

		// s, err := token(signature, "", data.AnonymousID, "23d", "THCOVID", time.Now())
		// if err != nil {
		// 	jsonw.InternalServerError(w, err)
		// 	return
		// }

		png, err := qrcode.Encode(data.AnonymousID, qrcode.Medium, 256)
		if err != nil {
			jsonw.InternalServerError(w, err)
			return
		}

		json.NewEncoder(w).Encode(&Qr{
			Type:   "image/png",
			Base64: base64.StdEncoding.EncodeToString(png),
		})
	}
}

func token(signature, aid, code, tagID, life, iss string, now time.Time) (string, error) {
	signingKey := []byte(signature)

	nonce, err := generateRandomBytes(10)
	if err != nil {
		return "", errors.Wrap(err, "generate nonce")
	}
	// Create the Claims
	claims := QRFormat{
		Empty: []string{
			aid,
			code,
			tagID,
			life,
			string(nonce),
		},
		StandardClaims: jwt.StandardClaims{
			IssuedAt: now.Unix(),
			Issuer:   iss,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

type QRFormat struct {
	Empty []string `json:"_"` // anonymousId,code,tagId,age,nounce
	jwt.StandardClaims
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

// {
// 	"_": [
// 	  "zB84S_2Ovt5Z",
// 	  "G",
// 	  "",
// 	  "23d",
// 	  "Po7kEB48Zx"
// 	],
// 	"iat": 1587997998,
// 	"iss": "THCOVID"
//   }

// { _: [anonymousId: string, code: string, tagId: string, age: string, nounce: string], iat: string }

package login

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/codeforpublic/morchana-static-qr-code-api/internal/jsonw"
	"github.com/gorilla/mux"
)

func LoginOTP(client *http.Client, baseUrl *url.URL, params url.Values) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// uri := "https://ohgikdu5ed.execute-api.ap-southeast-1.amazonaws.com/smsgw-api"

		// baseUrl, err := url.Parse(uri)
		// if err != nil {
		// 	jsonw.InternalServerError(w, err)
		// 	return
		// }

		// params := url.Values{}
		// params.Add("user", "morchana2")
		// params.Add("pass", "Y8adfQzJfwUKGwUY")
		// params.Add("from", "Morchana")
		params.Add("to", vars["subr"])
		// params.Add("msg", "ทดสอบ ข้อความ ภาษาไทย")

		baseUrl.RawQuery = params.Encode()

		req, err := http.NewRequest(http.MethodPost, baseUrl.String(), nil)
		if err != nil {
			jsonw.InternalServerError(w, err)
			return
		}

		_, err = client.Do(req)
		if err != nil {
			jsonw.InternalServerError(w, err)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}
}

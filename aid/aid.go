package aid

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/codeforpublic/morchana-static-qr-code-api/internal/jsonw"
	"github.com/google/uuid"
)

type Identity struct {
	DeviceID string `json:"deviceId"`
}

type Anonymous struct {
	Status      string `json:"status"`
	AnonymousID string `json:"anonymousId"`
}

type storeAnonymousIDFunc func(context.Context, time.Time, string, string) error

func AnonymousID(store storeAnonymousIDFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id Identity

		if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
			jsonw.BadRequest(w, err)
			return
		}

		aid := uuid.New().String()
		now := time.Now()

		if err := store(r.Context(), now, id.DeviceID, aid); err != nil {
			jsonw.InternalServerError(w, err)
			return
		}

		json.NewEncoder(w).Encode(&Anonymous{
			Status:      "200",
			AnonymousID: aid,
		})
	}
}

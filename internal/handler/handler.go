package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/pkg/models"
	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/processor"
)

func CreateOrderHandler(w http.ResponseWriter, r *http.Request, pool *processor.Pool) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var o models.Order
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&o); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	o.CreatedAt = time.Now().Unix()

	if err := o.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o.ID = generateID()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(o)

	pool.Orders <- o

}

func generateID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

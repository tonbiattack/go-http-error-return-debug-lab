package debuglab

import (
	"encoding/json"
	"net/http"
)

type createOrderRequest struct {
	ProductID string `json:"productId"`
}

type createOrderResponse struct {
	ID string `json:"id"`
}

// CreateOrderHandler は注文作成の最小HTTPハンドラーです。
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "リクエスト本文が不正です", http.StatusBadRequest)
		return
	}

	if request.ProductID == "" {
		http.Error(w, "productId は必須です", http.StatusBadRequest)
		// 意図的な不具合: エラー応答後に処理を終了していない。
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createOrderResponse{ID: "order-001"})
}

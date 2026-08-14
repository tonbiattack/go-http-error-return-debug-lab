package debuglab

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateOrderHandler_商品IDが空ならエラー本文だけを返す(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"productId":""}`))
	responseRecorder := httptest.NewRecorder()

	CreateOrderHandler(responseRecorder, request)

	response := responseRecorder.Result()
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("レスポンス本文の読み取りに失敗しました: %v", err)
	}

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusBadRequest)
	}

	const wantBody = "productId は必須です\n"
	if string(body) != wantBody {
		t.Fatalf("body = %q, want %q", body, wantBody)
	}
}

func TestCreateOrderHandler_商品IDがあれば注文IDを返す(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"productId":"product-123"}`))
	responseRecorder := httptest.NewRecorder()

	CreateOrderHandler(responseRecorder, request)

	response := responseRecorder.Result()
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("レスポンス本文の読み取りに失敗しました: %v", err)
	}

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusCreated)
	}

	const wantContentType = "application/json"
	if response.Header.Get("Content-Type") != wantContentType {
		t.Fatalf("Content-Type = %q, want %q", response.Header.Get("Content-Type"), wantContentType)
	}

	const wantBody = "{\"id\":\"order-001\"}\n"
	if string(body) != wantBody {
		t.Fatalf("body = %q, want %q", body, wantBody)
	}
}

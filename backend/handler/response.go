package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// NOTE: レスポンスデータをJSONに変換してステータスコードと一緒にhttp.ResponseWriterインターフェースを満たす型の値に書き込む
//       HTTPハンドラーの共通処理を，ヘルパー関数として実装
func RespondJSON(ctx context.Context, w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// https://ja.foobarninja.com/golang/standard-library/encoding-json/marshal/#:~:text=Marshal%E9%96%A2%E6%95%B0%E3%81%AFjson%E3%83%91%E3%83%83%E3%82%B1%E3%83%BC%E3%82%B8,error%E5%9E%8B%E3%81%AE%E5%80%A4%E3%81%A7%E3%81%99%E3%80%82
	// json.Marshal(): 引数として与えた値をJSONに変換（エンコーディング）して返り値として返す
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rsp := ErrResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		// https://medium-company.com/encode/
		// データを別の形式に変換する（データ圧縮や暗号化も含む）ことをエンコード（encode），エンコードしたデータを元に戻すことをデコード（decode）
		// https://zenn.dev/hsaki/articles/go-convert-json-struct
		// エンコーディングはインメモリ表現からバイト列表現への変換のこと．Go構造体からjsonを生成 and 平文から暗号文を生成
		// エンコーディングと同じ意味の言葉として、シリアライゼーション(serialization)・マーシャリング(marshalling)
		// https://zenn.dev/hsaki/articles/go-convert-json-struct#encoder%E3%81%AE%E5%88%A9%E7%94%A8
		// json.NewEncoder関数から、引数で指定された場所にjsonを出力するエンコーダーを作成
		// Encode(Go構造体)メソッドを実行してエンコード
		// Marshal関数とエンコーダーの違いとしては、前者はエンコード結果が[]byteになるのに対し後者はio.Writerの形で自由に指定することができるという点
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("write error respponse error: %v\n", err)
		}
		return
	}

	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write error respponse error: %v\n", err)
	}
}

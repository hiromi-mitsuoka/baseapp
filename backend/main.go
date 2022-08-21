package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	err := http.ListenAndServe(
		":8000", // localhostを省略した記法
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// https://qiita.com/taji-taji/items/77845ef744da7c88a6fe#%E6%8E%A5%E9%A0%AD%E8%BE%9E-f
			// 接頭辞F : 書き込み先を明示的に指定
			// fmt.Println(r.URL)
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	)
	if err != nil {
		// https://qiita.com/taji-taji/items/77845ef744da7c88a6fe#%E6%8E%A5%E5%B0%BE%E8%BE%9E-f
		// 接尾辞f : フォーマットを指定
		fmt.Printf("failed to terminate server: %v", err)
		// https://qiita.com/umisama/items/7be04949d670d8cdb99c
		// https://tech-up.hatenablog.com/entry/2018/12/13/154143
		// OSに値を返してプロセスを切る．deferも実行されない．
		// 0 = 正常終了．1 = 以上終了
		os.Exit(1)
	}
}

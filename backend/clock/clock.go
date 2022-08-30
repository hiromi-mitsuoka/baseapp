package clock

import "time"

// NOTE: Goのtime.Time型はナノ秒単位の時刻制度の情報を持っているが，永続化したデータを比較すると，ほぼ確実に時刻情報が不一致となる

// インターフェース経由のみでメソッドを実行するように
type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

// NOTE: テスト用の固定時刻を返すインターフェース型
type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	return time.Date(2022, 1, 1, 12, 00, 00, 0, time.UTC)
}

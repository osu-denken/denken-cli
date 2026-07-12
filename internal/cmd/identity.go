package cmd

import "strings"

// studentDomain は大産大の学生メールドメイン。
const studentDomain = "@ge.osaka-sandai.ac.jp"

// resolveEmail は入力をログイン用メールアドレスに正規化する。
// '@' を含めばメールとしてそのまま扱い、含まなければ学籍番号とみなして
// 小文字化・接頭辞 s の補完・ドメイン付与を行う (例: 24H034 -> s24h034@ge.osaka-sandai.ac.jp)。
func resolveEmail(input string) string {
	v := strings.TrimSpace(input)
	if v == "" || strings.Contains(v, "@") {
		return v
	}
	v = strings.ToLower(v)
	if !strings.HasPrefix(v, "s") {
		v = "s" + v
	}
	return v + studentDomain
}

# 開発ルール

## 構造ルール
- 1ファイル300行以内。超過見込みなら実装前に分割案を提示して停止
- 1メソッド/関数40行以内、ネスト3段以内
- 1ファイル1責務。「Utils」「Helper」への追記より新規クラスを優先提案

## 例外処理
try/catchの許可箇所:
- 外部I/O境界: ネットワーク、ファイル、DB、外部API
- エントリポイント最上位（main, リクエストハンドラ）
- 明示的に回復処理が存在する場所（リトライ、フォールバック）

それ以外は例外を伝播させる。

```text
// GOOD: 境界でcatchし、文脈を付けて再送出
if (a == b) return;
try {
	return httpClient.send(req);
} catch (e) {
	throw new Error("API fetch failed: " + url, e);
}

// BAD: ロジック内で握りつぶし
try {
	if (a == b) return; // try内で分岐する必要のないif文
	process(data);
} catch (e) {
	console.error(e); // 呼び出し元は失敗を知れない
}
```

## コーディング規約
- プロジェクト固有の命名、フォーマット規約は今後の実装で明文化する
- コメントは「なぜ」のみ

# tdd-moshi-share
テスト駆動開発の実践：最終コードの公開用

## このリポジトリはなに
テスト駆動開発の実践訓練をした時の最終成果物公開リポジトリ

## 題材
こちらを解きました

https://www.yumemi.co.jp/serverside_recruit

## 途中過程
テストを書きながら少しずつ前進するテスト駆動開発の実践ということで、以下に過程を残しています。参考になれば幸いです。

https://zenn.dev/foxtail88/scraps/17e94c540e0771

## 実行方法
### ビルドして実行
for linux, unix(max)
```shell
go build main.go
./main ./testdata/large.csv
```

windowsの場合は出力ファイル名をexeにして、パスも適切に与え直してあげてください。

### テスト
`main`のテストはないです。全て`ranking`パッケージのテストです。

```shell
go test ./ranking
```

詳細が見たい場合は`-v`オプションをつけてあげると良いです。
# terraform-provider-nifcloud
* [nifcloud-sdk-go][1]を使用させていただき、ニフクラ用の terraform provider を作成してみる。
* [こちら][2]を参考にさせていただく。
* モジュール管理は [Go Modules][4] に変更。

## 環境
```
go version go1.13.1 linux/amd64
Terraform v0.12.9
```

## 作業環境準備
```
cd $GOPATH/src/github.com/
mkdir -p shztki/terraform-provider-nifcloud
cd shztki/terraform-provider-nifcloud/
go mod init
git init
```

## 利用方法
`$GOPATH/src/github.com/shztki/terraform-provider-nifcloud` へリポジトリをクローンして、 `go build` を実施しバイナリファイルを作成してください。

```
$ mkdir -p $GOPATH/src/github.com/shztki; cd $GOPATH/src/github.com/shztki
$ git clone https://github.com/shztki/terraform-provider-nifcloud.git; cd terraform-provider-nifcloud
$ go build
```

ビルドしたバイナリファイルを `~/.terraform.d/plugins/` に配置してインストール完了です。

```
$ mkdir -p ~/.terraform.d/plugins/
$ mv terraform-provider-nifcloud ~/.terraform.d/plugins/
```

## 作成コメント
##### nifcloud/resources
1. ニフクラには「月額/従量」の課金タイプがある。設定する際は `AccountingType` にパラメータを渡すだけだが、変更は翌月からとなる関係で、最新の状態は `NextMonthAccountingType` となる。このため `accounting_type` として tfstate に残す値は `NextMonthAccountingType` にした方がよい。
1. インスタンスで `IpType` を残していると、 `NetworkInterfaces` で作成した場合にも Describe したあとに値が入ってしまい、tfstate の差分が生まれるため、無しとした。 `NetworkInterfaces` の指定で全パターン作成可能(共通グローバル/共通プライベート、共通グローバル/プライベートLAN、共通プライベートのみ、プライベートLANのみ)
1. プライベートLAN に所属させるインスタンスで、 `userdata` を利用してプライベートIPアドレスを設定しない場合、サーバー作成完了までにかなり時間がかかる(サーバーのステータスは「異常あり」で完了)。この場合、サーバー自体は作成されても、terraform の実行はタイムアウトでエラー終了することがある(あとで import は可能)。気に入らない場合は、 `Create: schema.DefaultTimeout(15 * time.Minute)` をもっと延ばしてもいいかもしれない。
1. ファイアウォールグループについて、ログ取得件数を変更する処理は実装していません。必要なら `UpdateSecurityGroup` API を使う形で別途実装が必要。
1. バックアップルールの初回作成時には、最初のバックアップ処理も走ります。完了までに時間がかかるため、 status が available になるまで待つ処理は入れていません。
1. OSイメージの作成完了までは時間がかかるため、 State が available になるまで待つ処理は入れていません。

##### aws-sdk-go
1. いくつか修正しないといけない点があったため、ブランチを分けて[自分の環境に作成][3]。
1. 上記を利用する形とする。修正点は上記コミット履歴を参照。

##### examples/tffiles
1. terraform 0.12 で動作確認中...

## 作成状況
| リソース | ステータス |
|---|---|
| SSHキーインポート | ok |
| プライベートLAN | ok |
| サーバー | ok |
| ディスク | ok |
| ファイアウォール | ok |
| ファイアウォールグループルール | ok |
| バックアップ | ok |
| OSイメージ | ok |
| 拠点間VPNゲートウェイ | 作成中... |
| ロードバランサ | 検討中... |
| マルチロードバランサー | 検討中... |
| 付替IPアドレス | 検討中... |
| 追加NIC | 検討中... |
| 基本監視 | 検討中... |
| ルーター | 検討中... |
| サーバーセパレート | 検討中... |
| RDB | 検討中... |
| NAS | 検討中... |


[1]:https://github.com/alice02/nifcloud-sdk-go
[2]:https://github.com/kzmake/terraform-provider-nifcloud
[3]:https://github.com/shztki/nifcloud-sdk-go
[4]:https://blog.golang.org/using-go-modules
# terraform-provider-nifcloud
* [nifcloud-sdk-go][1]を使用させていただき、ニフクラ用の terraform provider を作成してみる。
* [こちら][2]を参考にさせていただく。

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
dep init
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
1. インスタンスで `IpType` を残していると、 `NetworkInterfaces` で作成した場合にも Describe したあとに値が入ってしまい、tfstate の差分が生まれるため、無しとした。
1. プライベートLAN に所属させるインスタンスで、 `userdata` を利用してプライベートIPアドレスを設定しない場合、サーバー作成完了までにかなり時間がかかる(サーバーのステータスは「異常あり」で完了)。この場合、サーバー自体は作成されても、terraform の実行はタイムアウトでエラー終了することがある(あとで import は可能)。気に入らない場合は、 `Create: schema.DefaultTimeout(5 * time.Minute)` を 15分くらいに延ばした方がいいかもしれない。
1. ファイアウォールグループについて、ログ取得件数を変更する処理は実装していません。必要なら `UpdateSecurityGroup` を使う形で別途実装が必要。

##### aws-sdk-go
1. いくつか修正しないといけない点があったため、ブランチを分けて[自分の環境に作成][3]。
1. 上記を利用する形とする。修正点は上記コミット履歴を参照。

##### examples/tffiles
1. terraform 0.12 で動作確認中...

[1]:https://github.com/alice02/nifcloud-sdk-go
[2]:https://github.com/shztki/terraform-provider-nifcloud
[3]:https://github.com/shztki/nifcloud-sdk-go

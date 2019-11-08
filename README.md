# Terraform Provider for NIFCLOUD
* [nifcloud-sdk-go][1]を使用させていただき、ニフクラ用の terraform provider を作成してみる。
* [こちら][2]を参考にさせていただく。

## 環境
* [Terraform][5] 0.12+
* [Go][4] 1.13 (to build the provider plugin)

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
1. ファイアウォールグループルールの追加について、かなり時間がかかることがあるようで、追加されないままタイムアウトして終了することもあります。ただ、タイムアウト時間を延ばしたり、再作成処理を実施したりするのもあまり意味が無さそうだったので、対応していません。
1. バックアップルールの初回作成時には、最初のバックアップ処理も走ります。完了までに時間がかかるため、 status が available になるまで待つ処理は入れていません。
1. OSイメージの作成完了までは時間がかかるため、 State が available になるまで待つ処理は入れていません。
1. `nifcloud_vpn_connection` について、 `NiftyIpsecConfiguration` 情報を `ipsec` へ、 `NiftyTunnel` 情報を `tunnel` へ戻すようにしている。PreSharedKey が自動生成可能なので、情報を取ってくる意味も含めて実装してある。ただしこのせいで、自動生成された部分との差分が検出され、Terraform側で毎回再作成処理が走ってしまうため、 ignore_changes の指定が必須となる。
	* VPNコネクションは一切の変更が不可なリソースなので、厳密な管理は不要と思われるので、情報が不要であれば `resourceNifcloudVpnConnectionRead` から該当の処理をコメントアウトしても可。
	* `mtu` を `tunnel` に含めるかどうか、非常に悩んだ末に含めた。変更可能なのは `NiftyTunnel` を使うときのみなのに、パラメータとしては `NiftyIpsecConfiguration` に含まれる。。。
	* 一切変更不可のリソースなので、全体的に ForceNew にしたかったが、できていない。 Default を指定したい部分もあり、そうすると ForceNew できない。これも悩ましい。。。Update処理は実装していないので、変更を検知したとしても、何も処理は実施されない。
1. VPNゲートウェイについて、RouteTable関連の処理は実装していません。
1. RDBは「新規作成、スナップショットからの作成、リードレプリカとしての作成」の 3パターンが可能です。
	* ニフクラ独自仕様で、MySQLにのみ冗長化に `性能優先` というのがある。これを選ぶと、フェールオーバー可能なリードレプリカが追加でできあがる。作成したリソースは 1つなのに、実際には 2個の RDB が存在することになる。おかしな感じだが、とりあえずそのままに。操作はできないが、このリードレプリカがいると削除できなくなってしまうので、 `replica_identifier` がある場合はそれを先に削除するようにしている。
	* 「スナップショットからの作成、リードレプリカとしての作成」時に、初回作成時に指定はできなくても、変更が可能なパラメータについては、反映できるようにしてあります(パラメータグループの変更時には再起動も実行)。
	* 原因がよくわかりませんでしたが、「スナップショットからの作成」時に `InternalFailure: System Error.` や `SerializationError: failed decoding Query response` で異常終了するものの、RDB自体は無事作成される、ということがあったため、これらのエラー時は無視して継続するようにしてあります。

##### aws-sdk-go
1. いくつか修正しないといけない点があったため、ブランチを分けて[自分の環境に作成][3]。
1. 上記を利用する形とする。修正点は上記コミット履歴を参照。
1. 基本は `models/apis/` 配下の jsonファイルを修正し、 `go generate ./service` して、タグのバージョンを 1つずつ上げている。そのたび、こちらの `go.mod` のバージョン指定を修正。

##### examples/tffiles
1. terraform 0.12 で動作確認中...

## 作成状況
| リソース | ステータス | 備考 |
|---|---|---|
| SSHキーインポート | ok | Createは作っていません |
| プライベートLAN | ok | |
| サーバー | ok | インポートやコピーは作っていません |
| ディスク | ok | |
| ファイアウォール | ok | ログ取得件数変更処理は作っていません |
| バックアップ | ok | |
| OSイメージ | ok | |
| 拠点間VPNゲートウェイ | ok | ルートテーブルは作っていません |
| RDB | ok | イベント通知は作っていません |
| ロードバランサー | n/a | |
| マルチロードバランサー | n/a | |
| 付替IPアドレス | n/a | |
| 追加NIC | n/a | |
| 基本監視 | n/a | |
| ルーター | n/a | |
| サーバーセパレート | n/a | |
| NAS | n/a | |


[1]:https://github.com/alice02/nifcloud-sdk-go
[2]:https://github.com/kzmake/terraform-provider-nifcloud
[3]:https://github.com/shztki/nifcloud-sdk-go
[4]:https://golang.org/doc/install
[5]:https://www.terraform.io/downloads.html
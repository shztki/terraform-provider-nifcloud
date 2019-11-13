# Terraform Provider for NIFCLOUD examples
nifcloud に各種リソースを作成するためのサンプルコード

### 環境変数
```
AWS_ACCESS_KEY_ID=`ニフクラアカウントのアクセスキー`
AWS_SECRET_ACCESS_KEY=`ニフクラアカウントのシークレットアクセスキー`
TF_VAR_def_pass=`WindowsServerのパスワード`
TF_VAR_ssh_pubkey_path=`SSH公開鍵のファイルパス`
TF_VAR_allow_cidr_001=`FWで許可したい許可したいアドレス`
TF_VAR_pre_shared_key_001=`IPSecの事前共有鍵`
```

### 実行
```
terraform init
terraform plan
terraform apply
```

### サンプルイメージ
こんな構成ができあがります。

![examples_001](https://raw.githubusercontent.com/shztki/terraform-provider-nifcloud/images/nifcloud_examples_001.png)

* イメージとしては、このあと手動でリモートアクセスVPNGW を `192.168.2.245` で作成し、その際の配布NW が `10.168.201.0/24` になるイメージで `example_server_kanri` に userdata でルーティングを設定したり、 `example_firewallgroup_006` にアクセス許可ポリシーを入れたりしています。
* `.disable` にしたりコメントアウトしたりしていますが、RDB やバックアップ、カスタマイズイメージの作成も可能です。


### コメント
* イメージID(スタンダード)については、 [CLI][1] を入れて以下コマンド実行するなどして特定する。
	* 環境変数に `AWS_DEFAULT_REGION` も指定必要。

```
nifcloud-debugcli computing describe-images --query 'ImagesSet[?ImageOwnerId==`niftycloud`].[Name,ImageOwnerId,ImageId]'
```

* 月末に近い場合は、従量で作成後、すぐに月額に一括変更可能(charge_type)。
* 以下のようなコマンドで SSH鍵は生成し、 `test-ssh-key.rsa.pub` の方をインポートするイメージです。
```
ssh-keygen -t rsa -C "" -f test-ssh-key.rsa -N ""
```
* userdata の処理により、プライベートLAN の IPアドレスが設定された状態で作成可能(スクリプトの中身は都度修正必要)。
	* グローバルIP無しにする場合は、userdata でプライベートアドレスを設定する対象NIC が 1個目になる。
	* 処理は入れていないが、コントロールパネルのコンソールからログインできるように、アカウントにはパスワードを設定しておくのがよい。
* ディスクはサーバー作成後に作成・アタッチされるため、自動のマウント処理がしたい場合は別途 Ansible等での対応が必要。
* ファイアウォールルールは、ニフクラ仕様により同時指定できないパラメータや、設定できない値(/32はつけてはダメ)などがあるので注意。
* `nifcloud_securitygroup` でも `rules` の指定でポリシー作成が可能だが、変更するとすべて削除、改めて全体を新規作成、という仕様になります。このため `nifcloud_securitygroup` ではグループを作成するのみにとどめ、 `nifcloud_securitygroup_rule` にて別途ルールをアタッチしていくやり方を推奨。
	* ファイアウォールルールの追加には結構[時間がかかる][2]ことがあるようです。タイムアウトしたら、再度 apply してください。
* `nifcloud_instancebackup_rule` でディスクが増設されたサーバーを対象にしたい場合は、 `depends_on` でボリュームが作成されるまで待つようにすること。でないとタイミングによってはボリューム作成より前にバックアップルールが作成されてしまい、ボリュームが作成できなくなります。
* バックアップルールの変更については、なんらかの処理中だとエラーが返されます。初回設定時は同時にバックアップも走るため、ステータスが available にならない限りは変更できないので注意。
* カスタマイズイメージの変更については、作成中だとエラーが返されます。ステータスが available にならない限りは変更できないので注意。
* VPNコネクションは一切変更不可のリソースのため、Update処理は実装されていない。しかし、完全な ForceNew にはできていないため、 tunnel や ipsec に変更があると、変更処理をしようとしてしまうので、 ignore_changes が必須です。
	* Describeした情報をtfstateに戻している関係で、作成した瞬間から差分ありの状態になります。
* RDBに関して厳密なパラメータのチェックは行っていないため、指定不可な組み合わせにした場合、異常が発生することがあるので、注意してください。
	* プライベートLAN利用時
		* master と virtual の両IPアドレス指定が必須です(かつ /** もつけること)。
		* 冗長化(データ優先)を選択した場合、 slave の IPアドレスも指定が必要です。
		* 冗長化(性能優先)は MySQL時のみ選択可能ですが、この場合は slave の指定は不可で、 replica の名前と IPアドレス指定が必須です。
	* サンプルコード上ではいくつか `ignore_changes` を指定していますが、これには変更不可のものもあれば、特に指定しなかったために State の値と差分が生まれてしまったので抑制しているだけのものがあります。「新規作成、スナップショットからの作成、リードレプリカとしての作成」のそれぞれで、変更不可なものは異なるので、ご注意ください。
* VPN Gateway および Router について、同じプライベートLAN に所属するなど、作成処理が競合する場合はエラーになった。このため、必ずどちらかに `depends_on` を入れること。
	* AssociateRouteTable系にも `depends_on` は入れておくこと。
	* VpnConnection にも `depends_on` は入れておくこと。
	* ネットワーク系の処理が同時実行されないような注意が必要です。

[1]:https://github.com/nifcloud/nifcloud-sdk-python
[2]:https://pfs.nifcloud.com/api/rest/AuthorizeSecurityGroupIngress.htm
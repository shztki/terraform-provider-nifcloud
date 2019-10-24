# terraform-provider-nifcloud
nifcloud に各種リソースを作成するためのサンプルコード

### 環境変数
```
AWS_ACCESS_KEY_ID=`ニフクラアカウントのアクセスキー`
AWS_SECRET_ACCESS_KEY=`ニフクラアカウントのシークレットアクセスキー`
TF_VAR_def_pass=`WindowsServerのパスワード`
TF_VAR_ssh_pubkey_path=`SSH公開鍵のファイルパス`
```

### コメント
* イメージID(スタンダード)については、 [CLI][1] を入れて以下コマンド実行するなどして特定する。
	* 環境変数に `AWS_DEFAULT_REGION` も指定必要。

```
nifcloud-debugcli computing describe-images --query 'ImagesSet[?ImageOwnerId==`niftycloud`].[Name,ImageOwnerId,ImageId]'
```

* 月末に近い場合は、従量で作成後、すぐに月額に一括変更可能(charge_type)。
* userdata の処理により、プライベートLAN の IPアドレスが設定された状態で作成可能(スクリプトの中身は都度修正必要)。
	* グローバルIP無しにする場合は、userdata でプライベートアドレスを設定する対象NIC が 1個目になる。
* ディスクはサーバー作成後に作成・アタッチされるため、自動のマウント処理がしたい場合は別途 Ansible等での対応が必要。
* ファイアウォールルールは、ニフクラ仕様により同時指定できないパラメータや、設定できない値(/32の指定)などがあるので注意。
* `nifcloud_securitygroup` でも `rules` の指定でポリシー作成が可能だが、変更するとすべて削除、改めて全体を新規作成、という仕様になります。このため `nifcloud_securitygroup` ではグループを作成するのみにとどめ、 `nifcloud_securitygroup_rule` にて別途ルールをアタッチしていくやり方を推奨。
	* ファイアウォールルールの追加には結構[時間がかかる][2]ことがあるようです。タイムアウトしたら、再度 apply してください。
* `nifcloud_instancebackup_rule` でディスクが増設されたサーバーを対象にしたい場合は、 `depends_on` でボリュームが作成されるまで待つようにすること。でないとタイミングによってはボリューム作成より前にバックアップルールが作成されてしまい、ボリュームが作成できなくなります。
* バックアップルールの変更については、なんらかの処理中だとエラーが返されます。初回設定時は同時にバックアップも走るため、ステータスが available にならない限りは変更できないので注意。
* カスタマイズイメージの変更については、作成中だとエラーが返されます。ステータスが available にならない限りは変更できないので注意。
* VPNコネクションは一切変更不可のリソースのため、Update処理は実装されていない。しかし、完全な ForceNew にはできていないため、 tunnel や ipsec に変更があると、変更処理をしようとしてしまうので、 ignore_changes が必須です。
	* Readした情報をtfstateに戻している関係で、作成した瞬間から差分ありの状態になります。


[1]:https://github.com/nifcloud/nifcloud-sdk-python
[2]:https://pfs.nifcloud.com/api/rest/AuthorizeSecurityGroupIngress.htm
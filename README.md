# hakaru [![Build Status][travis-img]][travis-url]

[travis-img]: https://travis-ci.com/voyagegroup/hakaru.svg?token=iBCGFnZyWWvHWvMJnnx3&branch=master
[travis-url]: https://travis-ci.com/voyagegroup/hakaru

hakaru: 素朴な計測サーバ

## 1st step

- デプロイを実施する
- AMIをビルドする

## deployment

1. ビルドを実施し、成果物をアップロードする

```bash
$ make upload
```

1. blue/green or in-place のどちらかを実施する

### build AMI

```bash
$ cd provisioning/ami
$ make
```

### launch EC2 instance

- インスタンスタイプ: c5.large
- サブネット: プライベートサブネット
- iam: hakaru
- セキュリティグループ: hakaru
- ユーザデータに ./user_data.sh の内容を記述する

### blue/green deployment

1. AMI をビルドする
1. AMIからEC2インスタンスを起動する
1. 起動するEC2インスタンスの User data に ./user_data.sh の内容をコピペする
1. EC2インスタンスをロードバランサーに紐付る
1. 古いEC2インスタンスを終了する

### in-place deployment

1. 既にEC2インスタンスを起動していること
1. インスタンス上でユーザデータ ./user_data.sh の内容を実行する

## よくやるオペレーション

#### デプロイ！

ローカルのhakaruで

```
make upload
```

AWS システムマネージャーのセッションマネージャーでhakaruに入る。

```
sudo su
cd /root/hakaru
make deploy ARTIFACTS_COMMIT=latest
```

#### 負荷書ける

AWS システムマネージャーのセッションマネージャーでkakeruに入る。
シナリオを変えるのを忘れずに

```
sudo su
cd /opt/hakeru
make -C app kakeru upload HOST=hakaru-1043865531.ap-northeast-1.elb.amazonaws.com SCENARIO=2
```

## マルチノードやりかた

ssm で kakeru-worker にログインする(複数のworkerがあるがどれでも良い)
```
sudo su ubuntu
cd /opt/kakeru
make -C app kakeru/multinode HOST=hakaru-1043865531.ap-northeast-1.elb.amazonaws.com SCENARIO=2
```


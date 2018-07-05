ここでは`kubectx` と `kubens` を提供します。


**`kubectx`** を使うことでクラスターの行き来がしやすくなります。
![kubectx demo GIF](img/kubectx-demo.gif)

**`kubens`** を使うことでネームスペース間の移動がスムーズになります。
![kubens demo GIF](img/kubens-demo.gif)

# kubectx(1)

kubectxはcontextsの管理や切り替えのためのツールです。

```
USAGE:
  kubectx                   : contextsリストを表示
  kubectx <NAME>            : <NAME>というcontextにスイッチ
  kubectx -                 : 直近のcontextにスイッチ
  kubectx <NEW_NAME>=<NAME> : contextの<NAME>を<NEW_NAME>に変更
  kubectx <NEW_NAME>=.      : 現在のcontextを<NEW_NAME>に変更
  kubectx -d <NAME>         : <NAME>というcontextを削除 ('.' は現在のcontext)
                              (このコマンドはconrtextで使用されているユーザーやクラスタのエントリは削除しません)
```

### 使い方

```sh
$ kubectx minikube
context "minikube"にスイッチします。

$ kubectx -
context "oregon"にスイッチします。

$ kubectx -
context "minikube"にスイッチします。

$ kubectx dublin=gke_ahmetb_europe-west1-b_dublin
context "dublin"を設定します。
"gke_ahmetb_europe-west1-b_dublin"を"dublin"としてエイリアスを作成します。
```

`kubectx`は長いcontext名の管理のためにbash/zsh/fishにて<kbd>Tab</kbd>をサポートしています。 
context名を丸暗記する必要はないのです。

-----

# kubens(1)

kubensはnamespace間の移動のためのツールです。

```
USAGE:
  kubens                    : namespaceの一覧を表示します。
  kubens <NAME>             : アクティブなnamespaceを<NAME>に変更します。
  kubens -                  : 直近のnamespaceにスイッチします。
```


### 使い方

```sh
$ kubens kube-system
Context "test" をセットします。
アクティブなnamespaceは "kube-system"です。

$ kubens -
Context "test"をセットします。
アクティブなnamespaceは"default"です。
```

`kubens`も同様にbash/zsh/fishでの<kbd>Tab</kbd>をサポートしています。

-----

## インストール

### macOS

:confetti_ball: [Homebrew](https://brew.sh/)を使う。

    brew install kubectx

このコマンドは自動的にbash/zsh/fishをセットアップします。

- `brew install`を`--with-short-names`オプションを付けて実行することで`kctx`と`kns`というコマンドでインストールすることができます。`kubectl`との名前の衝突を避けるのが目的です。

- もしあなたがシェルプロンプト(`$PS1`)にcontextやnamespaceの情報を追加したい場合、[kube-ps1](https://github.com/jonmosco/kube-ps1)を試してみることをオススメします。

### Linux

`kubectx`と`kubens`はBashで書かれているため、Bashが使えるPOSIX環境でインストールが可能です。

- `kubectx`と`kubens`スクリプトをダウンロードします。
- どちらかを実行:
  - スクリプトを`PATH`の通っている場所に移動させます。
  - もしくはスクリプトをディレクトリに移動させます。そして、`/usr/local/bin`など`PATH`の通っている場所から`kubectx`や`kubens`へシンボリックリンクを作成します。
- `kubectx`と `kubens`に実行権限を付与します。 (`chmod +x ...`)
- bash/zsh/fishでのインストール方法を見つけよう[completion scripts](completion/).

インストール例

``` bash
sudo git clone https://github.com/ahmetb/kubectx /opt/kubectx
sudo ln -s /opt/kubectx/kubectx /usr/local/bin/kubectx
sudo ln -s /opt/kubectx/kubens /usr/local/bin/kubens
```
#### Arch Linux

非公式 [AURパッケージ](https://aur.archlinux.org/packages/kubectx) `kubectx`を利用できます。
インストール手順はこちらから [Arch 
wiki](https://wiki.archlinux.org/index.php/Arch_User_Repository#Installing_packages).

-----

### アクティブなcontextの色を変更する。

アクティブなnamespaceやcontextの色を変更したい場合は`KUBECTX_CURRENT_FGCOLOR`と`KUBECTX_CURRENT_BGCOLOR`という変数に値を設定してください。

```
export KUBECTX_CURRENT_FGCOLOR=$(tput setaf 6) # 青字
export KUBECTX_CURRENT_BGCOLOR=$(tput setaf 7) # 白背景
```

カラーコードは[こちら](https://linux.101hacks.com/ps1-examples/prompt-color-using-tput/)

-----

####  ユーザー

| kubectxについてどう思いますか? |
| ---- |
| _“Thank you for kubectx & kubens - I use them all the time & have them in my k8s toolset to maintain happiness :) ”_ – [@pbouwer](https://twitter.com/pbouwer/status/925896377929949184) |
| _“I can't imagine working without kubectx and especially kubens anymore. It's pure gold.”_ – [@timoreimann](https://twitter.com/timoreimann/status/925801946757419008) |
| _“I'm liking kubectx from @ahmetb, makes it super-easy to switch #Kubernetes contexts [...]”_ &mdash; [@lizrice](https://twitter.com/lizrice/status/928556415517589505) |
| _“Also using it on a daily basis. This and my zsh config that shows me the current k8s context 😉”_ – [@puja108](https://twitter.com/puja108/status/928742521139810305) |
| _“Lately I've found myself using the kubens command more than kubectx. Both very useful though :-)”_ – [@stuartleeks](https://twitter.com/stuartleeks/status/928562850464907264) |
| _“yeah kubens rocks!”_ – [@embano1](https://twitter.com/embano1/status/928698440732815360) |
| _“Special thanks to Ahmet Alp Balkan for creating kubectx, kubens, and kubectl aliases, as these tools made my life better.”_ – [@strebeld](https://medium.com/@strebeld/5-ways-to-enhance-kubectl-ux-97c8893227a)

> もし`kubectx`を気に入ってくれたなら、もう一つの僕のプロジェクト[`kubectl-aliases`](https://github.com/ahmetb/kubectl-aliases)も見てね。

-----

免責事項: Googleの公式プロダクトではありません。


#### Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/ahmetb/kubectx.svg)](https://starcharts.herokuapp.com/ahmetb/kubectx)


---
title: コマンドライン
description: tuna CLI で、シェルやスクリプトから tuna のセッションとペインを操作します。
---

# コマンドライン

`tuna` コマンドは、あなたのシェル、スクリプト、エージェントなど、どのシェルからでもセッションを操作できます。**tuna › コマンドラインツールをインストール…** でインストールします。

## 対象

ほとんどのコマンドは**対象**を取ります: セッションかペインです。対象には次のものを指定できます:

- ID、またはその一意な前方部分
- セッションの名前
- セッション（フォーカスのあるペインを表す）

tuna の中では、`$TUNA_SESSION` と `$TUNA_PANE` がコマンドを実行するセッションとペインを示し、それがデフォルトの対象になります。

## セッションとペイン

```sh
tuna ls                          # すべてのセッションとそのペイン
tuna new -n api -d ~/code/api    # セッションを開始し、その ID を出力
tuna new -- npm run dev          # コマンドを実行するセッションを開始
tuna split --down -- npm test    # このペインの横にペイン
tuna split --url localhost:3000  # ページを表示するブラウザペイン
tuna rename api backend          # セッション名を変更
tuna focus backend               # アプリで表示
tuna close backend               # セッションを終える（ペインを指定するとそのペイン）
tuna close fix-login --remove-worktree   # 専用のワークツリーも削除する
tuna restart <pane>              # 終了したペインのシェルを再び開始
```

`tuna new` のオプション:

| フラグ                      | 意味                                            |
| --------------------------- | ----------------------------------------------- |
| `-p`, `--project <name\|id>` | 入れるプロジェクト（デフォルト: ディレクトリから決まる） |
| `-n`, `--name <name>`       | セッションに名前を付ける                        |
| `-d`, `--dir <dir>`         | そこで開始する（デフォルト: 現在のディレクトリ） |
| `-f`, `--focus`             | アプリで表示する                                |
| `-w`, `--wait`              | コマンドの終了を待ち、その終了ステータスで終了する |
| `-u`, `--url <url>`         | シェルの代わりに、ページを表示するブラウザを開く |
| `--worktree <branch>`       | 専用の git ワークツリーで、新しいブランチで開始する（[詳しく](./changes.md#sessions-in-their-own-worktree)） |
| `--base <ref>`              | ワークツリーのブランチの起点（既定: デフォルトブランチ） |

## 入力と読み取り

```sh
tuna send api "npm run migrate" --enter   # テキストを入力してから Enter を押す
tuna send api -e 'y\n'                    # \n \r \t \e \xHH のエスケープを解釈
tuna key api ctrl+c                       # キーを押す
tuna read api -n 40                       # 画面の最後の 40 行
tuna read api -a                          # スクロールバックを含む画面
```

`tuna key` は `enter`、`tab`、`escape`、`up`、`down`、`left`、`right`、`backspace` と、`ctrl+c` や `ctrl+d` のような修飾キーを理解します。

## 待機

```sh
tuna wait api                             # 作業が止まるまで
tuna wait api --for waiting,error --timeout 10m
tuna prompt review "summarize the diff"   # 送信 + Enter、待機、画面を出力
```

`tuna wait` は戻るときにペインの状態を出力します。`tuna prompt` はテキストと Enter を送り、ペインの作業が止まるまで待ってから画面を出力します。ほかのエージェントに何かを尋ねる最も手早い方法です。

### コマンドの終了を待つ

`make test` のようなコマンドは*作業中*と表示されないので、終了を待つには `--wait` を付けて開始します。すると `tuna` はそのコマンド自身の終了ステータスで終了します:

```sh
if tuna split --down --wait -- make test; then
  echo "tests passed"
fi
```

`tuna wait --exit` は、開始時のコマンドをすでに実行しているペインに対して同じことをします:

```sh
pane=$(tuna split --down -- go test ./...)
# … ほかのことをする …
tuna wait "$pane" --exit --timeout 10m
echo "go test exited with $?"
```

コマンドの終了後も出力とともにペインは開いたままなので、`tuna read` で読めます。コマンドの終了を判断できない `cmd` をシェルにしている場合、コマンドの終了待ちは使えません。

## イベントの監視

`tuna events` は、セッションとペインに起きたことを起きたときに 1 イベント 1 行で出力し、止めるまで続けます。スクリプトは `tuna status` をポーリングせずに、すべてのエージェントを一度に追えます。

```sh
tuna events                               # すべてのセッション
tuna events api                           # 1 つのセッション、または 1 つのペイン
tuna events -p tuna                       # 1 つのプロジェクトのセッション
tuna events --for waiting,error           # これらの状態になったペインだけ
tuna events --json | jq -r .name          # 1 行に 1 つの JSON オブジェクト
tuna events --for waiting --once          # 最初の 1 件で終了
```

```
12:04:31  api x78n9ss  working → waiting  Allow Bash: rm -rf dist?
12:04:40  api  renamed api → backend
```

| イベント | いつ | フィールド |
| --- | --- | --- |
| `pane.status` | ペインの状態が変わったとき | `pane`、`from`、`to`、`message`、`kind`、`agent`。exited になったときは `exitCode` も |
| `command.exited` | ペインの開始時のコマンドが終わり、シェルが続くとき | `pane`、`exitCode` |
| `pane.created`、`pane.closed` | ペインが開始または終了したとき | `pane` |
| `session.created`、`session.closed` | セッションが開始または終了したとき | |
| `session.renamed` | セッションの名前が変わったとき | `from`、`to` |
| `session.moved` | セッションが別のプロジェクトに移ったとき | `from`、`to`（プロジェクトの ID） |
| `notice` | セッション内のプログラムがデスクトップ通知を送ったとき | `pane`、`title`、`body` |

`--json` では、どのイベントにも `event`、`time`（unix ミリ秒）、`session`、`project`、`name`（サイドバーに表示されるセッション名）が付きます。

`tuna events` が報告するのは今からの変化です。すでにあるものは `tuna ls --json` で確認します。対象を指定すると、対象が閉じたときに止まります。`--once` は最初に一致したイベントで終了し、対象が先に閉じるか `--timeout` が過ぎると `1` で終了します。tunad は状態が落ち着いてから送るので、一瞬だけの状態（短い作業中など）は出ないことがあります。

<a id="browsing"></a>

## ブラウザの操作

`tuna browser` は、`tuna send` と `tuna read` がターミナルを操作するのと同じように、ブラウザペインのページを操作します: クリック、入力、ページの読み取り、スクリーンショット。対象はブラウザペイン、またはセッション（そのブラウザ）です。対象を指定しなければ、実行したセッションのブラウザを操作します。

```sh
tuna split --url localhost:3000            # ページを開き、ペインの ID を出力
tuna browser open localhost:3000/todos     # ページへ移動し、読み込みを待つ
tuna browser type '#new-todo' "Buy milk" --enter
tuna browser wait '.todo-list li'          # 要素が表示されるまで
tuna browser wait --text "1 item left"     # テキストが表示されるまで
tuna browser text                          # ページのテキスト
tuna browser text h1                       # 要素のテキスト
tuna browser click 'button[type=submit]'
tuna browser screenshot -o page.png
tuna browser eval 'document.title'
```

| コマンド                        | 動作 |
| ------------------------------- | ---- |
| `open <url>`                    | ページへ移動し、読み込みを待つ |
| `back`、`forward`、`reload`     | 戻る、進む、または再読み込みして、それを待つ |
| `text [selector]`               | ページの表示テキスト、または最初に一致した要素のテキストを出力 |
| `html [selector]`               | ページ、または要素の HTML を出力 |
| `click <selector>`              | マウスと同じように要素をクリック |
| `type <selector> <text…>`       | フィールドの内容を置き換えてテキストを入力。`--enter` を付けるとその後 Enter を押し、フォームを送信する |
| `wait [selector] [--text text]` | 要素かテキストが表示されるまで待つ。どちらもなければページの読み込みまで |
| `eval <js…>`                    | ページで JavaScript を実行し、その値を出力: 文字列はそのまま、それ以外は JSON |
| `screenshot [-o file.png]`      | ページの画像を保存し、その場所を出力。`-o -` で標準出力に書き出す |

セレクタは CSS セレクタです。シェルでは引用符で囲んでください。要素が必要なコマンドは要素が現れるのを待ち、`--timeout`（デフォルト `10s`）でそれらとページの読み込みにかけてよい時間を設定します。`type` と `eval` は、`-` で標準入力からテキストを受け取ります。

クリックと入力は本物のマウスとキーボードと同じイベントを送るので、React などのフレームワークで作られたページもそれを受け取ります。スクリーンショットにはペインが画面に表示されている必要があります: 先に `tuna focus` で表示してください。それ以外は、表示されているかどうかにかかわらず、どのブラウザペインでも動きます。

`tuna do` と同じく、`tuna browser` には **設定 › エージェント** のリモート操作と、アプリが開いていることが必要です。

## 過去の会話

```sh
tuna history                     # ペインで動いたエージェントの会話を、最後に動いたものから順に
tuna history -p tuna -n 50       # ひとつのプロジェクトの会話を、最大 50 件
tuna history login --json        # 「login」に一致するものを JSON で
tuna resume 13b9                 # 再開する: ID、その先頭部分、または名前で
tuna resume fix-login -f         # 再開して、アプリに表示する
```

`tuna resume` は、その会話を再開するセッションを会話のディレクトリとプロジェクトで開始し、セッションの ID を表示します。まだペインで動いている会話なら、そのセッションを返します。[過去の会話](./sessions.md#past-conversations)を参照してください。

## 状態

```sh
tuna status api --json                    # 状態、プログラム、ディレクトリ
tuna report working "Running migrations"  # 自分の状態を報告
tuna notify "Need your input"             # このペインを対応が必要なものとして示す
```

状態については[状態と通知](./status.md)を参照してください。

## 変更

```sh
tuna diff                                 # ペインの未コミットの変更を git diff の形式で
tuna diff --turn                          # エージェントが直前のターンで変更したもの
tuna diff review --branch                 # デフォルトブランチ以降
```

[変更](./changes.md)を参照してください。

## 構成とセットアップ

```sh
tuna config path                 # 設定ファイルの場所
tuna config show                 # すべてのキーとその値
tuna config show --defaults      # すべてのキーの初期値
tuna config edit                 # $EDITOR で開く
tuna config check                # 不明なキーと不正な値を報告
tuna projects                    # プロジェクトを一覧
tuna skill setup                 # コーディングエージェントに tuna の使い方を教える
tuna daemon status|stop|restart  # tunad を管理
tuna daemon install-service      # マシンの起動とともに tunad を実行（systemd、launchd）
tuna daemon serve                # サービスマネージャ向けに tunad をフォアグラウンドで実行
tuna clients                     # tunad に接続しているもの
tuna pets                        # ウィンドウ上端のペット、表示中のものが先頭
tuna pets use <名前|tuna|codex>  # ペットを表示する：~/.codex/pets のもの、または tuna のもの
tuna pets install <pack.zip>     # Codex ペットパック（.codex-pet.zip）をインストールする
tuna pets rm <名前>              # tuna がインストールしたペットを削除する
```

## Feishu

```sh
tuna feishu setup                # アプリの ID と secret を尋ね、保存して接続し、次にすることを示す
tuna feishu status               # ボットが接続しているか
tuna feishu test                 # Feishu がアプリを認識するか確認する
tuna feishu send <チャット> <テキスト…>  # チャット ID でメッセージを送る
```

Feishu のチャットは、結び付いたセッションのコーディングエージェントを動かします: `/new` で始め、`/use` で開いているセッションに結び付け、そのほかのテキストはプロンプトです。[Feishu](./feishu.md)を参照してください。

## ほかのマシン

コマンドの前に `-H <host>` を付けるか `$TUNA_HOST` を設定すると、セッションに対するコマンドを別のマシンの tunad に向けられます: 追加したホスト、または `me@devbox` や `wsl://Ubuntu` のようなアドレスです。[リモートホスト](./remote.md)を参照してください。

```sh
tuna hosts                       # 追加したホスト
tuna hosts add devbox me@devbox  # 追加する（アプリにも表示される）
tuna hosts check devbox          # 接続して、結果を表示
tuna hosts install devbox        # SSH でそこに tuna をインストール
tuna ls --all                    # このマシンとすべてのホストのセッション
tuna -H devbox ls
tuna -H devbox new -n review -- claude
tuna -H devbox daemon update     # そこの tuna をこちらのバージョンにアップデート
```

ペインの中では、`tuna` はアプリがどこにあっても、ペインが動いているマシンの tunad と通信します。

`tuna ls --all` はすべてのホストに同時に問い合わせ、それぞれのセッションを `-H` に渡す名前（このマシンは `local`）の下に一覧表示します。届かないホストは、セッションの代わりにその理由を表示します。`--json` を付けると `{host, addr, error, sessions}` のリストを出力します。

## アプリのアクションを実行する

`tuna do` は、メニューと同じアプリのアクションをどれでも実行します。**設定 › エージェント** のリモート操作（デフォルトでオン）が必要です。

```sh
tuna actions                     # 一覧を表示
tuna do pane.split.right
tuna do pane.browser url=localhost:3000
tuna do sidebar.toggle
```

## 終了ステータス

| コード | 意味                                 |
| ------ | ------------------------------------ |
| `0`    | 成功                                 |
| `1`    | 失敗                                 |
| `2`    | 使い方の誤り                         |
| `3`    | 対象が存在しない、またはあいまい     |
| `4`    | tunad に接続できない                 |

プログラムを待つコマンド（`tuna wait --exit` と `--wait`）は、代わりにそのプログラムの終了ステータスで終了します。ペインのシェルがエラーで終了したときの `tuna wait` も同じです。タイムアウトした待機は `1` で終了します。

`ls`、`status`、`projects` に `--json` を付けると、機械で読める形式で出力します。

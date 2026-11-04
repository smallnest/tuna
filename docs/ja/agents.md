---
title: コーディングエージェント
description: Claude Code、Codex などのエージェントに、自分の横で作業を動かし、tuna に進捗を報告することを教えます。
---

# コーディングエージェント

tuna で動くエージェントは、あなたと同じように `tuna` コマンドを使えます: 下のペインで開発サーバーを起動し、2 つ目のエージェントにレビューを任せ、テストを待ち、行き詰まったらあなたに知らせます。

<a id="install-the-skill"></a>

## スキルをインストールする

```sh
tuna skill setup
```

これで 2 つのものがインストールされます:

- エージェントに CLI を教える **スキル**。Claude Code は `~/.claude/skills/tuna` から、Codex、OpenCode などのエージェントは `~/.agents/skills/tuna` から読み込みます。
- エージェントがセッションを始めるときに実行される **フック**。エージェントに、どの tuna セッションとペインにいるかを伝えるので、それを調べる手間が省けます。また、名前のないセッションに `fix-login-redirect` のようなタスクにちなんだ名前を付けさせるので、あなたが付ける必要はありません。さらに、エージェント自身のセッションを tuna に伝えるので、再起動のあとなど tunad が再起動したとき、エージェントは[会話を再開](./sessions.md#tunad)します。

フックは、インストールされているエージェントごとに、それぞれが探す場所に置かれます:

| エージェント | フックの場所                                          |
| ------------ | ----------------------------------------------------- |
| Claude Code  | `~/.claude/settings.json` の `SessionStart` フック     |
| Codex        | `~/.codex/hooks.json` の `SessionStart` フック         |
| Pi           | 拡張機能 `~/.pi/agent/extensions/tuna.ts`              |
| OpenCode     | プラグイン `~/.config/opencode/plugins/tuna`           |

`CLAUDE_CONFIG_DIR`、`CODEX_HOME`、`PI_CODING_AGENT_DIR`、`XDG_CONFIG_HOME` でエージェントのファイルの場所が変わっている場合、tuna はそれに従い、各設定ファイルのほかの部分はそのままにします。

[Feishu](./feishu.md) のチャットでは、Claude Code に `Stop` フックも入ります。各ターンが終わると、フックはエージェントの最後のメッセージを、ボットがそのペインに渡したファイルに書き込みます。ボットはターミナルの画面ではなく、そのメッセージで返信します。有効にするには `tuna skill setup` をもう一度実行します。ほかのエージェントはターミナルのテキストを使います。

1 つのエージェントだけをセットアップするには、`--agent` に `claude`、`codex`、`pi`、`opencode` のいずれかを渡します。セットアップは一度だけです: スキルは tuna 自身のコピーへのリンクで、tuna のアップデートに合わせて更新されます。Pi と OpenCode のプラグインも tuna が最新に保ちます。`tuna skill status` で、各部分が最新かどうかがわかります。

```sh
tuna skill status    # インストール済みのものと、それが最新かどうか
tuna skill show      # スキルを出力
tuna skill remove    # スキルとフックをアンインストール
```

フックが働くのはセットアップ後に始めたエージェントのセッションだけで、エージェントが先に確認を求めることもあります:

- **Claude Code** は、信頼したフォルダでのみフックを実行します。
- **Codex** は新しいフックの承認を求めます: `/hooks` で tuna のフックを承認してください。
- **OpenCode** は、バックグラウンドサービスの起動時にプラグインを読み込みます: 作業中のセッションがないときに、一度 `opencode service restart` を実行してください。

OpenCode はエージェントのコマンドをペインではなくバックグラウンドサービスで実行するので、`TUNA_SESSION` と `TUNA_PANE` がありません。代わりに tuna がエージェントにペインを伝え、エージェントはそれを対象として指定します。同じ理由で、OpenCode のエージェントからの `tuna report` はそのペインに届きません。

tuna の外では、フックは何もしません。

## エージェントにできること

自分の横で何かを動かす:

```sh
tuna split --down -- npm run dev
tuna new -n review -- claude
```

それを操作して結果を読む:

```sh
tuna prompt review "review the changes in src/auth"
tuna wait tests --for idle,error --timeout 10m
tuna read tests -n 60
```

自分の横でテストを実行して、失敗した内容を読む:

```sh
pane=$(tuna split --down -- go test ./...)
tuna wait "$pane" --exit || tuna read "$pane" -n 80
```

状況をあなたに伝える:

```sh
tuna report working "Migrating the database" --progress 40
tuna report blocked "Approve the production deploy?"
tuna report done "Deployed to staging"
```

`blocked` の報告でセッションは*待機中*になり、あなたに通知されます。エージェントがこれを使うのは、あなたなしでは続けられないときだけにしてください。

## リモート操作

**リモート操作** がオン（デフォルト、**設定 › エージェント**）なら、エージェントはアプリ自体も操作できます: `tuna focus` でセッションを表示したり、`tuna do` でメニューのアクションを実行したりします。エージェントを自分のセッションだけにとどめるには、`remote-control = false` でオフにします。

## 環境変数

すべてのペインには次の変数があります:

| 変数           | 値                             |
| -------------- | ------------------------------ |
| `TUNA_SESSION` | ペインのセッションの ID        |
| `TUNA_PANE`    | ペインの ID                    |

スクリプトが tuna の中で動いているかどうかは、`$TUNA_PANE` で確認できます。

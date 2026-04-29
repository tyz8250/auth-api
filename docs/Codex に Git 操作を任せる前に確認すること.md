## Codex に Git 操作を任せる前に確認すること

横着して、ちゃんと確認せずに承認してしまった。

今回は、Codex に渡せる作業フローを作ろうと思い、新しい Go プロジェクト `auth-api` を作成した。
ローカルにプロジェクト用のフォルダを作り、必要なファイルを作成し、それを GitHub に `push` するところまで Codex に任せようとした。

正直、そこまで難しい作業ではないと思っていた。

新しいフォルダで `git init` して、初回コミットを作って、GitHub に新規リポジトリを作って `push` するだけ。
そう考えていた。

実際、`push` 自体はできた。

しかし、その後 VSCode の Git 表示を見ると、変更数が +100 を超えているように見えた。
明らかにおかしい。

本来であれば、`auth-api` 配下のファイルだけが Git の管理対象になるはずだった。
それなのに、`.DS_Store` や Anaconda 関連のファイル、キャッシュファイルのような、プロジェクトとは関係ないファイルまで Git の変更として表示されていた。

そこで履歴を追うと、どうやら `auth-api` 自体が独立した Git リポジトリとして扱われておらず、親ディレクトリである `/Users/yourname` 側の `.git` を参照していた。

かなりまずい状態だったので、急いで GitHub のリポジトリを削除し、作り直した。

## 何が起こっていたのか

原因は、`auth-api` の中に `.git` がない状態で Git 操作を始めたことだった。

Git は、現在のディレクトリに `.git` がない場合、親ディレクトリへさかのぼって `.git` を探す。

今回の場合、`/Users/yourname/auth-api` で作業しているつもりだったが、Git は親ディレクトリの `/Users/yourname` にある `.git` を見ていた。

つまり、Git から見ると作業対象はこうなっていた。

自分の認識:

```text
/Users/yourname/auth-api がリポジトリ
```

Git の認識:

```text
/Users/yourname がリポジトリ
auth-api はその中の未追跡ディレクトリ
```

この状態で `git status` や `git remote -v` を実行すると、`auth-api` ではなく、親ディレクトリ側の Git 情報を見てしまう。

さらに危険なのは、この状態で `git add .` や `git commit` を実行すると、プロジェクト外のファイルまで Git の管理対象に入れてしまう可能性があること。

今回は VSCode 上で大量の変更が見えたため、異常に気づくことができた。

## Codex は最初に警告していた

今回一番反省すべき点は、Codex の最初の確認結果をちゃんと読んでいなかったこと。

Codex は作業前に、次のような確認をしていた。

```bash
git branch --show-current # 現在のブランチ名を表示
git remote -v # リモートリポジトリのURLを表示
git status --short --branch # ステータスを短縮形式で表示
gh auth status # GitHubの認証状態を表示
```

そして、こう指摘していた。

```text
今の auth-api は単独リポジトリではなく、
親ディレクトリ /Users/yourname 側の Git 管理下に
未追跡ディレクトリとして見えています。
```

これはかなり重要な警告だった。

ここで一度止まって、`git rev-parse --show-toplevel` を確認し、
`auth-api` がリポジトリのルートディレクトリであることを確認すべきだった。
それをせずに「たぶん大丈夫だろう」と進めてしまった。

## sandbox のエラーについて

途中で、Codex は Git 操作中に次のようなエラーを出していた。

```text
Unable to create '.git/HEAD.lock': Operation not permitted
Unable to create '.git/index.lock': Operation not permitted
```

これは Git が `.git` 内部の lock ファイルを作ろうとしたが、Codex の sandbox 制限により書き込みが許可されなかった、という意味だと思われる。

lock ファイルは、Git の内部ファイルを安全に更新するための一時的なロック札のようなもの。
Git は `HEAD` や `index` のような大事な管理情報を書き換えるとき、同時に別の Git 操作が走って壊れないように、まず `.lock` ファイルを作る。

最初は単なる権限エラーに見えた。

しかし、今振り返ると、Codex が書き換えようとしていた `.git` が `auth-api` のものではなく、親ディレクトリ `/Users/yourname` 側の `.git` だった可能性がある。

そう考えると、sandbox の制限で止まったことは、むしろ事故を防いでくれた可能性がある。

## 今後の対策

Codex に Git 操作を任せる前に、必ず次のコマンドを確認する。

```bash
pwd
git rev-parse --show-toplevel
git status --short --branch
git remote -v
```

特に大事なのはこれ。

```bash
git rev-parse --show-toplevel
```

この結果が、今 `push` したいプロジェクトのディレクトリと一致しているか確認する。

今回であれば、期待する結果はこれ。

```text
/Users/yourname/auth-api
```

もしこれが次のようになっていたら、その時点で Git 操作を止める。

```text
/Users/yourname
```

また、Codex に依頼するときは、最初に次の条件を入れる。

```text
作業前に必ず pwd、git rev-parse --show-toplevel、git status --short --branch、git remote -v を確認してください。

git rev-parse --show-toplevel が今回のプロジェクトディレクトリと一致しない場合は、git add、commit、push は実行せず、そこで止まってください。

また、git add . は使わず、必要なファイルだけを明示して add してください。
```

新規プロジェクトの初回 `git init` は、自分で実行してもよいかもしれない。
少なくとも、`git init` 後に `git rev-parse --show-toplevel` を確認してから Codex に任せる。

## 学んだこと

Git は、今 VSCode で開いているフォルダを自動的にリポジトリとして扱うわけではない。

Git は、現在地から親ディレクトリへさかのぼり、最初に見つけた `.git` を基準に動く。

だから、作業前に Git がどこをリポジトリのルートとして認識しているかを確認する必要がある。

今回の教訓はこれ。

```text
Codex に Git 操作を任せる前に、
Git が見ている作業範囲を確認する。
```

便利だからこそ、承認前に見るべきものは見る。
特に Git 操作は、コード生成よりも事故の影響範囲が広い。

だから Codex に任せるとしても、
「どのリポジトリを触っているか」だけは人間が確認する。

## まだ残っている謎

というか、そもそも `/Users/yourname/` でいつ `git init` したんだ俺は……。

ただ、これも今回の学びだと思う。

Git は、一度どこかに `.git` ができると、その配下のディレクトリにも影響する。
だから「いつ作ったか覚えていない `.git`」が親ディレクトリにあるだけで、あとから作った別プロジェクトにも影響してしまう。

今後は、新しいプロジェクトを作ったら、最初に次を確認する。

```bash
pwd
git rev-parse --show-toplevel
```

この2つを見て、自分が作業している場所と、Git がリポジトリだと思っている場所が一致しているか確認する。

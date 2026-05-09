# VRCLaunch

VRChat 用の軽量ランチャー。Go + [Gio](https://gioui.org) で実装、単一バイナリ約 10MB（Electron 製の VRCQL 比で大幅に軽量）。

複数プロファイルの管理、VR / Desktop モード切替起動、起動オプション設定をシンプルな UI で提供します。

## 特徴

- **単一バイナリ・CGO 不要**: Windows 用にクロスコンパイルして配布できる
- **VR モード時の自動軽量化**: VR で遊んでいる間はデスクトップ側ウィンドウを最小サイズ (320x240) で起動して GPU リソースを節約
- **Desktop モードはプロファイル別解像度**: 保存された解像度・全画面設定で起動
- **プロファイル別起動引数**: `--profile=N`、`--fps=N`、Custom Args
- **永続化**: `%APPDATA%\VRCLaunch\config.json`（Linux/macOS は `os.UserConfigDir()` 準拠）

## 要件

- Steam 経由でインストールされた VRChat
- Steam の `launch.exe` のフルパス（例: `F:\SteamLibrary\steamapps\common\VRChat\launch.exe`）

## ビルド

```bash
# Windows 用クロスビルド (Linux/WSL から)
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -s -w" -o VRCLaunch.exe .

# ネイティブビルド (Windows 上)
go build -ldflags="-H windowsgui -s -w" -o VRCLaunch.exe .
```

## 使い方

1. 起動して右上の **Settings** から `launch.exe` のパスを設定
2. **+ Profile** でプロファイルを追加（Name と Profile Index は必須、index は VRChat の `--profile=N` に渡される）
3. プロファイルを選択し **Launch VR** または **Launch Desktop** で起動

## 設定ファイル

`%APPDATA%\VRCLaunch\config.json`:

```json
{
  "version": 1,
  "launch_path": "F:\\SteamLibrary\\steamapps\\common\\VRChat\\launch.exe",
  "last_selected": "abc123...",
  "profiles": [
    {
      "id": "abc123...",
      "name": "Main",
      "index": 1,
      "options": {
        "fps": 90,
        "screen_width": 1920,
        "screen_height": 1080,
        "screen_fullscreen": true,
        "custom_args": ""
      }
    }
  ]
}
```

## 開発

```bash
# テスト + カバレッジ
go test -race -cover ./internal/...

# 静的解析
go vet ./...
```

### パッケージ構成

| パッケージ | 役割 | カバレッジ |
|-----------|------|------------|
| `internal/config` | プロファイル/設定スキーマと JSON 永続化 | 83.9% |
| `internal/launcher` | 起動引数組み立て + プロセス spawn | 96.2% |
| `internal/uistate` | UI 状態機械（純粋ロジック）+ フォームバリデーション | 100% |
| `internal/ui` | Gio レンダリング（main / profile editor / settings 各 view） | — |

UI レンダリング層は Gio ウィジェットへのレイアウト委譲が中心のため、状態機械を `internal/uistate` に切り出して単体テスト可能にしています。

## 既存ランチャーとの関係

このリポジトリには参考実装として `VRCQL/`（Electron 製、フル機能、稼働中）と `VRCL/`（Tauri 製、軽量、不安定）があります（gitignore 対象）。VRCLaunch はそれらのコア機能のみを Go + Gio で再構築した第三の実装です。

## ライセンス

未指定

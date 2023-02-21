# dropbox-uploader
ローカルのディレクトリを指定すると、そのディレクトリをまるごと dropbox にアップロードする。
その際に、`refresh_token` を環境変数で指定したものを利用し、都度 `access_token` を更新する。

```bash
APP_KEY=*** \ 
APP_SECRET=*** \
REFRESH_TOKEN=*** \
./dropbox-uploader start
```


## 制約
- 1ファイルサイズ150MBまで
- 現状は、既存ファイルがある場合は上書きしない仕様

## Dropbox App登録（参考）
1. Regist Your app
```
https://www.dropbox.com/oauth2/authorize?client_id=<your App key>&response_type=code&token_access_type=offline
```

2. Get Refresh Token
```
curl https://api.dropbox.com/oauth2/token \
    -d code=<your got 1. code> \
    -d grant_type=authorization_code \
    -u <your app key>:<your secret>
```

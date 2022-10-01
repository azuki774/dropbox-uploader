# dropbox-uploader
dropbox-uploader provides a server for uploading dropbox.
- Usage
  ```
  $ curl -XPOST -F "file=@<upload file>" http://localhost:8080/upload?path=<dropbox path>
  ```
  - Ex. `curl -XPOST -F "file=testimage.jpg" http://localhost:8080/upload?path=/testpath/testimage.jpg`

This program uses unonfficial SDK for Go : https://github.com/dropbox/dropbox-sdk-go-unofficial

## Docker Usage (Use docker compose)
```
$ make start
```

## Get Token and Setting Environment
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

3. Set Your information
- `/deployment/envfile`

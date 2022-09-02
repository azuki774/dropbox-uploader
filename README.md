# dropbox-uploader

## Usage
```
~/work/dropbox-uploader/build/bin$ ./dropbox-uploader upload -h
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  dropbox-uploader upload [flags]

Flags:
      --dry-run          dry run
  -d, --dst-dir string   Dropbox target directory
  -h, --help             help for upload
  -o, --overwrite        overwrite if exists same name file
  -s, --src-dir string   source file or directory
  -t, --token string     Dropbox access token
  -u, --update           update when exists same name file if it has an newer timestump.
```

## Get Token
1. Create you app
https://www.dropbox.com/developers/apps/create?_tk=pilot_lp&_ad=ctabtn1&_camp=create
    - Scoped access _> Full Dropbox -> Name your app

2. Add Permission
    - Permissions -> Check "files.content.write, files.content.read"

3. Get Token
    - Settings -> Generated access token

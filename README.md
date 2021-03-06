# g0chainwebserver
Go based 0chain webserver
Serves files stored on 0chain from specified authticket
 
### PRE-REQUSITES
#### - Existing $HOME/.zcn folder (the location for this executable)
#### - Containing zbox/zwallet cli tools plus existing wallet and allocation
#### - Build and Install this executable in above folder

## USAGE

Open dedicated console window in above folder

    ./g0chainwebserver --port <port.no> (default 6942)

### TYPICAL CONSOLE OUTPUT

    Listening on port: 6942
    SRV: /img.png  (localhost:6942)
    GET: /img.png  637.4 KB in 4.180 secs (152.5 KB/sec)
    SND: /img.png  637.4 KB in 0.007 secs (91.39 MB/sec)

### FILE SHARES

#### For File Authtickets (filename is decoded from filename included in authticket) :-

browser (or app) point to:-

    http://IPaddress:port/authticket/xxxxxxxxxxxxxx
    e.g. localhost:6942/authticket/xxxxxxxxxxxxxx

### FOLDER SHARES

#### For Folder Share Authtickets via full remote path (owner wallet only) :-

    http://IPaddress:port/authticket/xxxxxxxxxxxxxx/remote/path/file.ext
    e.g. http://192.168.1.50:6942/authticket/xxxxxxxxx/video/cat.m3u8 

#### For Folder Share Authtickets + authhash (file hash) (any wallet) :-

    http://IPaddress:port/authhash/xxxxxxxxxxxxxx/yyyyyyyyyyyyyy_file.ext
    e.g. http://192.168.1.50:6942/authhash/xxxxxxxxx/yyyyyyyyyyyyyy_cat.m3u8
    e.g. localhost:6942/authhash/xxxxxxxxx/yyyyyyyyyyyyyy_cat.m3u8

# STATUS
## Command Line parameters
#### v0.0.6
    --proto (default http) specify http or https
    --certpath specify path to ssl certificate if https proto selected
    --keypath specify path to ssl private key if https proto selected
#### v0.0.3
    --wallet (default wallet.json) specify different wallet file
    --allocation (default allocation.txt) specify different allocation file
#### v0.0.2
    --debug <1/0> (default 0) - 1 shows command line output to console
#### v0.0.1
    --config <configfile.yaml> (default config.yaml)
    --port <portno.> (default 6942)

## URL info paths
#### v0.0.5
    /enckeys (Show client_id and encryption_public_key for server wallet to enable encrypted file shares)
#### v0.0.1
    /status (Status OK)
    /config (Show currently selected config.yaml file contents)
    /version (show version number)

## UPDATES
### Fixes / Changes

#### v0.0.6
- Added --proto parameter for http/https plus --certpath, --keypath for https option
#### v0.0.5
- Added /enckeys url info path
#### v0.0.4
- Changed filehash separator to _ on folder shares to preserve path relative to authticket
- Fixed file size accuracy
#### v0.0.3
- Improved file size handling
#### v0.0.2
- Better error handling for incorrect paths
#### v0.0.1
- Fixed differences in file/folder handling (folder still requires allocation even though encoded in authticket)

### Additions
#### v0.0.3
- Added support for filehash for folder shares
#### v0.0.2
- Added file size / xfer speed to console output
#### v0.0.1
- Added support for default file for folder share. If no file specified, looks for .default file and then performs (301) redirect to file specified within

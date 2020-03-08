# g0chainwebserver
Go based 0chain webserver
Serves files stored on 0chain from specified authticket
 
PRE-REQUSITES
Existing $HOME/.zcn folder (the location for this executable)
Containing zbox/zwallet cli tools plus existing wallet and allocation
Build and Install this executable in above folder

## USAGE

Open dedicated command line window in above folder

./g0chainwebserver --port <port.no> (default 6942)

browser (or app) point to:-

http://IPaddress:port/authticket/xxxxxxxxxxxxxx

for files (filename and has included in authticket)


http://IPaddress:port/authticket/xxxxxxxxxxxxxx/remote/path/file.ext

for shared folders with specified full remote paths (owner wallet)

e.g. http://192.168.1.50:6942/authticket/xxxxxxxxx/video/cat.m3u8

OR

http://IPaddress:port/authhash/xxxxxxxxxxxxxx/yyyyyyyyyyyyyy/file.ext

for shared folders with file hash parameter (any wallet)

e.g. http://192.168.1.50:6942/authhash/xxxxxxxxx/yyyyyyyyyyyyyy/cat.m3u8


# STATUS
## Command Line parameters
#### v0.0.3
- --wallet (default wallet.json) specify different wallet file
- --allocation (default allocation.txt) specify different allocation file
#### v0.0.2
- --debug <1/0> (default 0) - 1 shows command line output to console
#### v0.0.1
- --config <configfile.yaml> (default config.yaml)
- --port <portno.> (default 6942)

## URL info paths
#### v0.0.1
- /status (Status OK)
- /config (Show currently selected config.yaml file contents)
- /version (show version number)

## UPDATES
### Fixes
#### v0.0.3
Improved file size handling
#### v0.0.2
Better error handling for incorrect paths
#### v0.0.1
Fixed differences in file/folder handling (folder still requires allocation even though encoded in authticket)
### Additions
#### v0.0.3
Added support for filehash for folder shares
#### v0.0.2
Added file size / xfer speed to console output
#### v0.0.1
Added support for default file for folder share. If no file specified, looks for .default file and then performs (301) redirect to file specified within

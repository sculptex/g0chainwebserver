# g0chainwebserver
Go based 0chain webserver
Serves files stored on 0chain from specified authticket
 
PRE-REQUSITES
Existing $HOME/.zcn folder (the location for this executable)
Containing zbox/zwallet cli tools plus existing wallet and allocation
Build and Install this executable in above folder

## USAGE

Open dedicated command line window in above folder

./g0chainwebserver <port> (default 6942)

browser (or app) point to:-

http://<IPaddress>:port/authticket/xxxxxxxxxxxxxx

for files

http://<IPaddress>:port/authticket/xxxxxxxxxxxxxx/remote/path/file.ext

for shared folders with specified full remote paths

e.g. http://192.168.1.50:6942/authticket/xxxxxxxxx/video/cat.m3u8

# Updates

v0.0.1

## Fixes

Fixed differences in file/folder handling (folder still requires allocation even though encoded in authticket)

## Additions

Added support for default file for folder share. If no file specified, looks for .default file and then performs (301) redirect to file specified within

## Added url info paths

- /status (Status OK)
- /config (Show currently selected config.yaml file contents)
- /version (show version number)

## Command line parameters

- --config <configfile.yaml> (default config.yaml)
- --port <portno.> (default 6942)

# g0chainwebserver
Go based 0chain webserver
Serves files stored on 0chain from specified authticket
 
PRE-REQUSITES
Existing $HOME/.zcn folder (the location for this executable)
Containing zbox/zwallet cli tools plus existing wallet and allocation
Build and Install this executable in above folder

USAGE
Open dedicated command line window in above folder
./g0chainwebserver <port> (default 6942)
browser (or app) point to:-
http://<IPaddress>:port/authticket/xxxxxxxxxxxxxx
for files
http://<IPaddress>:port/authticket/xxxxxxxxxxxxxx/remote/path/file.ext
for shared folders with specified full remote paths
e.g. http://192.168.1.50:6942/authticket/xxxxxxxxx/video/cat.m3u8

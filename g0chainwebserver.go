// Go based 0chain webserver
// Serves files stored on 0chain from specified authticket
// by Sculptex
// 
// PRE-REQUSITES
// Existing $HOME/.zcn folder (the location for this executable)
// Containing zbox/zwallet cli tools plus existing wallet and allocation
// Build and Install this executable in above folder
//
// USAGE
// Open dedicated command line window in above folder
// ./g0chainwebserver <port> (default 6942)
// browser (or app) point to:-
// http://<IPaddress>:port/authticket/xxxxxxxxxxxxxx
// for files
// http://<IPaddress>:port/authticket/xxxxxxxxxxxxxx/remote/path/file.ext
// for shared folders with specified full remote paths
// e.g. http://192.168.1.50:6942/authticket/xxxxxxxxx/video/cat.m3u8

package main

import (
    "log"
    "os"
    "net/http"
    "flag"
    "fmt"
    //"html"
    "path"
    "time"
    "os/exec"
    "strings"
    "io/ioutil"
    "encoding/json"
    "encoding/base64"
)

const tmppath="tmp"

func FilenameWithoutExtension(fn string) string {
      return strings.TrimSuffix(fn, path.Ext(fn))
}

func microTime() float64 {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	micSeconds := float64(now.Nanosecond()) / 1000000000
	return float64(now.Unix()) + micSeconds
}

func getTmpPath() string {
	tmp := fmt.Sprintf("%d", time.Now().UnixNano())
    return (tmp)
}
	
func main() {

    http.HandleFunc("/config/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "config.yaml")
    })

    http.HandleFunc("/authticket/", func(w http.ResponseWriter, r *http.Request) {
		urlfull := r.URL.Path
		urlonly := strings.Replace(urlfull, "/authticket/", "", 1)
		// Split into Max two parts, authticket plus (optional) path
		urlsplit := strings.SplitN(urlonly, "/", 2)
		authticket := urlsplit[0];
		// Convert to json 
		atjson, err := base64.StdEncoding.DecodeString(authticket)
		if err != nil {
			fmt.Println("decode error:", err)
			return
		}
		
		var dat map[string]interface{}
		if err := json.Unmarshal(atjson, &dat); err != nil {
			panic(err)
		}
		// DEBUG
		//fmt.Println(dat)
		
		referencetype := dat["reference_type"].(string)
		var filename string
				
		allocationid, aerr := ioutil.ReadFile("allocation.txt")
		if aerr != nil {
			log.Fatal(aerr)
		}
		tmpdir := tmppath
		if len(tmpdir)>0 {
			tmpdir = tmpdir+string(os.PathSeparator)
		}
		tmpdir = tmpdir+getTmpPath()
		os.Mkdir(tmpdir, 0700)
		
		if referencetype == "d" {
			// If Authticket is for a directory, then file path is extracted as rest of url
		    filename = "/"+urlsplit[1]
		}
		
		if referencetype == "f" {
			// If Authticket is for a directory, then file path is extracted from authticket
			filename = dat["file_name"].(string)
			if(len(filename)==0) {
				// Failing that, file path is extracted as rest of url
				filename = "/"+urlsplit[1]
			}
		}
		
		if(len(filename)>0) {
			fmt.Println("\nServing "+filename)
			relfilename := tmpdir+filename
			cmdarray := []string{
				"./zbox",
				"download",
				"--allocation",
				string(allocationid),
				"--authticket",
				string(authticket),
				"--remotepath",
				string(filename),
				"--localpath",
				string(relfilename) }
			// DEBUG
			//output := strings.Join(cmdarray, " ")
			//fmt.Printf("OUTPUT, %s\n", output)
			head := cmdarray[0]
			parts := cmdarray[1:len(cmdarray)]
			_ , err = exec.Command(head,parts...).Output()
			if err != nil {
				fmt.Printf("%s", err)
			}
			// Server Downloaded File
			// opportunity for file validation
			http.ServeFile(w, r, relfilename)
		}	
        // (opportunity to cache files)
        // Delete folder and file
        os.RemoveAll(tmpdir) 
    })

    // Allow user to specify port    
    var port string
    flag.StringVar(&port, "port", "6942", "Port Number (default 6942)")
    flag.Parse()

    // Advise listening on port
    fmt.Println("Listening on port:", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

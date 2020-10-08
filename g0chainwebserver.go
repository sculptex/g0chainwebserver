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
//
// AUTHTICKET FILE SHARE
// For any wallet (authticket contains file hash)
// http://<IPaddress>:port/authticket/xxxx
//
// AUTHTICKET FOLDER SHARE
// For Owner Wallet (full path required)
// http://<IPaddress>:port/authticket/xxxx/remote/path/file.ext
// e.g. http://192.168.1.50:6942/authticket/xxxx/video/cat.m3u8
//
// For Other Wallet (File Hash required) xxxx authticket yyyy filehash
// http://<IPaddress>:port/authhash/xxxx/yyyy_file.ext
// e.g. http://192.168.1.50:6942/authhash/xxxx/yyyy_index.html
// (file.ext required else hash is used for filename with no extension)
// Reason for change to _ separator, file path will still be relative
// to authticket

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

const tmppath = "tmp"
const defaultconfig = "config.yaml"
const defaultwallet = "wallet.json"
const defaultallocation = "allocation.txt"
const version = "0.0.4"

var configfile string
var allocationfile string
var walletfile string
var debug string

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func showfilesize(bytes int64) string {
	if(bytes < 1000) {
		return(fmt.Sprintf("%d bytes", bytes)) 
	}
	fbytes := float64(bytes)
	if(bytes < 1000000) {
		return(fmt.Sprintf("%0.1f KB", float64(fbytes/1000))) 
	}
	if(bytes < 1000000000) {
		return(fmt.Sprintf("%0.2f MB", float64(fbytes/1000000))) 
	}			
	return(fmt.Sprintf("%0.3f GB", float64(fbytes/1000000000))) 
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

    http.HandleFunc("/version/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Version %s", version)
    })
    	
    http.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK, config = %s, wallet = %s, allocation = %s", configfile, walletfile, allocationfile)
    })	

    http.HandleFunc("/config/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, defaultconfig)
    })

    http.HandleFunc("/enckeys/", func(w http.ResponseWriter, r *http.Request) {
		var cmdarray []string
		cmdarray = []string{
			"./zbox",
			"getwallet",
			"--json" }

		head := cmdarray[0]
		parts := cmdarray[1:len(cmdarray)]
		
		jsonres , err := exec.Command(head,parts...).Output()
		
		if err != nil {
			fmt.Printf("%s", err)
		}
		
		var dat map[string]interface{}
		if err := json.Unmarshal(jsonres, &dat); err != nil {
			panic(err)
		}
		
		clientid := dat["client_id"].(string)
		encpubkey := dat["encryption_public_key"].(string)
								
		fmt.Fprintf(w, "OK, client_id = %s, encryption_public_key = %s", clientid, encpubkey)
    })
    
    mainhandle := func(w http.ResponseWriter, r *http.Request) {
		urlfull := r.URL.Path
		urlhost := r.Host
		urlonly := urlfull
		var urlsplit []string
		isauthticket := strings.Index(urlfull, "/authticket/")==0
		if(isauthticket) {
		  urlonly = strings.Replace(urlonly, "/authticket/", "", 1)
		  // Split into Max two parts, authticket plus (optional) filepath
		  urlsplit = strings.SplitN(urlonly, "/", 2)
		} else
		{
		  urlonly = strings.Replace(urlonly, "/authhash/", "", 1)
		  // Split into Max three parts, authticket, hash plus (optional) filename
		  urlsplit = strings.SplitN(urlonly, "/", 2)
		}

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
		var cmdarray []string
		var relfilename string
		var redirect bool
		redirect = false
						
		allocationid, aerr := ioutil.ReadFile(allocationfile)
		if aerr != nil {
			log.Fatal(aerr)
		}
		tmpdir := tmppath
		if len(tmpdir)>0 {
			tmpdir = tmpdir+string(os.PathSeparator)
		}
		tmpdir = tmpdir+getTmpPath()//+string(os.PathSeparator)
		os.Mkdir(tmpdir, 0700)

		lookuphash := ""

		if referencetype == "d" {
		    filename = "/"
			if(isauthticket) {
				// If Authticket is for a directory, then file path is extracted as rest of url
				if(filename[len(filename)-1:] == "/") {
					redirect = true
					filename = filename+".default"
				}
		    } else
		    {
				// authhash, hash is 2nd parameter
				lookuphash = urlsplit[1]
				hashsplit := strings.SplitN(lookuphash, "_", 2)
				lookuphash = hashsplit[0]
				if(len(hashsplit)>1) {
					// use after _ as filename
					filename = filename+hashsplit[1]
				} else
				{
					// just use filehash as filename if not specified 
					filename = filename+hashsplit[0]
				}
			}
		    relfilename = tmpdir+filename
			cmdarray = []string{
				"./zbox",
				"download",
				"--allocation",  // existing bug in zboxcli, allocation required for folder event though included in authticket
				string(allocationid),		
				"--authticket",	
				string(authticket),
				"--localpath",
				string(relfilename) }
			if(len(lookuphash) > 0) {
				cmdarray = append( cmdarray, "--lookuphash", lookuphash )	
			} else
			{
				cmdarray = append( cmdarray, "--remotepath", string(filename) )	
			}
			    
		}
		
		if referencetype == "f" {
			// If Authticket is for a file, then file path is extracted from authticket
			filename = dat["file_name"].(string)
			relfilename = tmpdir+string(os.PathSeparator)+filename
			redirect = false
			cmdarray = []string{
				"./zbox",
				"download",
				"--authticket",
				string(authticket),
				"--localpath",
				string(relfilename) }			
		}
		
		if(configfile != defaultconfig) {
			cmdarray = append( cmdarray, "--config", configfile )	
		}
		
		if(walletfile != defaultwallet) {
			cmdarray = append( cmdarray, "--wallet", walletfile )	
		}					
						
		if(len(filename)>0) {
			fmt.Printf("SRV: %s  (%s)\n", filename, urlhost)

			if(debug=="1") {
				output := strings.Join(cmdarray, " ")
				fmt.Printf("CMD: %s\n", output)
		    }
			head := cmdarray[0]
			parts := cmdarray[1:len(cmdarray)]
			
			var starttime float64
			var endtime float64
			var elapsedtime float64
			var rate int64
			var fsize int64

			starttime = microTime()
			_ , err = exec.Command(head,parts...).Output()
			endtime = microTime()
			elapsedtime = endtime-starttime
			
			if err != nil {
				fmt.Printf("%s", err)
			}
			// Server Downloaded File
			// opportunity for file validation
			if(fileExists(relfilename)) {
				if(redirect) {
				
					b, ferr := ioutil.ReadFile(relfilename) // just pass the file name
				    if ferr != nil {
				        fmt.Print(ferr)
				    }
				    newpath := string(b) // convert content to a 'string'
					fmt.Printf("RED: %s\n", newpath)
					http.Redirect(w, r, newpath, 301)
				} else
				{
					fi, serr := os.Stat(relfilename)
					if serr != nil {
					  	fmt.Printf("%s", serr)
					}
					fsize = fi.Size()
					rate = int64( float64(fsize) / elapsedtime )
					fmt.Printf("GET: %s  %s in %0.3f secs (%s/sec)\n", filename, showfilesize(fsize), elapsedtime, showfilesize(rate))
					
					starttime = microTime()
					http.ServeFile(w, r, relfilename)
					endtime = microTime()
					elapsedtime = endtime-starttime
					rate = int64( float64(fsize) / elapsedtime )
					fmt.Printf("SND: %s  %s in %0.3f secs (%s/sec)\n", filename, showfilesize(fsize), elapsedtime, showfilesize(rate))
	
				}
			} else
			{
				fmt.Printf("ERR: %s  NOT FOUND\n", filename)
				http.NotFound(w, r)
			}
			
		}	
        // (opportunity to cache files)
        // Delete folder and file
        os.RemoveAll(tmpdir) 
    }
    
    http.HandleFunc("/authticket/", mainhandle)
    http.HandleFunc("/authhash/", mainhandle)

    // Allow user to specify port    
    var port string
    flag.StringVar(&port, "port", "6942", "Port Number (default 6942)")
    
    // Allow user to specify config file    
    flag.StringVar(&configfile, "config", string(defaultconfig), "config file (default "+defaultconfig+")")

    // Allow user to specify wallet file    
    flag.StringVar(&walletfile, "wallet", string(defaultwallet), "wallet file (default "+defaultwallet+")")
    
    // Allow user to specify allocation file    
    flag.StringVar(&allocationfile, "allocation", string(defaultallocation), "allocation file (default "+defaultallocation+")")
        
    // Allow debug    
    flag.StringVar(&debug, "debug", "0", "debug (1 or 0, default 0)")
    
    flag.Parse()

    // Advise listening on port
    fmt.Println("Listening on port:", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

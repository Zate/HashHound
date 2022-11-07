package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/jasonlvhit/gocron"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var HashMap = make(map[string]int)

func errCheck(msg string, err error) {
	if err != nil {
		log.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

func hashtags(c echo.Context) error {

	payload := `<!doctype html>
	<html lang="en">
	  <head>
		<!-- Required meta tags -->
		<link rel="apple-touch-icon" sizes="57x57" href="/apple-icon-57x57.png">
		<link rel="apple-touch-icon" sizes="60x60" href="/apple-icon-60x60.png">
		<link rel="apple-touch-icon" sizes="72x72" href="/apple-icon-72x72.png">
		<link rel="apple-touch-icon" sizes="76x76" href="/apple-icon-76x76.png">
		<link rel="apple-touch-icon" sizes="114x114" href="/apple-icon-114x114.png">
		<link rel="apple-touch-icon" sizes="120x120" href="/apple-icon-120x120.png">
		<link rel="apple-touch-icon" sizes="144x144" href="/apple-icon-144x144.png">
		<link rel="apple-touch-icon" sizes="152x152" href="/apple-icon-152x152.png">
		<link rel="apple-touch-icon" sizes="180x180" href="/apple-icon-180x180.png">
		<link rel="icon" type="image/png" sizes="192x192"  href="/android-icon-192x192.png">
		<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
		<link rel="icon" type="image/png" sizes="96x96" href="/favicon-96x96.png">
		<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
		<link rel="manifest" href="/manifest.json">
		<meta name="msapplication-TileColor" content="#ffffff">
		<meta name="msapplication-TileImage" content="/ms-icon-144x144.png">
		<meta name="theme-color" content="#ffffff">
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		
		<title>hashtags.fyi</title>
		<style>
		body {text-align: center;}
		.wrapper {display: inline-block;margin-top: 25px;position: relative;}
		
		.wrapper img {
			display: block;
			max-width: 100%;
		}
		
		.wrapper .overlay {
			position: absolute;
			top: 70%;
			left: 38%;
			transform: translate(-50%, -50%);
			color: black;
		}
		</style>
	  </head> 
	  <body>
	  <div class="wrapper">
		<div class="overlay"><h2>`
	payload += getHashMap()
	payload += `</h2></div>
	</div>
	<a rel="me" href="https://infosec.exchange/@zate"></a>
</body>
</html>`
	return c.HTMLBlob(http.StatusOK, []byte(payload))
}

func getHashMap() string {
	s, err := os.ReadFile(".hashmap")
	errCheck("Not able to read .hashmap", err)
	return string(s)
}

func updateHashTagFile() {
	for k := range HashMap {
		delete(HashMap, k)
	}
	getMastodonHashTags()
	s := ""
	keys := make([]string, 0, len(HashMap))

	for key := range HashMap {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return HashMap[keys[i]] > HashMap[keys[j]]
	})
	for _, k := range keys {
		s += "#" + k + " => " + fmt.Sprint(HashMap[k]) + "</br>\n"
	}
	f := []byte(s)
	err := os.WriteFile(".hashmap", f, 0644)
	errCheck("Problem Writing to file", err)
	log.Println("This task will run periodically")
}

func executeCronJob() {
	gocron.Every(15).Minute().Do(updateHashTagFile)
	<-gocron.Start()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	updateHashTagFile()
	go executeCronJob()
	e := echo.New()
	e.File("/favicon.ico", "favicon.ico")
	e.File("/favicon-16x16.png", "favicon-16x16.png")
	e.File("/favicon-32x32.png", "favicon-32x32.png")
	e.File("/favicon-96x96.png", "favicon-96x96.png")
	e.File("/manifest.json", "manifest.json")
	e.File("/ms-icon-70x70.png", "ms-icon-70x70.png")
	e.File("/ms-icon-144x144.png", "ms-icon-144x144.png")
	e.File("/ms-icon-150x150.png", "ms-icon-150x150.png")
	e.File("/ms-icon-310x310.png", "ms-icon-310x310.png")
	e.File("/android-icon-36x36.png", "android-icon-36x36.png")
	e.File("/android-icon-48x48.png", "android-icon-48x48.png")
	e.File("/android-icon-72x72.png", "android-icon-72x72.png")
	e.File("/android-icon-96x96.png", "android-icon-96x96.png")
	e.File("/android-icon-144x144.png", "android-icon-144x144.png")
	e.File("/android-icon-192x192.png", "android-icon-192x192.png")
	e.File("/apple-icon-57x57.png", "apple-icon-57x57.png")
	e.File("/apple-icon-60x60.png", "apple-icon-60x60.png")
	e.File("/apple-icon-72x72.png", "apple-icon-72x72.png")
	e.File("/apple-icon-76x76.png", "apple-icon-76x76.png")
	e.File("/apple-icon-114x114.png", "apple-icon-114x114.png")
	e.File("/apple-icon-120x120.png", "apple-icon-120x120.png")
	e.File("/apple-icon-144x144.png", "apple-icon-144x144.png")
	e.File("/apple-icon-152x152.png", "apple-icon-152x152.png")
	e.File("/apple-icon-180x180.png", "apple-icon-180x180.png")
	e.File("/apple-icon-precomposed.png", "apple-icon-precomposed.png")
	e.File("/apple-icon.png", "apple-icon.png")
	e.File("/browserconfig.xml", "browserconfig.xml")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET},
	}))
	e.Use(middleware.Secure())
	e.GET("/", hashtags)
	e.Logger.Fatal(e.Start(":3003"))
}

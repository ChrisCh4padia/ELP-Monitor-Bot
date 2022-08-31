package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

//go:embed output.txt
var content string

//go:embed counter.txt
var countercontent string

// Twittercreds
var TwconsumerKey string = ""
var TwconsumerSecret string = ""
var TwaccessToken string = ""
var TwaccessTokenSecret string = ""

// vars
var outstr string
var out []byte
var err error
var txtdiff string
var Tweetcontent string
var counter int = 1
var counterstr string
var countercontentint int

func ClearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		fmt.Println(err, "1")
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			fmt.Println(err, "2")
		}
	}
	return nil
}

func getcontent() {

	file, err := os.Open("/etc/ELPmon/output.txt")
	if err != nil {
		fmt.Println(err, "3")
	}
	contentcon, err := ioutil.ReadFile("/etc/ELPmon/output.txt")
	if err != nil {
		fmt.Println("Err")
	}

	content = string(contentcon)
	defer file.Close()
}

func compare() {

	outstr = string(out)

	if outstr == content {
	} else {
		getdiff()
	}

}

func callHR() {

	cmd := exec.Command("/bin/bash", "/etc/ELPmon/handelsregister-main/run.sh")
	out, err = cmd.Output()

	if err != nil {
		fmt.Println(string(out))
		println(err.Error())
		return

	}
}

func getdiff() {
	txtdiff = (strings.Replace(outstr, content, "", 1))
}

func removeleftover() {

	rm := ClearDir("/etc/ELPmon/cache/")
	if rm != nil {
		fmt.Println(rm, "4")
	}

}

func getcountercontent() {

	file, err := os.Open("/etc/ELPmon/counter.txt")
	if err != nil {
		fmt.Println(err, "5")
	}
	content2, err := ioutil.ReadFile("/etc/ELPmon/counter.txt") // the file is inside the local directory
	if err != nil {
		fmt.Println("Err")
	}
	countercontent = string(content2)
	countercontentint, _ = strconv.Atoi(countercontent)
	counter = countercontentint
	defer file.Close()

}

func setcounter() {

	s := countercontentint + 1
	sstring := strconv.Itoa(s)
	sbyte := []byte(sstring)
	ioutil.WriteFile("/etc/ELPmon/counter.txt", sbyte, 0777)
}

func setcontent() {
	contentbyte := []byte(outstr)
	ioutil.WriteFile("/etc/ELPmon/output.txt", contentbyte, 0777)
}

func Tweetoutput() {
	counterstr = strconv.Itoa(counter)
	if txtdiff == "" {
		if counter == 1 {
			Tweetcontent = "No changes in the last day"
			counter = countercontentint + 1
		} else {
			Tweetcontent = "No changes in the last " + counterstr + " days"
			counter = countercontentint + 1
		}
	} else {
		Tweetcontent = "New change in handelsregister:" + "\n" + outstr
		countercontentint = 0
	}

	config := oauth1.NewConfig(TwconsumerKey, TwconsumerSecret)
	token := oauth1.NewToken(TwaccessToken, TwaccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	user, _, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("Account: @%s (%s)\n", user.ScreenName, user.Name)
	_, resp, err := client.Statuses.Update(Tweetcontent, nil)
	if err != nil {
		fmt.Print(err)
		fmt.Println(resp)
	}
}

func main() {

	for {
		getcontent()
		getcountercontent()
		callHR()
		setcontent()
		compare()
		getdiff()
		removeleftover()
		Tweetoutput()
		setcounter()
		setcontent()
		txtdiff = ""
		time.Sleep(24 * time.Hour)
	}

}

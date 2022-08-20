package main

import (
	_ "embed"
	"fmt"
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

func ClearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func compare() {

	outstr = string(out)

	if outstr == content {
	} else {
		getdiff()
	}

}

func callHR() {

	cmd := exec.Command("bash", "-c", "./handelsregister-main/run.sh")
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

	rm := ClearDir("./cache/")
	if rm != nil {
		fmt.Println(rm)
	}

}

func Tweetoutput() {
	counterstr = strconv.Itoa(counter)
	if txtdiff == "" {
		if counter == 1 {
			Tweetcontent = "No changes in the last day"
			counter = counter + 1
		} else {
			Tweetcontent = "No changes in the last " + counterstr + " days"
			counter = counter + 1
		}
	} else {
		Tweetcontent = "New change in handelsregister:" + "\n" + txtdiff
		counter = 1
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
		callHR()
		compare()
		getdiff()
		removeleftover()
		Tweetoutput()
		txtdiff = ""
		time.Sleep(24 * time.Hour)
	}
}

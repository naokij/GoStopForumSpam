package main

import (
	"flag"
	"fmt"
	"github.com/naokij/GoStopForumSpam/stopforumspam"
)

var (
	key      *string = flag.String("key", "", "stopforumspam.com API key")
	ip       *string = flag.String("ip", "", "ip")
	email    *string = flag.String("email", "", "email")
	username *string = flag.String("username", "", "username")
	evidence *string = flag.String("evidence", "", "evidence")
	c        *stopforumspam.Client
)

func main() {
	flag.Parse()
	c = stopforumspam.New(*key)
	searchData := stopforumspam.SearchData{Ip: *ip, Email: *email, Username: *username}
	spamData := stopforumspam.SpamData{SearchData: searchData, Evidence: *evidence}
	if err := c.Add(spamData); err != nil {
		fmt.Println("Add error: ", err)
	} else {
		fmt.Print("successful")
	}
	response, err := c.Search(searchData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.Ip.LastSeen.ToStdTime(), response.Email.LastSeen.ToStdTime(), response.Username)
}

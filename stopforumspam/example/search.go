package main

import (
	"flag"
	"fmt"
	"github.com/naokij/GoStopForumSpam/stopforumspam"
)

var (
	key      *string = flag.String("key", "", "stopforumspam.com API key")
	ip       *string = flag.String("ip", "", "ip to be searched")
	email    *string = flag.String("email", "", "email to be searched")
	username *string = flag.String("username", "", "username to be searched")
	c        *stopforumspam.Client
)

func main() {
	flag.Parse()
	c = stopforumspam.New(*key)
	searchData := stopforumspam.SearchData{Ip: *ip, Email: *email, Username: *username}
	response, err := c.Search(searchData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.Success, response.Ip, response.Email, response.Username, response.Error)
}

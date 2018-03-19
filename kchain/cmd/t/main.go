package main

import (
	dbm "github.com/tendermint/tmlibs/db"
	"fmt"
)

func main() {

	db, err := dbm.NewGoLevelDB("hello", "hello")
	if err != nil {
		fmt.Println(err.Error())
	}

	db.Set([]byte("hello:1"), []byte("1"))
	db.Set([]byte("hello:2"), []byte("1"))
	db.Set([]byte("hello:3"), []byte("1"))
	db.Set([]byte("hello:0"), []byte("1"))
	db.Set([]byte("hello:2"), []byte("1"))
	db.Set([]byte("hello:9"), []byte("1"))
	db.Set([]byte("tello:4"), []byte("1"))
	db.Set([]byte("3ello:4"), []byte("1"))
	db.Set([]byte("oello:4"), []byte("1"))
	db.Set([]byte("aello:4"), []byte("1"))
	db.Set([]byte("height:4"), []byte("1nuhuhuhu"))

	i := db.Iterator([]byte("hello:3"), []byte("hello:999999999999999999999999999"))
	for {
		if i.Valid() {
			fmt.Println(string(i.Key()))
			i.Next()
		}else{
			break
		}

	}
}

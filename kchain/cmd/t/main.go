package main

import (
	//dbm "github.com/tendermint/tmlibs/db"
	"github.com/json-iterator/go"
	"fmt"
	"encoding/hex"
)

func main() {

	//db, err := dbm.NewGoLevelDB("hello", "hello")
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//db.Set([]byte("hello:1"), []byte("1"))
	//db.Set([]byte("hello:2"), []byte("1"))
	//db.Set([]byte("hello:3"), []byte("1"))
	//db.Set([]byte("hello:0"), []byte("1"))
	//db.Set([]byte("hello:2"), []byte("1"))
	//db.Set([]byte("hello:9"), []byte("1"))
	//db.Set([]byte("tello:4"), []byte("1"))
	//db.Set([]byte("3ello:4"), []byte("1"))
	//db.Set([]byte("oello:4"), []byte("1"))
	//db.Set([]byte("aello:4"), []byte("1"))
	//db.Set([]byte("height:4"), []byte("1nuhuhuhu"))
	//
	//i := db.Iterator([]byte("hello:3"), []byte("hello:999999999999999999999999999"))
	//for {
	//	if i.Valid() {
	//		fmt.Println(string(i.Key()))
	//		i.Next()
	//	}else{
	//		break
	//	}
	//
	//}

	a1 := map[string]string{"sss": "4", "ss1s": "4", "ss2s": "4"}

	a, _ := jsoniter.MarshalToString(a1)
	fmt.Println(a)

	a3, _ := jsoniter.Marshal(a1)
	fmt.Println(hex.EncodeToString(a3))

	aa := "WyJhY2NvdW50OmIzYzUzZDljZTVhZjEyMTlkYTY5ZDBiMTI1M2RmZjk2NDdlZjM2NTYiLCJtZXRhZGF0YToyRUJGSlQyUURYUTZEN1gwOTVXQjNYSUdPUk42TERTWk5PRlVXUU5OSzA4TjJERlJJSSIsIm1ldGFkYXRhOjJFQkZKVDJRRFhRNkQ3WDA5NVdCM1hJR09STjZMRFNaTk9GVVdRTk5LMDhOMkRGUklJIiwibWV0YWRhdGE6MkVCRkpUMlFEWFE2RDdYMDk1V0IzWElHT1JONkxEU1pOT0ZVV1FOTkswOE4yREZSSUkiLCJtZXRhZGF0YToyRUJGSlQyUURYUTZEN1gwOTVXQjNYSUdPUk42TERTWk5PRlVXUU5OSzA4TjJERlJJSSIsIm1ldGFkYXRhOjJFQkZKVDJRRFhRNkQ3WDA5NVdCM1hJR09STjZMRFNaTk9GVVdRTk5LMDhOMkRGUklJIiwibWV0YWRhdGE6MkVCRkpUMlFEWFE2RDdYMDk1V0IzWElHT1JONkxEU1pOT0ZVV1FOTkswOE4yREZSSUkiLCJtZXRhZGF0YToyRUJGSlQyUURYUTZEN1gwOTVXQjNYSUdPUk42TERTWk5PRlVXUU5OSzA4TjJERlJJSSJd"
	dd, _ := hex.DecodeString(aa)

	dd1 := new([]string)
	jsoniter.Unmarshal(dd, &dd1)

	fmt.Println(dd)
	fmt.Println(dd1)

	a22, _ := jsoniter.Marshal("hello")
	fmt.Println(string(a22))

	var a34 interface{}
	jsoniter.Unmarshal(a22,&a34)
	fmt.Println(a34)

}

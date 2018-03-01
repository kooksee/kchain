package crypto

import (
	crypto "github.com/tendermint/go-crypto"
	"sort"
	"strings"
)

func MapHasher(a map[string]string) ([]byte, []string) {
	var (
		ks []string
		vs []string
	)

	for key, _ := range a {
		ks = append(ks, key)
	}

	sort.Strings(ks)

	for _, k := range ks {
		vs = append(vs, a[k])
	}

	return crypto.Ripemd160([]byte(strings.Join(vs, ""))), ks
}

func StringHasher(a string) []byte {
	return crypto.Ripemd160([]byte(a))
}

//func hashStringMap(m map[string]interface{}) []byte {
//	hash := sha3.New512()
//	encoder := jsoniter.NewEncoder(hash)
//	keys := make([]string, len(m))
//	i := 0
//	for id := range m {
//		keys[i] = id
//		i++
//	}
//	sort.Strings(keys)
//	for _, key := range keys {
//		encoder.Encode(key)
//		encoder.Encode(m[key])
//	}
//	return hash.Sum(nil)
//}
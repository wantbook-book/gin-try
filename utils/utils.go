package utils

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {

	strarr := []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	s := make([]byte, n)
	s1 := rand.NewSource(time.Now().Unix())
	r1 := rand.New(s1)
	for i := 0; i < n; i++ {
		s[i] = strarr[r1.Intn(len(strarr))]
	}
	return string(s)
}

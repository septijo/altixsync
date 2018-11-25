package main

import (
	"fmt"
	"crypto/sha256"
)

// https://www.fileformat.info/tool/hash.htm
func npSha256_salted(key1 []byte, key2 string) []byte {
  /* fmt.Println(sha256.Sum256([]byte("1234"))) // 03 AC 67 42 16 F3 E1 5C 76 1E E1 A5 E2 55 F0 67 95 36 23 C8 B3 88 B4 45 9E 13 F9 78 D7 C8 46 F4

  x := sha256.New()
  x.Write([]byte("1234"))
  fmt.Println(x.Sum(nil)) // 03 AC 67 42 16 F3 E1 5C 76 1E E1 A5 E2 55 F0 67 95 36 23 C8 B3 88 B4 45 9E 13 F9 78 D7 C8 46 F4

  y := sha256.New()
  y.Write([]byte("12"))
  fmt.Println(y.Sum([]byte("34"))) // hasil jadi 34 byte
  
  z := sha256.New()
  z.Write([]byte("12"))
  z.Write([]byte("34"))
  fmt.Println(z.Sum(nil)) // 03 AC 67 42 16 F3 E1 5C 76 1E E1 A5 E2 55 F0 67 95 36 23 C8 B3 88 B4 45 9E 13 F9 78 D7 C8 46 F4 */
  
  s := sha256.New()
  s.Write(key1)
  s.Write([]byte(key2))
  fmt.Println(key1, key2) 
  fmt.Println(s.Sum(nil)) // 03 AC 67 42 16 F3 E1 5C 76 1E E1 A5 E2 55 F0 67 95 36 23 C8 B3 88 B4 45 9E 13 F9 78 D7 C8 46 F4 */
  return s.Sum(nil)
}
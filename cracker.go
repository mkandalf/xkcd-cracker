package main

import (
    "fmt"
  "crypto/skein" //https://github.com/wernerd/Skein3Fish
  "encoding/hex"
  "math/rand"
  "strconv"
  "time"
  "net/http"
  "net/url"
  "bytes"
)

const goal string = "5b4da95f5fa08280fc9879df44f418c8f9f12ba424b7757de02bbdfbae0d4c4fdf9317c80cc5fe04c6429073466cf29706b8c25999ddd2f6540d4475cc977b87f4757be023f19b8f4035d7722886b78869826de916a79cf9c94cc79cd4347d24b567aa3e2390a573a373a48a5e676640c79cc70197e1c5e7f902fb53ca1858b6"

func hash (in string) []byte {
  var empty []byte
  skein, _ := skein.NewMac(1024, 1024, empty)
  skein.Update([]byte(in))
  res := skein.DoFinal()
  return res
}

func compBytes (a, b []byte) int {
  res := 0
  for i := 0; i < len(a); i++ {
    res += compByte(a[i], b[i])
  }
  return res
}

func compByte (a, b byte) int {
  res := 0
  for i := uint32(0); i < 8; i++ {
    if ((a >> i) & 1) != ((b >> i) & 1) {
      res += 1;
    }
  }
  return res
}

func run() {
  rand.Seed(time.Now().UTC().UnixNano())
  best := 1000
  goalbytes, _ := hex.DecodeString(goal)
  for i := 0; i < 1000000; i++ {
    str := strconv.FormatInt(rand.Int63(), 10)
    res := compBytes(hash(str), goalbytes)
    if res < best {
      fmt.Printf("%d '%s'\n", res, str)
      if res < 430 {
        http.Get("http://lewisjellis.webscript.io/skeinlog?bits=" + strconv.Itoa(res) + "&number=" + str)
        if res < 406 {
          resp, _ := http.PostForm("http://almamater.xkcd.com/?edu=seas.upenn.edu",
                        url.Values{"hashable": {str}})
          buf := new(bytes.Buffer)
          buf.ReadFrom(resp.Body)
          s := buf.String()
          fmt.Printf("%s\n", s)
          defer resp.Body.Close()
        }
      }
      best = res
    }
  }
}
func main() {
  for i := 0; i < 7; i++ {
    go run()
  }
  run()
}


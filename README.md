# Detect repeating content

```sh
$ xxd -g3 < ./testdata/foobarbaz.txt
00000000: 666f6f 626172 626172 626172 626172 62  foobarbarbarbarb
00000010: 617262 617262 617262 617262 617a0a     arbarbarbarbaz.

$ go run main.go -min-size=3 ./testdata/foobarbaz.txt
length: 31, offset: 3, repeats: 3*8=24
prefix: 666f6f
repeat: 626172
suffix: 62617a0a
```

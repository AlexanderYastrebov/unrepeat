# Detect repeating content

```sh
$ xxd -g3 < ./testdata/foobarbaz.txt
00000000: 666f6f 626172 626172 626172 626172 62  foobarbarbarbarb
00000010: 617262 617262 617262 617262 617a0a     arbarbarbarbaz.

$ go run main.go -min-size=3 ./testdata/foobarbaz.txt
2023/02/18 01:00:37 len: 31, offset: 3, repeats: 3*8=24
2023/02/18 01:00:37 prefix: 666f6f
2023/02/18 01:00:37 repeat: 626172
2023/02/18 01:00:37 suffix: 62617a0a
```

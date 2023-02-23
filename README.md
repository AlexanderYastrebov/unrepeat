# Detect repeating content

```sh
$ xxd -g3 -c15 < ./testdata/foobarbaz.txt
00000000: 666f6f 626172 626172 626172 626172  foobarbarbarbar
0000000f: 626172 626172 626172 626172 62617a  barbarbarbarbaz

$ go run main.go -min-size=3 ./testdata/foobarbaz.txt
length: 30, offset: 3, repeats: 3*8=24
prefix: 666f6f
repeat: 626172
suffix: 62617a
```

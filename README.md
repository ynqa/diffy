# diffy

*diffy* prints diffs like GitHub pull requests (e.g. unified, split), and alternative to [diff](https://linuxjm.osdn.jp/html/gnumaniak/man1/diff.1.html).

<img src="https://user-images.githubusercontent.com/6745370/81488173-5cf1e200-92a0-11ea-91ab-6f05f93aef94.png" width="50%" height="50%"/>
<img src=https://user-images.githubusercontent.com/6745370/81488175-5ebba580-92a0-11ea-8721-17112cde89ec.png width="100%" height="100%">

## Installation

For MacOS:

```bash
$ brew tap ynqa/tap-archived
$ brew install diffy
```

From source codes:

```bash
$ go get -u github.com/ynqa/diffy
```

## Usage
```
Print colored diff more readable

Usage:
  diffy [flags] FILE1 FILE2

Flags:
  -c, --context int    number of context to print (default 3)
  -h, --help           help for diffy
      --no-header      no file name header
  -s, --style string   output style; one of unified|split (default "unified")
      --tab-size int   tab stop spacing (default 4)
  -v, --version        version for diffy
```

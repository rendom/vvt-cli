vvt cli
=======
Crossplatform client to post and get pastes from vvt.nu

```
Usage of ./vvt-cli:
  -o="": Output file (defaul is stdout)
```

####Example usage
```
// Post data from stdin and returns paste url
cat file.name | vvt

// Get paste from vvt
vvt iggqrfj

// Paste to file
vvt -o=outputfile.txt iggqrfj
```

####Todo
* Post encrypted paste
* Crossplatform clipboard support

# chainer


Library to generate binary stream from DSLed file

## Format of DSL

```

[+]<type>[s]: <value> [# comment]
```

```
<type>:
b: binary data <value> is hex string
d: decimal <value> is decimal number with base 10
s: string <value> is string
g: generated <value> is number of bytes to be generated randomly
```

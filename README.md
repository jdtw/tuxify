# tuxify

[Tuxify](https://words.filippo.io/the-ecb-penguin/) your jpeg and png images as a service!

```
curl -s -F 'img=@input.png' -o tuxified.png https://tuxify.art
```

Or install locally...

```
go install jdtw.dev/tuxify/cmd/tuxify@latest
tuxify --in input.png --out tuxified.png
```

![tux](https://github.com/jdtw/tuxify/blob/main/tuxified.png?raw=true)

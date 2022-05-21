# tuxify

Tuxify your jpeg and png images as a service!

```
curl -s -F 'img=@input.png' -o tuxified.png https://tuxify.art
```

Or build locally...

```
go build -o . ./...
./tuxify --in input.png --out tuxified.png
```

![tux](https://github.com/jdtw/tuxify/blob/main/tuxified.png?raw=true)

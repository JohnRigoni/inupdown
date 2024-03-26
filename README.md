
# inupdown

## Requires

golang >= v1.21

go install github.com/a-h/templ/cmd/templ@latest

npm install -D tailwindcss

## Building


```
templ generate
npx tailwindcss -i ./tailwind.css -o ./assets/tstyle.css
go build -o tmp/inupdown
```


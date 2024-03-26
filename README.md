
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

![pic](https://github.com/JohnRigoni/inupdown/assets/38547951/3852a66b-042d-4379-b1bd-a8bfde0843ed)

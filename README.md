pre-requisites:
```
docker
```

to run:
```
docker build -t app .
```
```
docker run -it --rm --name app -p 8080:8080 app.
```

open http://localhost:8080/front, open dev-tools and watch the session being created.

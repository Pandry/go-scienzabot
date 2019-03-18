# TODO
[] When sending the message, chech if the user exists in the gruop the list was called  
[] determine what admins can do and what not; eg. should it permit to ban someone?  
[] Find when to increment bookmark last update, to delete the bookmark after some time  
[] Implement subcategories (already implemented in DB)  
[] Refacor database/listHelpers.go  


#### Compilation for docker
Requires special flag
`CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-extldflags "-static"' .`

#### Sample Dockerfile
```
FROM alpine

ENV TELEGRAM_TOKEN 123:34556767867789

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true


COPY ./scienzabot /

CMD [ "/scienzabot", "-database", "/db/sqlite", "-vv" ]

```

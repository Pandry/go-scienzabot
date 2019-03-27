# TODO
- [ ] When sending the message, chech if the user exists in the gruop the list was called  
- [ ] determine what admins can do and what not; eg. should it permit to ban someone?  
- [ ] Find when to increment bookmark last update, to delete the bookmark after some time  
- [ ] Implement subcategories (already implemented in DB)  
- [ ] Refactor database/listHelpers.go methods  
- [ ] Improve user experience
- [ ] Set default group locale for the welcome message (or welcome message with no locale)
- [ ] Bookmarked messages should be seen from the latest one to the oldest.  

# ScienzaBot
This bot was made to tag the people when in a group where a lot of topics are treated.  
The bot look in every message for a "**list**".  
A **list** is basically a topic.  
The bot supports multiple lists.  
A list is called via a special character prepending the name of the list (`@`, `#`, `.` or `!`).  
When a list is called, the bot contacts via private message all the users **subscribed** to the list.  
A **subscription** is the relationship a user creates when he __joins__ a list.
If the group where the list was called is private, the bot will just say that a list was **invoked**, providing the possibility to see the message by tapping on a inline keyboard button to be tagged at the message:  

If instead the group is public, the bot will also provide 2 additional buttons; One that takes the user to the group, and the other one the the message in the group:  

The user also provide the possibility too "bookmark" a message.
Basically the message get saved and a user can see it later in time.

[To be completed in the future]

#### Compilation for Alpine to include all the libraries
Requires special flag
`CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-extldflags "-static"' .`

#### Sample Dockerfile (The used in production one)
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

# TODO
[ ] When sending the message, chech if the user exists in the gruop the list was called
[ ] determine what admins can do and what not; eg. should it permit to ban someone?
[ ] Check error messages logged ``` search "Requested a nickname not present in the database but the error is unknown" in the solution```
[ ] Add check for uniquie constraint on database (unsername table) - a user could have a nickname another used had in the past
[ ] Find when to increment bookmark last update, to delete the bookmark after some time
[ ] Implement subcategories (already implemented in DB)
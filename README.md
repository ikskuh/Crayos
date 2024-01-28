# Crayos!

A chaotic multiplayer crayon painting game.

## Build

```bash
# build and run backend
cd backend
go build
./crayos-backend
```


## Deployment

```sh-session
(cd backend && CGO_ENABLED=0 go build)
rsync -avhz --delete --delete-after "frontend/" phpfriends:/srv/crayos.random-projects.net 
scp backend/crayos-backend phpfriends:/opt/crayos.random-projects.net
```


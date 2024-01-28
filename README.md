```bash
# build and run backend
cd backend
go build
./crayos-backend

# open http://127.0.0.1:5500/frontend/?local in your browser



## Deployment

```sh-session
(cd backend && CGO_ENABLED=0 go build)
rsync -avhz --delete --delete-after "frontend/" phpfriends:/srv/crayos.random-projects.net 
scp backend/crayos-backend phpfriends:/opt/crayos.random-projects.net

```
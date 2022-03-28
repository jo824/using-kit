# https://ngrok.com/docs
# using to expose local drone.io server to the internet & test github integration/webhook.

docker run \
  --volume=/var/lib/drone:/data \
  --env=DRONE_GITHUB_CLIENT_ID=<$val> \
  --env=DRONE_GITHUB_CLIENT_SECRET=<$val2> \
  --env=DRONE_RPC_SECRET=<val3> \
  --env=DRONE_SERVER_HOST=90d1-69-126-160-199.ngrok.io \
  --env=DRONE_SERVER_PROTO=http \
  --publish=80:80 \
  --publish=443:443 \
  --restart=always \
  --detach=true \
  --name=drone \
  drone/drone:2

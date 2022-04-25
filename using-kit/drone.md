Drone.io

Drone is an open-source continuous integration and delivery platform based on container technology, encapsulating the environment and functionalities for each step of a build pipeline inside a temporary container. Thanks to the design, our builds and tests are decoupled from physical hosts and highly portable.

We'll look at setting up a single Drone server and pair it with a Drone runner. We'll integrate it with the example repository we've been working with on Github.


Let's start by looking at a high-level overview that'll introduce us to all the components involved.
Github integration
Drone server & runner
setting up each
Pipelines & Plugins
TODO: INSERT OVERVIEW DIAGRAM HERE

Drone Server
This server is responsible for authentication, repository configuration, users, secrets, and accepting webhooks.
Setup
From the documentation: The Drone server should be installed on a server or virtual machine (using your cloud provider of choice) with standard HTTP and HTTPS ports open. The instance must be publicly accessible by domain name or IP address to receive webhooks from GitHub.

Note: If you don't have access to a remote VM on a cloud provider, you can still test this locally by installing docker on your machine and starting the server/runner containers locally. Since the Drone server needs to be publicly accessible by GitHub(or another source control management tool), you can [use a tunneling tool like Ngrok](https://ngrok.com/docs/getting-started) to make this happen.

For my setup, I've:
Launched an instance on AWS via EC2 and installed Docker CE.
Made an elastic IP and assigned it to this instance.
Using Route53, I made a DNS entry pointing an existing domain to this IP.
As part of the GitHub integration setup, we need to register an OAuth app under developer settings in our GitHub account. We add the Homepage URL as `https://drone.yikyak.org` and the Authorization callback URL as`https://drone.yikyak.org/login` . We'll need the `Client ID` and `Client Secret.` for later when we start our drone server.

TODO: Add a screen for registration or link to the settings page.

One more piece we'll need is a shared secret. We'll use this secret to authenticate comms between our runner and the central Drone server. We can generate this shared secret from the docs like:
`openssl rand -hex 16`.

TODO: Phrase this side note better or remove it altogether. Ideally, it would read like it's nice to know the laborious setup and all components involved, even if you'd never actually do this.

<We could deploy to a service other than EC2 that would better suit running containers for the following steps. That's not the purpose of this example, though. We want to highlight the features of Drone and what workflows could be.>

Next, we will connect to our EC2 instance via SSH, pull in the Drone server container image, and run it.

Note: Since we didn't walk through the EC2 instance creation and docker install, it's worth mentioning a few things.
1) We set up a new security group, deleting the default launch options. It restricts access to the machine via SSH/HTTP to our IP.
2) Created an RSA key pair as a `.pem` format used to connect to the machine via SSH.

Connect to your instance by running
``` ssh -i ~/path/to/your/key ec2-user@Public IPv4 address/or Public IPv4 DNS
```
Let's pull the server image on this EC2 instance by running: `docker pull drone/drone:2`

I like to put my container run command into a `.sh` file. IMO it's easier to run later or edit over time. Here's what my initial run command looks like
```
drone_server.sh

docker run \
--volume=/var/lib/drone:/data \
--env=DRONE_GITHUB_CLIENT_ID<GH_CLIENT_ID> \
    --env=DRONE_GITHUB_CLIENT_SECRET=<YOUR_GITHUB_CLIENT_SECRET> \
        --env=DRONE_RPC_SECRET=<YOUR_SECRET> \
            --env=DRONE_SERVER_HOST=drone.yikyak.org \
            --env=DRONE_SERVER_PROTO=https \
            --env=DRONE_USER_CREATE=username:jo824,admin:true \
            --env=DRONE_REGISTRATION_CLOSED=true \
            --env=DRONE_TLS_AUTOCERT=true \
            --env=DRONE_LOGS_TEXT=true \
            --env=DRONE_LOGS_PRETTY=true \
            --env=DRONE_LOGS_COLOR=true \
            --env=DRONE_DEBUG=true \
            --publish=80:80 \
            --publish=443:443 \
            --restart=always \
            --detach=true \
            --name=drone \
            drone/drone:2
            ```

Above are many configuration options for our GitHub/Docker Drone server setup. You can find [the rest here](https://docs.drone.io/server/reference/).

TODO: Provide more context.
```
DRONE_USER_CREATE=username:jo824,admin:true
DRONE_REGISTRATION_CLOSED=true
DRONE_TLS_AUTOCERT=true
```
We create a Drone user. We explicitly give login access via my own GitHub account and provide it with admin privileges. We also close the registration of new users. Lastly, by setting `TLS_AUTOCERT` to true, we enable automated SSL configuration and updates using Let's Encrypt.

TODO: format these snippets/outputs.
After starting the container, we can confirm that it's up and running using `docker ps` command and should see output like:
```
CONTAINER ID  IMAGE             COMMAND         CREATED    STATUS   PORTS                                   NAMES
04a271f28534  drone/drone:2         "/bin/drone-server"   6 days ago  Up 6 days  0.0.0.0:80->80/tcp, :::80->80/tcp, 0.0.0.0:443->443/tcp, :::443->443/tcp  drone
```
Using the container's name, we can see the container's logs.
```
root@hostIP:~/scripts# docker logs drone
INFO[0000] starting the cron scheduler          interval=30m0s
INFO[0000] starting the http server           acme=true host=drone.yikyak.org port=":443" proto=https url="https://drone.yikyak.org"
INFO[0000] starting the zombie build reaper       interval=24h0m0s
```
The Drone server is up and running. We can use the host URL we defined to access the Done UI and sign in via our Github credentials.

TODO: insert the drone landing page + dashboard after the login window captures.


Drone Runner
A runner is a standalone daemon that polls the server for pending pipelines to execute.
Setup
We are going to ssh into the same ec2 host again and pull the runner packaged into a container provided by Drone. ` docker pull drone/drone-runner-docker:1`.
```
drone_runner.sh

docker run --detach \
--volume=/var/run/docker.sock:/var/run/docker.sock \
--env=DRONE_RPC_PROTO=https \
--env=DRONE_RPC_HOST=drone.yikyak.org \
--env=DRONE_RPC_SECRET=<Shared_Secret_W_Server> \
    --env=DRONE_RUNNER_CAPACITY=2 \
    --env=DRONE_RUNNER_NAME=my-first-runner \
    --env=DRONE_RPC_DUMP_HTTP=true \
    --env=DRONE_RPC_DUMP_HTTP_BODY=true \
    --env=DRONE_DEBUG=true \
    --publish=3000:3000 \
    --restart=always \
    --name=runner \
    drone/drone-runner-docker:1
    ```
    We can run this command directly or as a shell script, and again run `docker ps`, followed by `docker logs "$container_name"` to confirm we are good to go!

    Note: Docker runner: This runner is a daemon that executes a pipeline in steps inside ephemeral Docker containers. It can be installed as a single docker runner or across multiple machines to create a build cluster. There are other types of runners, just as there are different types of Server integrations/setups. You can see the [additional runner types here](https://docs.drone.io/runner/overview/).

    Now we've arrived at the fun part. It's time to define our pipeline.

    Drone Pipeline
    In this section, we will explore how to configure pipelines generally. Like servers and runners, multiple types of pipelines exist, and we aren't limited to using a single variety inside one repo. We'll mainly focus on the pipeline type used within our example project, a docker variant. After looking at the essential fields of a pipeline, we'll dig deeper by focusing on our example. To define a pipeline, we create a `.drone.yml` in the root directory of our repository.

    Triggers.
    Workspace & cloning.
    Steps.
    Plugins.
    Services
    Conditions.
    Volumes.
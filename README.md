# kiwiland

This program is a small utility I wrote to run on a raspberry pi which is connected to my TV via HDMI and my home ethernet. 

I wanted to be able to control the TV's HDMI CEC capabilities (it is a modern-ish Sony Bravia) as well as run wake on lan commands to my actual media PC (an old Toshiba laptop).

## Installation

You will need to `apt install cec-utils`. You will also need to install `wakeonlan`, but a raspberry pi has this preinstalled.

Then, you can `go get` and `go build` this project. It does not have any other dependencies. 

## Configuration

This program will start a server called kiwiland. You can sign into it with credentials you will provide on first run. You should also provide a non-default cookie salt in the appropriate environment variable (see the startup logs).

You will want to change the `wakeonlan` MAC address, this is located at `/kiwiserver/wolcommand.go`.

You can make this server run automatically using systemd. This is the service file I made:

Filename `/etc/systemd/system/go-kiwiland.service`
```
[Unit]
Description=go-kiwiland
After=network.target

[Service]
User=pi
ExecStart=/home/pi/go/src/github.com/kiwih/kiwiland/kiwiland
WorkingDirectory=/home/pi/go/src/github.com/kiwih/kiwiland
Restart=always

[Install]
WantedBy=multi-user.target
```

Then, you can expose the go server to the internet via nginx, if you'd like. Mine is only exposed to the local intranet.
Filename `/etx/nginx/sites-available/kiwi.land`
```
server {
        listen 80 default_server;
        server_name kiwi.land;
        location / {
                allow 192.168.2.0/24;
                deny all;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $remote_addr;
                proxy_set_header Host $host;
                proxy_pass http://127.0.0.1:3000;
        }
}
```

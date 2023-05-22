# Documentation is a work in progress
Goal of this Project was to enable remote controlling the Microphone of the target computer.

In a setting where the whole room shares one Room-Microphone in a Teams-Setting everyone could remote Mute-Toggle the microphone when raising a Question etc.

## Optinal Flags
- to specify the path /public folder with the index.html `--files ./src/public`
- to change the port listening on: `--port 8080`

## Building and running the Docker image
- `git clone https://github.com/vincepr/go_micmute_server.git`
- `cd go_micmute_server`
- `sudo docker build -t go_auth_proxy .`
- `sudo docker run -it --rm -p 3003:3003 go_auth_proxy`

now we can access it in localhost:3003

If we would want to expose it directly, we could do so with--network="host" . Or we use something like nginx (see below)
- `sudo docker run -it --rm --network="host" -p 3003:3003 go_auth_proxy`
### Auto restart for the docker image
If we want it to autostart after reboots we can add the `-d --restart unless-stopped` flag:
- `sudo docker run -it -d --restart unless-stopped -p 3003:3003 go_auth_proxy`

## Https via nginx & certbot:
Quickest way to add https with minimal effort. Use nginx redirect that upgrades connection with a certificate.
`sudo nano /etc/nginx/sites-available/mic.vprobst.de.conf`
```
server {
	  listen 80;
        listen [::]:80;
        server_name mic.vprobst.de;

        location / {            
                proxy_pass http://localhost:3003/;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection upgrade;
                proxy_set_header Accept-Encoding gzip;  
        }       
}
```

- then test and activate those settings
```
sudo ln -s /etc/nginx/sites-available/mic.vprobst.de.conf /etc/nginx/sites-enabled/mic.vprobst.de.conf
sudo nginx -t
systemctl status nginx sudo systemctl restart nginx
```
- afterwards run certbot

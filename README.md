# The Microphone Remote Volume Control Server
Check out the website at: https://mic.vprobst.de
## About

Goal of this project was to enable remote controlling the microphone of the target computer (running the [csharp-App](https://github.com/vincepr/Mic_Mute)).

In a setting where the whole room shares one Room-Microphone in a Teams-Setting everyone could remote Mute-Toggle the microphone. Enabling a shared Microphone for convenient use, while keeping background noise low.

This server:
- authentificates Users connecting via the [csharp-App](https://github.com/vincepr/Mic_Mute), keeps WebSocket connections to those open and removes inactive ones.
- hosts the website, listens for requests to send Volume-Up/Down/Toggle requests from the website.
- If the requests form website match the username and password for the open connection of a App-User it sends corresponding Volume Control signal over to the App.

### Optinal Flags
- to change the port listening on: `--port 8080`
- to specify the path to the /public folder with the index.html and ressources `--files ./src/public`

### Working on the project
```
git clone https://github.com/vincepr/go_micmute_server.git
cd ./go_micromute_server/src
go run . --port 5555 --files ./public
```
- then you can visit: ``

## Building and running the Docker image
Docker image containts only the Linux(64bit)-binary and the ressources for the website. Arround 7-8MB in size 
- `git clone https://github.com/vincepr/go_micmute_server.git`
- `cd go_micmute_server`
- `sudo docker build -t go_micmute_server .`
- `sudo docker run -it --rm -p 3003:3003 go_micmute_server`

now we can access it in localhost:3003

If we would want to expose it directly, we could do so with--network="host" . Or we use something like nginx (see below)
- `sudo docker run -it --rm --network="host" -p 3003:3003 go_micmute_server`

### Deploy the docker image with autostart
If we want it to autostart after rebooting we can add the `-d --restart unless-stopped` flag:
- `sudo docker run -it -d --restart unless-stopped -p 3003:3003 go_micmute_server`

If we want to shut that container down again we:
- Find ID for the go_micmute_server `sudo docker container ls`
- and shut it down using that Container-ID `sudo docker stop CONTAINERID`

## Https via nginx & certbot:
Quickest way to add https with minimal effort. Use nginx redirect that upgrades connection with a certificate.
- `sudo nano /etc/nginx/sites-available/mic.vprobst.de.conf`
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

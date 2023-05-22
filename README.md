# Documentation is a work in progress

## Https via nginx & certbot:
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

## Building and running with Docker image
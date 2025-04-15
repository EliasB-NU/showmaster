# Showmaster

I am currently rewriting the whole project in a more modular way and
with a new backend & web interface.

Stay tuned.

If you want to use the current version, this would be the docker compose file for it:

```yaml
# Showmaster Docker compose file

services:
  app:
    image: braunelias/showmaster:latest
    restart: always
    networks:
      - showmaster-network
    ports:
      - "80:80"
    depends_on:
      - db
    environment:
      - DBUser=showmaster
      - DBPassword=password
      - Database=showmaster

  db:
    image: postgres:16-alpine
    restart: always
    environment:
      - POSTGRES_USER=showmaster
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    networks:
      - showmaster-network
    volumes:
      - ./data-db:/var/lib/postgresql/data

  pgadmin:
    image: dpage/pgadmin4:latest
    networks:
      - showmaster-network
    environment:
      - PGADMIN_DEFAULT_EMAIL=admin@example.com
      - PGADMIN_DEFAULT_PASSWORD=SomeSecurePassword
    ports:
      - "8080:80"
    volumes:
      - ./data-pq:/var/lib/pgadmin/data

networks:
  showmaster-network:
    external: false
```
To add scenes you need to use a database tool like Datagrip or pgAdmin.
There is also currently no integration with other tools, but I am working on it.

## Description
A tool to manage your scenes across different audio, light and video consoles.
You can use the midi implementation for syncing your consoles or the OSC implementation to control your consoles from a single device.

## Features

- [ ] Midi implementation
- [ ] OSC implementation
- [ ] Web interface
- [ ] Web API
- [ ] Websockets

## Integrations
There will be an integration with my own low latency and lossless 
video streaming service for internal use, with an integrated router 
for video and video playback server. 
You can also create scenes in the llls web interface and use them in your showmaster scenes.
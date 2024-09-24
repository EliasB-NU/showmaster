# ShowMaster - V2

This is a cross-platform/browser synced list of scenes for your local event.  
  
You can highlight a specific scene by clicking on it or moving the currently highlighted scene with your up/down arrow keys. 
The currently highlighted row is synced with every device.  


To add data to showmaster, you need to use a database tool like Jetbrains DataGrip or pqAdmin, maybe I will create an admin panel in the future. 

## Showcase
I am too stupid to embed a video in the readme, you can download it in the media folder, sorry for any troubles

## Setup
Just copy the contents of the "docker compose" file to your local machine, I assume you have docker & docker compose installed, and run it. On the mainpage you can create and select projects.  
To easily add data, I included a pgadmin4 service, you can include it, to get a fast and easy option to view and edit the database in the web.

Example Docker Compose File:

```yaml
# Showmaster Docker compose file

services:
  app:
    image: braunelias/showmaster:latest
    restart: always
    ports:
      - "80:80"
    depends_on:
      - db
    environment:
      - DBUser=showmaster
      - DBPassword=password
      - Database=showmaster
      - TimeZone=Europe/Berlin
      - AdminUserName=admin
      - AdminPassword=1234

  db:
    image: postgres:16-alpine
    restart: always
    environment:
      - POSTGRES_USER=showmaster
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    volumes:
      - psql-data:/var/lib/postgresql/data

  pgadmin:
    image: dpage/pgadmin4:latest
    environment:
      - PGADMIN_DEFAULT_EMAIL=admin@example.com
      - PGADMIN_DEFAULT_PASSWORD=SomeSecurePassword
    ports:
      - "8080:80"

volumes:
  psql-data:
```

This docker compose file includes a pqadmin instance for easy management. If you want to use another piece of software, just remove it.

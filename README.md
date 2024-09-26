# ShowMaster - V3 (W.I.P. | latest stable version V2.3.1 (read readme in the v2.3.1 branch))

This is a cross-platform/browser synced list of scenes for your local event.  
  
You can highlight a specific scene by clicking on it or moving the currently highlighted scene with your up/down arrow keys. 
The currently highlighted row is synced with every device.  

It has a fully interactive web ui to move the different entries around and edit them,
watch the showcase video for more information

## Showcase
Coming soon

## Setup
Just copy the contents of the "docker compose" file to your local machine, 
I assume you have docker & docker compose installed, and run it.

**Example Docker Compose File:**

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
      - DBUser=showmaster      # Should correspond to the environment variables ->
      - DBPassword=password    # set in the db section of this docker compose file
      - Database=showmaster
      - TimeZone=Europe/Berlin # Change to your local timezone
      - AdminUserName=admin    # Email is admin@example.com
      - AdminPassword=1234     # Please change to a secure password

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

volumes:
  psql-data:
```

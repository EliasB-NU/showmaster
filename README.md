# ShowMaster - V1

This is an cross platform/browser synced list of scenes for your local event.  
  
You can highlight a specific scene by clicking on it or moving the currently highlighted scene with your up/down arrow keys. 
The currently highlighted row is synced with every device.  


To add data to showmaster, you need to use an database tool like Jetbrains DataGrip or pqAdmin, maybe I will create an admin panel in the future. 

## Setup
Just copy the contents of the "docker compose" file to your local machine, I assume you have docker & docker compose installed, and run it. By changing the project name and restarting, the showmaster programm creates a new table, so you can easily switch between projects.  
To easily add data, I included a pgadmin4 service, you can include it, to get an fast and easy option to view and edit the database in the web. 
# ShowMaster - V1

This is an cross platform/browser synced list of scenes for your local event.  
  
You can highlight a specific scene by clicking on it or moving the currently highlighted scene with your up/down arrow keys. 
The currently highlighted row is synced with every device.  


To add data to showmaster, you need to use an database tool like Jetbrains DataGrip or pqAdmin, maybe I will create an admin panel in the future. 

## Setup
To use this tool, you need a postgresql database and you just need to add all the details to the config file in the config folder (don't forget to change the name from "config.example.json" to "config.json").

The correspending table is created with the name given under the "Project" entry, if you allready created a table with the same name as the "Project" name, it will not be overwritten and the entries you already created will be used.

You have to replace the YOURIP statement in the "public/script.js" file with the ip the frontend will run on.

In the server folder is an .service file, which automaticly starts the programm after you compiled the go project. Please check folder location. 

If you have question, contact me via email (braunelias@tghd.email) or Discord (tg_eliasb)
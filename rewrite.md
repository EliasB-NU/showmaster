# Database

---
> 'id' is data type 'SERIAL PRIMARY KEY', whereas position is of type 'doubleprecision' to allow for supsteps
> and reordering without the need for reformating the whole table every time
## User

---

### Table

---

| id |    name    |         email          |      password      | permissionLevel |
|:--:|:----------:|:----------------------:|:------------------:|:---------------:|
| 1  |  testUser  |  norepley@example.com  |  exmaplePassword   |        4        |
| 2  | secondUser | secondUser@example.com | verySecurePassword |        1        |


### Permission Level

---

| id |  name   |                       Permissions                       |
|:--:|:-------:|:-------------------------------------------------------:|
| 1  | watcher |                can be added to projects                 |
| 2  |  user   |                 1 & create own projects                 |
| 3  |  admin  | 1+2 & you can see and edit all projects/users -> admin  |

## Projects

---
### Table

---
| id |    name     |      table       | creator  | users+permission level |
|:--:|:-----------:|:----------------:|:--------:|:----------------------:|
| 1  | testProject | testProjectTable | testUser |      secondUser+2      |
> Links to the project table

| id | position |     name     |      audio       |       light        | pptx |             note             |
|:--:|:--------:|:------------:|:----------------:|:------------------:|:----:|:----------------------------:|
| 1  |    1     | First Scene  | Commentary Mikes | Spot on Commentary |  1   | Wait til normal light is off |
| 2  |   1.1    | Second Scene |  Computer Audio  | Only Video Screen  |  2   |     first light -> pptx      |
| 3  |    2     | Thrid Scene  | Orchestra Audio  |     Orchestra      |  3   |    pptx -> light + audio     |

### Permission Level

---
| id |     name     |      Permission       |
|:--:|:------------:|:---------------------:|
| 1  |   watcher    |    can see project    |
| 2  | collaborator | 1 & can edit project  |
| 3  |    admin     | 1+2 & can edit users  |

# Login

---
Password Forgot button and register button
you get an email with a code to change your password

# Main Side

---
- top left: admin logo, to access admin settings, if you have the permission level
- top right: user logo, to logout (maybe more coming in the future)
- bottom right: plus logo for new project,if you have the permission level
- bottom middle: Impress & CC

tiles with each project you can see/have access to/created

# Project Side

---
- top left: admin logo, to access admin settings, if you have the permission level
- top right: user logo, to logout (maybe more coming in the future)
- bottom right: plus logo for new row, if you have the permission level
- bottom middle: Impress & CC
- top left, before admin (if you have the admin permission), a back arrow

# Project Update

---
```go
# Project Creation Update
type NewProjectStruct struct {
    Name String `json:"name"`
    Creator String `json:"creator"`
}

# Project Update
type ProjectUpdateStruct struct {
    Name String `json:"name"` # If empty, no action
    Users []String `json:"users"` # If empty, no action | Formatting 'email+permLVL'
    
}
```


# Row Update

---
```go
# New Row/Update Row | Null values will not be inserted on an update request
type RowStruct struct {
    Position float32 `json:"position"`
	Name String `json:"name"`
	Audio String `json:"audio"`
	Light String `json:"light"`
	PPTX String `json:"pptx"`
	Notes String `json:"pptx"`
}
```

# General Update Messages

---
```go
# Time and Highlighted Row | Null values will not be processed
type Update struct { # Incoming
	Project String `json:"project"`
	HighlightedRow float32 `json:"highlightedrow"`
	TimerStatus String `json:"timerstatus"` # Possible Values: running, stopped, reset
}

# The current time and status for new clients
type TimerStatus struct { # Outgoing
	Running bool `json:"status"`
    Time uint64 `json:"time"`
}
```

# Websocket

---
The websocket is used to send the new highlighted row and a refresh request to all clients
in case of a change. It also sends the status of the timer on an update

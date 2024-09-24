# Database

---
> 'id' is data type 'SERIAL PRIMARY KEY', whereas position is of type 'double precision' to allow for sup steps
> and reordering without the need for reformating the whole table every time
## User

---

### Table

---

| id |    name    |         email          |      password      | permissionLevel |
|:--:|:----------:|:----------------------:|:------------------:|:---------------:|
| 1  |  testUser  |  norepley@example.com  |  examplePassword   |        4        |


### Permission Level

---

| id |  name   |                      Permissions                       |
|:--:|:-------:|:------------------------------------------------------:|
| 1  | watcher |                    can see projects                    |
| 2  |  user   |                1 & create own projects                 |
| 3  |  admin  | 1+2 & you can see and edit all projects/users -> admin |

## Projects

---
### Table

---
| id |    name     |      table       | creator  | 
|:--:|:-----------:|:----------------:|:--------:|
| 1  | testProject | testProjectTable | testUser |
> Name of the project table

| id | position |     name     |      audio       |       light        | pptx |             note             |
|:--:|:--------:|:------------:|:----------------:|:------------------:|:----:|:----------------------------:|
| 1  |    1     | First Scene  | Commentary Mikes | Spot on Commentary |  1   | Wait til normal light is off |
| 2  |   1.1    | Second Scene |  Computer Audio  | Only Video Screen  |  2   |     first light -> pptx      |
| 3  |    2     | Third Scene  | Orchestra Audio  |     Orchestra      |  3   |    pptx -> light + audio     |

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
package database

// NewProjectStruct Project Creation Update
type NewProjectStruct struct {
	Name    string `json:"name"`
	Creator string `json:"creator"`
}

// ProjectUpdateStruct Project Update
type ProjectUpdateStruct struct {
	Name string `json:"name"` // If empty, no action
	Users []string `json:"users"` // If empty, no action | Formatting 'email+permLVL'
}
```


# Row Update

---

```go
package database

// RowStruct New Row/Update Row | Null values will not be inserted on an update request
type RowStruct struct {
	Position float32 `json:"position"`
	Name     string  `json:"name"`
	Audio    string  `json:"audio"`
	Light    string  `json:"light"`
	PPTX     string  `json:"pptx"`
	Notes    string  `json:"notes"`
}
```

# General Update Messages

---

```go
package web

// Update Time and Highlighted Row | Null values will not be processed
type Update struct { // Incoming
	Project        string  `json:"project"`
	HighlightedRow float32 `json:"highlightedrow"`
	TimerStatus    string  `json:"timerstatus"` // Possible Values: running, stopped, reset
}

// TimerStatus The current time and status for new clients
type TimerStatus struct { // Outgoing
	Running bool   `json:"status"`
	Time    uint64 `json:"time"`
}
```

# Websocket

---
The websocket is used to send the new highlighted row and a refresh request to all clients
in case of a change. It also sends the status of the timer on an update

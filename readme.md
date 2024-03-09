# Project: forum
## Overview
This project is a web forum designed to facilitate communication between users through posts and comments, allowing users to associate categories with posts, like or dislike posts and comments, and filter posts based on different criteria. The project is written in Go and utilizes SQLite as the database management system for storing user data, posts, comments, and categories.

## Installation
Clone the repository from Gitea: git clone https://learn.zone01dakar.sn/git/mouhamsylla/forum
Navigate to the project directory: cd project
Install any necessary dependencies.

## Usage
```
$ make app

$ ./app
```
## Features

* Authentication

Users can register with a unique email.
Passwords are encrypted for security .
Unique sessions are created using cookies with expiration dates.
* Communication

Registered users can create posts and associate one or more categories with them.
Registered users can create comments on posts.
Posts and comments are visible to all users, but only registered users can create them.
Likes and Dislikes
Registered users can like or dislike posts and comments.
Like and dislike counts are visible to all users.
* Filtering

Users can filter displayed posts by categories, created posts, and liked posts.
Filtering by categories functions as subforums, allowing users to view posts related to specific topics.
* Database Structure

The database consists of the following tables:

users: Stores user information including email, username, and encrypted password.
posts: Contains post content, associated categories, and user IDs of the creators.
comments: Stores comments on posts along with the user IDs of the commenters.
likes_dislikes: Tracks likes and dislikes on both posts and comments.

## Structure

```
.
├── Api
│   ├── controllers
│   │   ├── authController.go
│   │   ├── commentController.go
│   │   ├── connectController.go
│   │   ├── controllerError.go
│   │   ├── Controller.go
│   │   ├── homeController.go
│   │   └── postController.go
│   ├── models
│   │   ├── instance.go
│   │   └── models.go
│   ├── services
│   │   ├── authentication.go
│   │   └── sessions.go
│   └── storage
│       └── storage.go
├── app
├── App
│   ├── internal
│   │   ├── assets
│   │   │   ├── connection.html
│   │   │   ├── error.html
│   │   │   ├── home.css
│   │   │   ├── home.html
│   │   │   ├── images
│   │   │   │   ├── 4.jpg
│   │   │   │   ├── b.jpg
│   │   │   │   ├── sg.svg
│   │   │   │   └── sv2.svg
│   │   │   ├── index.html
│   │   │   ├── log.css
│   │   │   ├── main.wasm
│   │   │   ├── script.js
│   │   │   └── wasm_exec.js
│   │   ├── boot
│   │   │   └── init.go
│   │   ├── components
│   │   │   ├── comment.go
│   │   │   ├── post.go
│   │   │   └── profile.go
│   │   ├── directives
│   │   │   ├── commentDirective.go
│   │   │   ├── directive.go
│   │   │   ├── filterDirective.go
│   │   │   ├── inputDirective.go
│   │   │   ├── loginDirective.go
│   │   │   ├── postDirective.go
│   │   │   └── registerDirective.go
│   │   └── utils
│   │       ├── DOM.go
│   │       └── utils.go
│   └── main.go
├── cmd
│   └── main.go
├── db
│   ├── forumDB.db
│   └── migrates
│       ├── 2024-03-09-09-52-53-create-table-Categories.sql
│       ├── 2024-03-09-09-52-53-create-table-Comments.sql
│       ├── 2024-03-09-09-52-53-create-table-Post.sql
│       └── ..............................
├── dockerbash.sh
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── orm
│   ├── delete.go
│   ├── insert.go
│   ├── models.go
│   ├── orm.go
│   ├── scan.go
│   ├── sql.go
│   ├── tags.go
│   ├── types.go
│   └── update.go
├── readme.md
├── server
│   ├── middleware
│   │   └── middleware.go
│   ├── router
│   │   ├── models.go
│   │   ├── router.go
│   │   └── trie.go
│   └── server.go
├── start.sh
└── utils
    └── utils.go

```


## Authors
-  Aliou Niang
-   Ndieye Diop
-   Abdoul Aziz Ba
-   Mouhamadou Sylla

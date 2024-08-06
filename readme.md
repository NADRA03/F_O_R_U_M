# pages
1. login
2. sign up 
3. home (posts with comments and likes and categories filter)
4. created posts 
5. create post
6. liked post 
7. profile as side bar at home + button for liked posts + button for created posts
8. moderator 

# cookies and sessions 
the server can use the session ID stored in the cookie to retrieve the corresponding session data.

# Tables 
1. id username password email image    user
2. post-id id text media date category        post
3. comment-id id post-id comment date       comment
4. like-id id post-id                    like 

# tutorial 
1. download sqlite add it to path 
2. download it into go "go get github.com/mattn/go-sqlite3"
3. install gcc 
4. run database "sqlite3 forum.db" will open a terminal 
5. fake users: moderator 123, not-moderator 123
moderator can access moderator.html
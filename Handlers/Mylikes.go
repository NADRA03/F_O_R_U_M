package Forum

import (
    "database/sql"
    "html/template"
    "net/http"
    "sort"
)

func MyLikesHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        
        //err
				if r.URL.Path != "/mylikes" {
					RenderErrorPage(w, http.StatusNotFound)
					return
		}



        session, _ := store.Get(r, "mysession")

        if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        userID, _ := session.Values["id"].(int)
        username, _ := session.Values["username"].(string)
        image, ok := session.Values["profileImage"].(string)
    
        if !ok || image == "" {
            image = defaultProfileImage 
        }
        rows, err := db.Query(`
            SELECT p.post_id, p.text, p.media, p.date, p.category 
            FROM post p
            JOIN like l ON p.post_id = l.post_id
            WHERE l.user_id = ?
        `, userID)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError) 
            return
        }
        defer rows.Close()

        var posts []struct {
            PostID   int
            Text     string
            Media    string
            Date     string
            Category string
        }

        for rows.Next() {
            var post struct {
                PostID   int
                Text     string
                Media    string
                Date     string
                Category string
            }
            err := rows.Scan(&post.PostID, &post.Text, &post.Media, &post.Date, &post.Category)
            if err != nil {
                RenderErrorPage(w, http.StatusInternalServerError) 
                return
            }
            posts = append(posts, post)
        }
        


        //select liked comments
        commentRows, err := db.Query(`
            SELECT c.comment_id, c.comment, c.date, p.post_id, u.username 
            FROM comment c
            JOIN comment_like cl ON c.comment_id = cl.comment_id
            JOIN user u ON c.user_id = u.id
            JOIN post p ON c.post_id = p.post_id
            WHERE cl.user_id = ?
        `, userID)
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }
        defer commentRows.Close()

        var comments []struct {
            CommentID int
            Comment   string
            Date      string
            PostID    int
            Username  string
        }

        for commentRows.Next() {
            var comment struct {
                CommentID int
                Comment   string
                Date      string
                PostID    int
                Username  string
            }
            err := commentRows.Scan(&comment.CommentID, &comment.Comment, &comment.Date, &comment.PostID, &comment.Username)
            if err != nil {
                RenderErrorPage(w, http.StatusInternalServerError)
                return
            }
            comments = append(comments, comment)
        }

        // Create a unified structure for posts and comments
        type FeedItem struct {
            Type     string // "post" or "comment"
            ID       int
            Content  string
            Date     string
            Category string // For posts
            Username string // For comments
            PostID   int
        }

        var feed []FeedItem

        // Add posts to the feed
        for _, post := range posts {
            feed = append(feed, FeedItem{
                Type:     "post",
                ID:       post.PostID,
                Content:  post.Text,
                Date:     post.Date,
                Category: post.Category,
            })
        }

        // Add comments to the feed
        for _, comment := range comments {
            feed = append(feed, FeedItem{
                Type:     "comment",
                ID:       comment.CommentID,
                Content:  comment.Comment,
                Date:     comment.Date,
                Username: comment.Username,
                PostID:   comment.PostID,
            })
        }

        // Sort the feed by date in descending order
        sort.Slice(feed, func(i, j int) bool {
            return feed[i].Date > feed[j].Date // Descending order
        })

        // Render the template with the combined feed
        tmpl, err := template.ParseFiles("HTML/mylikes.html")
        if err != nil {
            RenderErrorPage(w, http.StatusInternalServerError)
            return
        }

        tmpl.Execute(w, struct {
            Feed     []FeedItem
            Username string
            Image    string
        }{
            Feed:     feed,
            Username: username,
            Image:    image,
        })
    }
}
package Forum

import (
    "database/sql"
	"fmt"
)

type Table struct {
    Name    string
    Columns []string
    Rows    [][]interface{}
}

// func createTables(db *sql.DB) error {
//     statements := []string{
//         `CREATE TABLE IF NOT EXISTS user (
//             id INTEGER PRIMARY KEY AUTOINCREMENT,
//             username TEXT NOT NULL,
//             password TEXT NOT NULL,
//             email TEXT NOT NULL,
//             image TEXT
//         );`,
//         `CREATE TABLE IF NOT EXISTS post (
//             post_id INTEGER PRIMARY KEY AUTOINCREMENT,
//             user_id INTEGER NOT NULL,
//             text TEXT,
//             media TEXT,
//             date TEXT,
//             category TEXT,
//             FOREIGN KEY (user_id) REFERENCES user(id)
//         );`,
//         `CREATE TABLE IF NOT EXISTS comment (
//             comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
//             user_id INTEGER NOT NULL,
//             post_id INTEGER NOT NULL,
//             comment TEXT,
//             date TEXT,
//             FOREIGN KEY (user_id) REFERENCES user(id),
//             FOREIGN KEY (post_id) REFERENCES post(post_id)
//         );`,
//         `CREATE TABLE IF NOT EXISTS like (
//             like_id INTEGER PRIMARY KEY AUTOINCREMENT,
//             user_id INTEGER NOT NULL,
//             post_id INTEGER NOT NULL,
//             FOREIGN KEY (user_id) REFERENCES user(id),
//             FOREIGN KEY (post_id) REFERENCES post(post_id)
//         );`,
//     }

//     for _, statement := range statements {
//         _, err := db.Exec(statement)
//         if err != nil {
//             return err
//         }
//     }

//     return nil
// }

func fetchTables(db *sql.DB) ([]Table, error) {
    var tables []Table

    rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var tableName string
        rows.Scan(&tableName)
        
        columns, rowsData, err := fetchTableData(db, tableName)
        if err != nil {
            return nil, err
        }

        tables = append(tables, Table{
            Name:    tableName,
            Columns: columns,
            Rows:    rowsData,
        })
    }

    return tables, nil
}

func fetchTableData(db *sql.DB, tableName string) ([]string, [][]interface{}, error) {
    columnsQuery := fmt.Sprintf("PRAGMA table_info(%s);", tableName)
    columnsRows, err := db.Query(columnsQuery)
    if err != nil {
        return nil, nil, err
    }
    defer columnsRows.Close()

    var columns []string
    for columnsRows.Next() {
        var cid int
        var name, ctype string
        var notnull, pk int
        var dfltValue interface{}
        columnsRows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk)
        columns = append(columns, name)
    }

    dataQuery := fmt.Sprintf("SELECT * FROM %s;", tableName)
    dataRows, err := db.Query(dataQuery)
    if err != nil {
        return nil, nil, err
    }
    defer dataRows.Close()

    var rowsData [][]interface{}
    for dataRows.Next() {
        columnsCount := len(columns)
        values := make([]interface{}, columnsCount)
        valuePtrs := make([]interface{}, columnsCount)
        for i := range values {
            valuePtrs[i] = &values[i]
        }

        dataRows.Scan(valuePtrs...)
        rowsData = append(rowsData, values)
    }

    return columns, rowsData, nil
}

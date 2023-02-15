
package config
 
import (
    "database/sql"
)
 
// @Desc: Open connection to mySQL database
func Connect() *sql.DB {
  
  
	// *Change accordingly:
    db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/sys")

    // if there is an error opening the connection, handle it
    if err != nil {
        panic(err.Error())
    }


	return db
}
 
 1 - Edit /etc/postgres/pg_hba.conf



    # TYPE  DATABASE    USER    ADDRESS                 METHOD
    local   all         postgres                        peer map=pg_root
    local   all         all                             peer
    #LAN access
    host    all         all     192.168.1.0/24          trust
    #localhost access
    host    all         all     127.0.0.1/32            trust

 2 - Reload postgres config:

    # synoservice --reload pgsql

3 - Connection from go:



    /*==========================
	    pgSynoConnection.go
    ===========================*/
    
    package main
    
    import (
    	"database/sql"
    	"log"
    	_ "github.com/lib/pq"
    )
    
    var db *sql.DB
    
    // connect to the Db
    func init() {
    	var err error
    	db, err = sql.Open("postgres", "postgres://username:passwd@localhost/database?sslmode=disable")
    	if err != nil {
    		log.Fatal(err)
    	}
    }
    
    func main() {
        if err := db.Ping(); err != nil {
        log.Fatal(err)
      }
    
    	db.Close()
    }

 


4 - Cross-compile for ARM:

    GOOS=linux GOARCH=arm GOARM=7 go build pgSynoConnection.go


5 - scp binary to Disk Station

    scp -P sshport pgSynoConnection user@192.168.1.100:/volume1/path/
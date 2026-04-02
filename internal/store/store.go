package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Package struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Version string `json:"version"`
	Registry string `json:"registry"`
	Ecosystem string `json:"ecosystem"`
	Size string `json:"size"`
	Checksum string `json:"checksum"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"parcel.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS packages(id TEXT PRIMARY KEY,name TEXT NOT NULL,version TEXT DEFAULT '',registry TEXT DEFAULT '',ecosystem TEXT DEFAULT '',size TEXT DEFAULT '',checksum TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Package)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO packages(id,name,version,registry,ecosystem,size,checksum,created_at)VALUES(?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Version,e.Registry,e.Ecosystem,e.Size,e.Checksum,e.CreatedAt);return err}
func(d *DB)Get(id string)*Package{var e Package;if d.db.QueryRow(`SELECT id,name,version,registry,ecosystem,size,checksum,created_at FROM packages WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Version,&e.Registry,&e.Ecosystem,&e.Size,&e.Checksum,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Package{rows,_:=d.db.Query(`SELECT id,name,version,registry,ecosystem,size,checksum,created_at FROM packages ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Package;for rows.Next(){var e Package;rows.Scan(&e.ID,&e.Name,&e.Version,&e.Registry,&e.Ecosystem,&e.Size,&e.Checksum,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM packages WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM packages`).Scan(&n);return n}

package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Package struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Version string `json:"version"`
	Registry string `json:"registry"`
	Checksum string `json:"checksum"`
	SizeBytes int `json:"size_bytes"`
	Downloads int `json:"downloads"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"parcel.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS packages(id TEXT PRIMARY KEY,name TEXT NOT NULL,version TEXT DEFAULT '',registry TEXT DEFAULT '',checksum TEXT DEFAULT '',size_bytes INTEGER DEFAULT 0,downloads INTEGER DEFAULT 0,status TEXT DEFAULT 'published',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Package)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO packages(id,name,version,registry,checksum,size_bytes,downloads,status,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Version,e.Registry,e.Checksum,e.SizeBytes,e.Downloads,e.Status,e.CreatedAt);return err}
func(d *DB)Get(id string)*Package{var e Package;if d.db.QueryRow(`SELECT id,name,version,registry,checksum,size_bytes,downloads,status,created_at FROM packages WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Version,&e.Registry,&e.Checksum,&e.SizeBytes,&e.Downloads,&e.Status,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Package{rows,_:=d.db.Query(`SELECT id,name,version,registry,checksum,size_bytes,downloads,status,created_at FROM packages ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Package;for rows.Next(){var e Package;rows.Scan(&e.ID,&e.Name,&e.Version,&e.Registry,&e.Checksum,&e.SizeBytes,&e.Downloads,&e.Status,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Package)error{_,err:=d.db.Exec(`UPDATE packages SET name=?,version=?,registry=?,checksum=?,size_bytes=?,downloads=?,status=? WHERE id=?`,e.Name,e.Version,e.Registry,e.Checksum,e.SizeBytes,e.Downloads,e.Status,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM packages WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM packages`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Package{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,version,registry,checksum,size_bytes,downloads,status,created_at FROM packages WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Package;for rows.Next(){var e Package;rows.Scan(&e.ID,&e.Name,&e.Version,&e.Registry,&e.Checksum,&e.SizeBytes,&e.Downloads,&e.Status,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM packages GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}

package server
import("fmt";"io";"mime";"net/http";"os";"path/filepath";"strconv";"strings";"time";"crypto/rand";"encoding/hex";"github.com/stockyard-dev/stockyard-parcel/internal/store")
func randToken()string{b:=make([]byte,16);rand.Read(b);return hex.EncodeToString(b)}
func fmtSize(n int64)string{if n<1024{return fmt.Sprintf("%d B",n)};if n<1024*1024{return fmt.Sprintf("%.1f KB",float64(n)/1024)};return fmt.Sprintf("%.1f MB",float64(n)/1024/1024)}
func(s *Server)handleUpload(w http.ResponseWriter,r *http.Request){
    if !s.limits.IsPro(){n,_:=s.db.CountUploads();if n>=10{writeError(w,403,"free tier: 10 files max");return}}
    r.ParseMultipartForm(50<<20)
    file,header,err:=r.FormFile("file");if err!=nil{writeError(w,400,"file required");return};defer file.Close()
    token:=randToken();ext:=filepath.Ext(header.Filename)
    filename:=token+ext
    dest:=filepath.Join(s.db.DataDir,"files",filename)
    out,err:=os.Create(dest);if err!=nil{writeError(w,500,"storage error");return};defer out.Close()
    size,err:=io.Copy(out,file);if err!=nil{writeError(w,500,"write error");return}
    mtype:=header.Header.Get("Content-Type");if mtype==""{mtype=mime.TypeByExtension(ext)};if mtype==""{mtype="application/octet-stream"}
    maxDL,_:=strconv.Atoi(r.FormValue("max_downloads"))
    var expiresAt *time.Time
    if exp:=r.FormValue("expires_hours");exp!=""{if h,err:=strconv.Atoi(exp);err==nil&&h>0{t:=time.Now().Add(time.Duration(h)*time.Hour);expiresAt=&t}}
    u:=&store.Upload{Filename:filename,OriginalName:header.Filename,Size:size,MimeType:mtype,ShareToken:token,MaxDownloads:maxDL,ExpiresAt:expiresAt}
    if err:=s.db.CreateUpload(u);err!=nil{writeError(w,500,err.Error());return}
    writeJSON(w,201,map[string]interface{}{"id":u.ID,"token":token,"download_url":"/d/"+token,"filename":header.Filename,"size":size})}
func(s *Server)handleListFiles(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListUploads();if list==nil{list=[]store.Upload{}};writeJSON(w,200,list)}
func(s *Server)handleDeleteFile(w http.ResponseWriter,r *http.Request){
    id,_:=strconv.ParseInt(r.PathValue("id"),10,64)
    list,_:=s.db.ListUploads()
    for _,u:=range list{if u.ID==id{os.Remove(filepath.Join(s.db.DataDir,"files",u.Filename));s.db.DeleteUpload(id,u.Filename);break}}
    writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleDownload(w http.ResponseWriter,r *http.Request){
    token:=r.PathValue("token");u,_:=s.db.GetUploadByToken(token);if u==nil{http.NotFound(w,r);return}
    if u.ExpiresAt!=nil&&time.Now().After(*u.ExpiresAt){http.Error(w,"Link expired",410);return}
    if u.MaxDownloads>0&&u.Downloads>=u.MaxDownloads{http.Error(w,"Download limit reached",410);return}
    s.db.IncrDownloads(u.ID)
    path:=filepath.Join(s.db.DataDir,"files",u.Filename)
    w.Header().Set("Content-Disposition","attachment; filename="+u.OriginalName)
    w.Header().Set("Content-Type",u.MimeType)
    http.ServeFile(w,r,path)}
func(s *Server)handleStats(w http.ResponseWriter,r *http.Request){n,_:=s.db.CountUploads();sz,_:=s.db.TotalSize();writeJSON(w,200,map[string]interface{}{"files":n,"total_bytes":sz})}
var _=strings.TrimSpace

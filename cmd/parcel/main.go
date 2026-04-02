package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-parcel/internal/server";"github.com/stockyard-dev/stockyard-parcel/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="8670"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./parcel-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("parcel: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Parcel — Self-hosted package registry proxy\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("parcel: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}

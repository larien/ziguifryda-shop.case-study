package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	//	"reflect"
)

var (
	addr  = flag.String("addr", ":9000", "http service address") // Q=17, R=18
	store = sessions.NewCookieStore([]byte("something-very-secret"))
)

func main() {
	store.Options.Secure = true
	flag.Parse()
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", http.HandlerFunc(index))
	// http://techblog.bozho.net/?p=401
	r.HandleFunc("/login", http.HandlerFunc(login))
	r.HandleFunc("/register", http.HandlerFunc(register))
	r.HandleFunc("/registered", http.HandlerFunc(registered))
	r.HandleFunc("/logout", http.HandlerFunc(logout))
	r.HandleFunc("/products", http.HandlerFunc(products))
	r.HandleFunc("/admin", http.HandlerFunc(admin_index))
	r.HandleFunc("/admin/products", http.HandlerFunc(admin_products))
	r.HandleFunc("/admin/users", http.HandlerFunc(admin_users))
	r.HandleFunc("/admin/orders", http.HandlerFunc(admin_orders))
	r.HandleFunc("/admin/login", http.HandlerFunc(admin_login))
	r.HandleFunc("/admin/logout", http.HandlerFunc(admin_logout))
	http.Handle("/", r)
	fs := JustFilesFilesystem{http.Dir("static/")}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))

	go func() {
		log.Fatalf("ListenAndServe: %v", http.ListenAndServe(":8080", http.HandlerFunc(notlsHandler)))
	}()
	go func() {
		fmt.Println("https://localhost" + *addr)
		// go run generate_cert.go --host localhost
		// this works for other domain than localhost, no need to supply domain name:
		// openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
		// http://stackoverflow.com/questions/10175812/how-to-build-a-self-signed-certificate-with-openssl
		log.Fatalf("ListenAndServeTLS: %v", http.ListenAndServeTLS(*addr, "tls/cert.pem", "tls/key.pem", nil))
	}()
	select {}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method != "GET" || r.URL.Path != "/" {
		serve404(w)
		return
	}

	session, _ := store.Get(r, "session-name")

	pageTemplate, err := template.ParseFiles("tpl/index.html", "tpl/header.html", "tpl/bar.html", "tpl/products.html", "tpl/footer.html")
	if err != nil {
		serveError(w, err)
	}

	tplValues := map[string]interface{}{"Header": "Arapuã", "HeaderMessage": "- O melhor site de compras online!"}
	db, err := sql.Open("sqlite3", "file:./db/app.db?foreign_keys=true")
	if err != nil {
		serveError(w, err)
		return
	}
	defer db.Close()

	sql := "select title, description, price, quantity, filename from products order by title"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		serveError(w, err)
		return
	}
	defer rows.Close()

	levels := []map[string]string{}
	var title, description, price, quantity, filename string
	for rows.Next() {
		rows.Scan(&title, &description, &price, &quantity, &filename)
		levels = append(levels, map[string]string{"title": title, "description": description, "price": price, "quantity": quantity, "filename": filename})
	}
	tplValues["levels"] = levels

	rows.Close()

	if _, ok := session.Values["login"]; ok {
		tplValues["login"] = session.Values["login"]
	}

	pageTemplate.Execute(w, tplValues)
	if err != nil {
		serveError(w, err)
	}
}

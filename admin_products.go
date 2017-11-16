package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func admin_products(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	tplValues := map[string]interface{}{"Header": "Products"}

	authorized := false
	if i, ok := session.Values["admin_login"]; ok {
		if i == "admin" {
			authorized = true
		}
		tplValues["admin_login"] = i
	}

	if !authorized {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	fVals := map[string]string{}
	if r.Method == "POST" {
		if false {
			hah, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("%s", err)
			}
			fmt.Printf("%v", string(hah))
			return
		}

		reader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			if part.FileName() == "" { // normal text field
				b, err := ioutil.ReadAll(part)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fVals[part.FormName()] = html.EscapeString(string(b))
			} else {
				assetsDir := "static/assets/"
				maxFileName := 20
				if err := os.MkdirAll(assetsDir, 0777); err != nil {
					log.Println(err)
					return
				}

				b, err := ioutil.ReadAll(part)
				if err != nil {
					log.Fatal(err)
				}
				ext := Extension(b)

				var path string
				var fileName string
				for {
					fileName = randSeq(maxFileName) + ext
					path = assetsDir + fileName
					if _, err := os.Stat(path); os.IsNotExist(err) {
						//got unique path, so break
						break
					}
				}

				dst, err := os.Create(path)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer dst.Close()

				if _, err := io.Copy(dst, bytes.NewReader(b)); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fVals[part.FormName()] = fileName
			}
		}
	}

	if fVals["sent"] == "yes" {

		message_error := []string{}
		if fVals["title"] == "" {
			message_error = append(message_error, "O título não pode ser vazio!")
		}
		if fVals["description"] == "" {
			message_error = append(message_error, "A descrição não pode ser vazia!")
		}
		if fVals["price"] == "" {
			message_error = append(message_error, "O preço não pode ser vazio!")
		}
		if fVals["quantity"] == "" {
			message_error = append(message_error, "A quantidade não pode ser vazia!")
		}
		if fVals["filename"] == "" {
			message_error = append(message_error, "O arquivo não pode ser vazio!")
		}

		if len(message_error) != 0 {
			tplValues["message_error"] = message_error
			tplValues["fVals"] = fVals
		} else {

			db, err := sql.Open("sqlite3", "file:./db/app.db?foreign_keys=true")
			if err != nil {
				fmt.Println(err)
				serveError(w, err)
				return
			}
			defer db.Close()

			tx, err := db.Begin()
			if err != nil {
				fmt.Println(err)
				serveError(w, err)
				return
			}

			stmt, err := tx.Prepare("insert into products(title, description, price, quantity, filename) values(?, ?, ?, ?, ?)")
			if err != nil {
				fmt.Println(err)
				serveError(w, err)
				return
			}

			defer stmt.Close()

			res, err := stmt.Exec(fVals["title"], fVals["description"], fVals["price"], fVals["quantity"], fVals["filename"])
			if err != nil {
				fmt.Println(err)
				serveError(w, err)
				return
			}
			last, err := res.LastInsertId()
			if err != nil {
				fmt.Println(err)
				serveError(w, err)
				return
			}
			fmt.Println("last", last)
			tx.Commit()

			session.Values["last_product"] = last
			session.Save(r, w)
			tplValues["message_info"] = []string{"Produto cadastrado com sucesso!"}
		}
	}

	levels, err := getProducts()
	tplValues["levels"] = levels
	if err != nil {
		serveError(w, err)
	}

	pageTemplate, err := template.ParseFiles("tpl/admin_products.html", "tpl/header.html", "tpl/admin_bar.html", "tpl/footer.html")
	if err != nil {
		serveError(w, err)
	}

	err = pageTemplate.Execute(w, tplValues)
	if err != nil {
		serveError(w, err)
	}
}

func getProducts() ([]map[string]string, error) {
	db, err := sql.Open("sqlite3", "file:./db/app.db?foreign_keys=true")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()

	sql := "select title, description, price, quantity, filename from products order by title"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		return nil, err
	}
	defer rows.Close()

	levels := []map[string]string{}
	var title, description, price, quantity, filename string
	for rows.Next() {
		rows.Scan(&title, &description, &price, &quantity, &filename)
		levels = append(levels, map[string]string{"title": title, "description": description, "price": price, "quantity": quantity, "filename": filename})
	}
	return levels, nil
}

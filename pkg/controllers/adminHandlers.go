package controllers

import (
	"context"
	"crudapp/pkg/models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	adminTemp, _ = template.ParseGlob("views/admin/*.html")
)

func AdminIndex(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	username := r.FormValue("username")

	_, err := r.Cookie(username)
	if err == nil {
		http.Redirect(w, r, "admin/adminHome", http.StatusSeeOther)
		return
	}
	adminTemp.ExecuteTemplate(w, "adminIndex.html", nil)
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	db := models.OpenDb()
	defer func() {
		sqlDb, _ := db.DB()
		sqlDb.Close()
		fmt.Println("db closed")
	}()

	user := &models.User{}

	//search db for the entered username
	db.Where("username=?", username).First(&user)

	//check whether any admin with the entered username
	if !user.IsAdmin {
		data := map[string]any{
			"err": "not an admin",
		}
		adminTemp.ExecuteTemplate(w, "adminIndex.html", data)
		return
	}

	//compare password stored and entered
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		data := map[string]any{
			"err": "invalid password",
		}
		adminTemp.ExecuteTemplate(w, "adminIndex.html", data)
		return
	}

	session, _ := store.Get(r, username)

	session.Values[username] = username
	session.IsNew = true
	session.Options = &sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
		SameSite: 0,
	}
	session.Save(r, w)

	//open admin home page after validation
	// data := map[string]any{
	// 	"username": username,
	// }
	// adminTemp.ExecuteTemplate(w, "adminHome.html", data)
	http.Redirect(w, r, "/admin/adminHome", http.StatusSeeOther)

}

func AdminHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	// params := mux.Vars(r)
	// username := params["username"]
	username := r.FormValue("username")

	if c, err := r.Cookie(username); err == http.ErrNoCookie {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		http.SetCookie(w, c)
		return
	}
	data := map[string]interface{}{
		"username": username,
	}
	adminTemp.ExecuteTemplate(w, "adminHome.html", data)
}

func ViewUsers(w http.ResponseWriter, r *http.Request) {
	db := models.OpenDb()
	sqlDb, _ := db.DB()

	rows, _ := sqlDb.QueryContext(context.Background(), "SELECT id,name,username,email FROM users")

	var (
		id       int64
		name     string
		username string
		email    string
		data     []models.User
	)

	for rows.Next() {
		if err := rows.Scan(&id, &name, &username, &email); err != nil {
			fmt.Println(err)
		}
		data = append(data, models.User{
			Id:       id,
			Name:     name,
			Username: username,
			Email:    email,
		})
	}
	adminTemp.ExecuteTemplate(w, "viewUsers.html", data)
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	adminTemp.ExecuteTemplate(w, "createUser.html", nil)

}
func InsertUser(w http.ResponseWriter, r *http.Request) {
	if err := AddUser(w, r); err != nil {
		data := map[string]any{
			"err": "something went wrong  Try again",
		}
		adminTemp.ExecuteTemplate(w, "adminHome.html", data)
	}
	data := map[string]any{
		"err":  "Login now",
		"errr": "signup success",
	}
	adminTemp.ExecuteTemplate(w, "adminHome.html", data)
}
func AddUser(w http.ResponseWriter, r *http.Request) error {

	name := r.FormValue("name")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	newUser := models.User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: string(hashPass),
		IsAdmin:  false,
	}
	db := models.OpenDb()
	defer func() {
		sqlDb, _ := db.DB()
		sqlDb.Close()
		fmt.Println("db closed")
	}()

	user := &models.User{}
	db.AutoMigrate(user)

	result := db.Create(&newUser)

	return result.Error

}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	// deleteUser(int64(id))

	db := models.OpenDb()
	defer func() {
		sqlDb, _ := db.DB()
		sqlDb.Close()
		fmt.Println("db closed")
	}()

	user := &models.User{}
	db.Where("id = ?", id).Delete(&user)
	http.Redirect(w, r, "/admin/view", http.StatusSeeOther)

}
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	fmt.Println("admin logout")

	username := r.FormValue("username")

	c, err := r.Cookie(username)
	if err != nil {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	c.Value = ""
	c.Path = "/"
	c.MaxAge = -1
	http.SetCookie(w, c)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

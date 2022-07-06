package controllers

import (
	"crudapp/models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	userTemp, _ = template.ParseGlob("views/user/*.html")
	store       = sessions.NewCookieStore([]byte("secure-cookie"))
)

//user index page in the root path
func UserIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	username := r.FormValue("username")

	_, err := r.Cookie(username)
	if err == nil {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}
	userTemp.ExecuteTemplate(w, "userIndex.html", nil)
}

//login function in post method from index page
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" {
		data := map[string]any{
			"err": "enter username",
		}
		userTemp.ExecuteTemplate(w, "userIndex.html", data)
		return
	}

	Db := models.OpenDb()
	sqlDb, _ := Db.DB()
	defer func() {
		sqlDb.Close()
		fmt.Println("db closed")
	}()
	user := models.User{}

	//check in the db if the username exist
	Db.Where("username=?", username).First(&user)

	//validate entered username
	if username != user.Username {
		data := map[string]any{
			"err":  "user not found",
			"errr": "signup now",
		}
		// http.Redirect(w, r, "/", http.StatusSeeOther)
		userTemp.ExecuteTemplate(w, "userIndex.html", data)
		return
	}

	//getting stored password by decrypting hash value and compare with entered password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		data := map[string]any{
			"err": "invalid password",
		}
		userTemp.ExecuteTemplate(w, "userIndex.html", data)
		return
	}

	//adding cookie and session
	session, _ := store.Get(r, username)

	session.Values[username] = username
	session.Options = &sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
		SameSite: 0,
	}
	session.Save(r, w)

	data := map[string]interface{}{
		"username": username,
		"id":       user.Id,
	}

	//execute user home page after validation
	userTemp.ExecuteTemplate(w, "userHome.html", data)
	// http.Redirect(w, r, "/user", http.StatusSeeOther)
}

//UserHome function execute only if there is a valid cookie
func UserHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	username := r.FormValue("username")

	if c, err := r.Cookie(username); err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		http.SetCookie(w, c)
		return
	}
	user := &models.User{}
	db := models.OpenDb()
	defer func() {
		sqlDb, _ := db.DB()
		sqlDb.Close()
		fmt.Println("db closed")
	}()
	db.Where("username=?", username).First(&user)

	data := map[string]interface{}{
		"username": username,
		"id":       user.Id,
	}
	userTemp.ExecuteTemplate(w, "userHome.html", data)
}

//user edit get method,
//display user details and allow edits
func UserEdit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	user := &models.User{}
	db := models.OpenDb()
	defer func() {
		sqlDb, _ := db.DB()
		sqlDb.Close()
		fmt.Println("db closed")
	}()
	db.Where("id=?", id).First(&user)

	data := map[string]interface{}{
		"id":       user.Id,
		"name":     user.Name,
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
	}

	userTemp.ExecuteTemplate(w, "updateUser.html", data)

}

//user update post method,
//update edits into the database
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	uid, _ := strconv.Atoi(id)
	name := r.FormValue("name")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	hashPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{}

	db := models.OpenDb()
	defer func() {
		sqlDb, _ := db.DB()
		sqlDb.Close()
		fmt.Println("db closed")
	}()
	db.Where("id=?", uid).First(&user)

	db.Model(&user).Updates(map[string]any{"name": name, "username": username, "email": email, "password": hashPass})
	db.Save(&user)
}

//get user signup page
func SignUp(w http.ResponseWriter, r *http.Request) {
	userTemp.ExecuteTemplate(w, "signUp.html", nil)
}

//signup post method function,
//insert into the database is done by AddUser method
func UserSignUp(w http.ResponseWriter, r *http.Request) {
	err := AddUser(w, r)
	if err != nil {
		data := map[string]any{
			"err": "something went wrong Try again",
		}
		userTemp.ExecuteTemplate(w, "signUp.html", data)
	}
	data := map[string]any{
		"err":  "Login now",
		"errr": "signup success",
	}
	userTemp.ExecuteTemplate(w, "userIndex.html", data)
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	username := r.FormValue("username")

	c, err := r.Cookie(username)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	c.Value = ""
	c.Path = "/"
	c.MaxAge = -1
	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

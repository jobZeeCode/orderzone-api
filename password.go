package main 
import (
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
)

type passwordStruct struct{
	Email interface{}
	Password interface{}
}

//HashPassword : create hash from password and email
func HashPassword(password string, email string) (string, error) {
	var newPassword = password+email;
	bytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	return string(bytes), err 
}

//CheckPassWordHash : check password is correct ?
func CheckPassWordHash(password string, hash string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	//return err == nil
	return err == nil 
}

//PasswordHandler : manage about password
func PasswordHandler(w http.ResponseWriter, r *http.Request) {
	client, ctx = Conect()
	w.Header().Set("content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, content-type")
	w.Header().Set("Access-control-Allow-Origin", "http://localhost:3000")
	defer client.Close()
	switch r.Method {
	case "POST":
		//check login
		var pass passwordStruct 
		var user userStruct
		err := json.NewDecoder(r.Body).Decode(&pass)
		if err != nil {
			log.Fatalf("Fail Decoder : %v ", err)
		}
		var newPassword = fmt.Sprintf("%v",pass.Password)+fmt.Sprintf("%v",pass.Email)
		data, err := serchData(fmt.Sprintf("%v",pass.Email), "Email", fetchData("customers"))
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"msg": "not found"})
		} else {
			var jsonData []byte
			jsonData , err := json.Marshal(data)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"msg": "not found"})
			}else {
				err := json.Unmarshal(jsonData, &user)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"msg": "not found"})
				}
				res := CheckPassWordHash(newPassword, fmt.Sprintf("%v",user.Password))
				if res{
					json.NewEncoder(w).Encode(map[string]interface{}{
						"Email": user.Email,
						"ID": user.ID,
						"IsLogin" : res,
					})
				}else{ 
					json.NewEncoder(w).Encode(map[string]interface{}{
						"Email": user.Email,
						"ID": user.ID,
						"IsLogin" : res,
					})
				}
			}
		}
	}
}
package main

import (
	"crypto/sha256"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type Registrazione struct{
	Nome     string `form:"nome" binding:"required`
	Cognome     string `form:"cognome" binding:"required`
	DataNascita string `form:"dataNascita" binding:"required`
	Nickname string `form:"nickname" binding:"required`
	Password string `form:"password" binding:"required`
	ConfermaPassword string `form:"confermaPassword" binding:"required`
}

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Utente struct{
	Nome    string
	Cognome string
	DataNascita string
	Nickname string
	Password string
	Model
}


func cripto(psw string) string{
	myByte := sha256.Sum256([]byte(psw))
	myString := string(myByte[:])
	return myString
}

func aggiungiRecord(utente Utente){
	db, err := GetGormDBConnection()
	if err != nil{
		log.Panicln(err)
	}
	db.Create(&utente)
}

func eliminaRecord(id int){
	db, err := GetGormDBConnection()
	if err != nil{
		log.Panicln(err)
	}
	db.Delete(&Utente{}, id)
}

func checkLogin(user string, password string) bool{
	db, err := GetGormDBConnection()
	if err != nil{
		log.Panicln(err)
	}

	c := Utente{}
	db.Find(&c, "nickname = '"+ user +"'")

	if  user == c.Nickname && password == c.Password{
		return true
	}
	return false
}

func registraUtente(form Registrazione) bool{
	var utente Utente
	if len(form.Nome) < 3  && len(form.Cognome) < 3 && len(form.Nickname) < 3 && len(form.DataNascita) < 10 {
		return false
	}else {
		if form.Password==form.ConfermaPassword {
			utente.Nome = form.Nome
			utente.Cognome = form.Cognome
			utente.Nickname = form.Nickname
			utente.Password = form.Password
			utente.DataNascita = form.DataNascita
			aggiungiRecord(utente)
			return true
		}
	}
	return false



}

func init(){
	err := InitializeDBConnection()
	if err != nil{
		log.Panicln(err)
	}

	db, err := GetGormDBConnection()
	if err != nil{
		log.Panicln(err)
	}

	err = db.AutoMigrate(&Utente{})
	if err != nil{
		log.Panicln(err)
	}
}

func main() {
	router := gin.Default()

	router.POST("/loginForm", func(c *gin.Context) {
		var form Login	//creo una variabile con la struct fatta apposta per il login
		if err := c.ShouldBind(&form); err != nil { //prendo i campi del form e li salvo nella struct
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		giusto := checkLogin(form.User, form.Password) //controllo con la funzione che i campi inseriti siano corretti



		if giusto { //se è tutto corretto stampo che è acceduto con success
			c.String(200, "ACCESSO AUTORIZZATO")
		}
		if !giusto { //se l'utente ha sbagliato qualcosa stampo che l'accesso è stato negato
			c.String(200, "ACCESSO NEGATO")
		}
	})

	router.POST("/registerform", func(c *gin.Context) {
		var registrazione Registrazione //creo una variabile con la struct fatta apposta per la registrazione
		if err := c.ShouldBind(&registrazione); err != nil { //prendo i campi del form e li salvo nella struct
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !registraUtente(registrazione){ //se l'utente sbaglia qualcosa stampo gli eventuali problemi
			c.String(200, "LE PASSWORD NON COINCIDONO OPPURE HAI DIMENTICATO DI COMPILARE QUALCOSA")
		}else{ //se l'utente ha fatto tutto correttamente stampo che si è registrato con successo
			c.String(200, "REGISTRATO CON SUCCESSO")
		}
	})
	router.Run(":8080")
}


/*func main() {
	router := gin.Default()



	// Example for binding JSON ({"user": "manu", "password": "123"})
	/*router.POST("/loginJSON", func(c *gin.Context) {
		var json Login

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})*/

	// Example for binding XML (
	//	<?xml version="1.0" encoding="UTF-8"?>
	//	<root>
	//		<user>manu</user>
	//		<password>123</password>
	//	</root>)
	/*router.POST("/loginXML", func(c *gin.Context) {
		var xml Login
		if err := c.ShouldBindXML(&xml); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if xml.User != "manu" || xml.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})


	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		giusto := checkLogin(form.User, form.Password)



		if giusto {
			c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
		}
		if !giusto {
			c.JSON(http.StatusOK, gin.H{"status": "user o password wrong"})
		}
	})

	router.POST("/registerform", func(c *gin.Context) {
		var registrazione Registrazione
		if err := c.ShouldBind(&registrazione); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !registraUtente(registrazione){
			c.String(200, "Le password non coincidono")
		}else{
			c.String(200, "Registrato con successo")
		}
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}*/



/*func main() {				//custom log format
	router := gin.New()

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %s %s %s %s",
			param.ClientIP, //indirizzo ip del cliente
			param.TimeStamp.Format(time.RFC1123),	//il momento esatto in cui viene mandata la richiesta
			param.Method,	//il metodo della richiesta (GET POST ecc..)
			param.Path,	//il parametro del path
			param.Request.Proto,	//il protocollo
			param.StatusCode,	//stato della richiesta
			param.Latency,	//velocità della risposta
			param.Request.UserAgent(),	//specifiche teniche del cliente
			param.ErrorMessage,	//eventuali messaggi d'errore
		)
	}))
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Run(":8080")
}*/





/*func main() {			//How to write log file
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Run(":8080")
}*/






/*func main() {
	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}

	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}

	router.Run(":8080")
}*/










/*func main() {				//Map as querystring or postform parameters
	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {

		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")

		fmt.Printf("ids: %v; names: %v", ids, names)
	})
	router.Run(":8080")
}*/







/*func main() {					//Another example:query + post form
	router := gin.Default()

	router.POST("/post/:name/", func(c *gin.Context) {

		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		name2 := c.Query("name")
		name3 := c.Param("name")
		message := c.PostForm("message")

		fmt.Printf("id: %s; page: %s; name: %s, message: %s", id, page, name, name2, name3, message)
	})
	router.Run(":8080")
}*/






/*func main() {			//Multipart/Urlencoded Form
	router := gin.Default()

	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	router.Run(":8080")
}*/



/*func main() {         //Querystring parameters
	router := gin.Default()

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.Run(":8080")
}*/


/*func main() {          //Parameters in path
router := gin.Default()

// This handler will match /user/john but will not match /user/ or /user
router.GET("/user/:name", func(c *gin.Context) {
name := c.Param("name")
c.String(http.StatusOK, "Hello %s", name)
})

// However, this one will match /user/john/ and also /user/john/send
// If no other routers match /user/john, it will redirect to /user/john/
router.GET("/user/:name/*action", func(c *gin.Context) {
name := c.Param("name")
action := c.Param("action")
message := name + " is " + action
c.String(http.StatusOK, message)
})

// For each matched request Context will hold the route definition
router.POST("/user/:name/*action", func(c *gin.Context) {
	name := c.Param("name")
	action := c.Param("action")
	c.String(http.StatusOK, "qualcosa in post " + name + action)
})

router.POST("/user", func(c *gin.Context) {
	name := c.Param("name")
	action := c.Param("action")
	c.String(http.StatusOK, "qualcosa in post " + name + action)})

// This handler will add a new router for /user/groups.
// Exact routes are resolved before param routes, regardless of the order they were defined.
// Routes starting with /user/groups are never interpreted as /user/:name/... routes
router.GET("/user/groups", func(c *gin.Context) {
	c.String(http.StatusOK, "The available groups are: \n[GIN-debug] GET    /user/:name               --> main.main.func1 (3 handlers)\n[GIN-debug] GET    /user/:name/*action       --> main.main.func2 (3 handlers)\n[GIN-debug] POST   /user/:name/*action       --> main.main.func3 (3 handlers)\n[GIN-debug] GET    /user/groups")
})

router.Run(":8080")
}*/
package main

import (
	"encoding/json"
	"github.com/johnidm/owl-links-api/database"	
	"github.com/johnidm/owl-links-api/controllers"
	"github.com/johnidm/owl-links-api/utils"	
	"github.com/julienschmidt/httprouter"	
	"net/http"
	"strconv"
	"flag"
	"log"
	"errors"
	"io"
  	"io/ioutil"
  	"gopkg.in/mgo.v2"
//  	"fmt" 
)

var db = database.OpenDB()

var collecttion *mgo.Collection

func main() {

	port := flag.String("port", "8000", "HTTP Port")
	flag.Parse()

	sess, err := controllers.Connect()
	if err != nil {
		panic(err)
	}

	collecttion = sess.DB("owl-links").C("links")

	defer db.Close()
	defer sess.Close()

	router := httprouter.New()

	router.GET("/", RunProject)

	router.GET("/links", GetLinks)
	router.GET("/link/:id", GetLink)
	router.DELETE("/link/:id", DeleteLink)
	router.PUT("/link/", PutLink)
	router.POST("/link", PostLink)

	router.GET("/contacts", GetContacts)
	router.GET("/contact/:id", GetContact)
	router.DELETE("/contact/:id", DeleteContact)

	router.GET("/collectlinks", GetCollectlinks)	
	router.DELETE("/collectlink/:id", DeleteCollectlink)

	router.GET("/newslatters", GetNewslatters)	
	router.DELETE("/newslatter/:id", DeleteNewslatter)

	log.Println("Stating Server on ", *port)
	log.Fatal(http.ListenAndServe(":" + *port, router))
}

func RunProject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.Write([]byte("<h2><font color=\"green\">Owl Link API is running!</font></h2>"))
}
    
func GetLinks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if !utils.APIKeyIsValid(w, r) {
		return
	}
		
	links, err := controllers.GetLinks(collecttion)
	if err != nil {
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.GetLinks")
		return
	}

	js, err := json.Marshal(links)
	if err != nil {
		utils.DefineReturnRequestError(w, err, "Erro ao fazer o marshal dos links",)		
		return
	}	

	utils.WriteJson(w, js)	
}

func PutLink(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if !utils.APIKeyIsValid(w, r) {
		return
	}

	defer r.Body.Close()
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {            
      utils.DefineReturnRequestError(w, err, "Falha ao ler os dados da requisição")	
      return
    }

    var link controllers.Link

    if err := json.Unmarshal(body, &link); err != nil {
      utils.DefineReturnRequestError(w, err, "Falha fazer unmarshal dos dados da requisição")
      return	
    }

    err = controllers.UpdateLink(&link, collecttion)
    if err != nil {		
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.UpdateLink")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func PostLink(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if !utils.APIKeyIsValid(w, r) {
		return
	}

	defer r.Body.Close()
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {            
      utils.DefineReturnRequestError(w, err, "Falha ao ler os dados da requisição")	
      return
    }

    var link controllers.Link

    if err := json.Unmarshal(body, &link); err != nil {
      utils.DefineReturnRequestError(w, err, "Falha fazer unmarshal dos dados da requisição")
      return	
    }

    err = controllers.CreateLink(&link, collecttion)
    if err != nil {		
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.CreateLink")
		return
	}
		
	w.WriteHeader(http.StatusCreated)

}

func GetLink(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if !utils.APIKeyIsValid(w, r) {
		return
	}
		
	_id := p.ByName("id")
	if _id == "" {
		utils.DefineReturnRequestIdInvalid(w, errors.New("Id is blank"))	
		return
	}
	

	link, err := controllers.GetLink(_id, collecttion)
	if err != nil {		
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.GetLink")
		return
	}

	js, err := json.Marshal(link)
	if err != nil {
		utils.DefineReturnRequestError(w, err, "Erro ao fazer o marshal do link")					
		return
	}

	utils.WriteJson(w, js)
	
}

func DeleteLink(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if !utils.APIKeyIsValid(w, r) {
		return
	}

	_id := p.ByName("id")
	if _id == "" {
		utils.DefineReturnRequestIdInvalid(w, errors.New("Id is blank"))	
		return
	}

	err := controllers.DeleteLink(_id, collecttion)
	if err != nil {
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.DeleteLink")
		return
	}

	w.WriteHeader(http.StatusOK)
	
}

func GetContacts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !utils.APIKeyIsValid(w, r) {
		return
	}
	
	contacts, err := controllers.GetContacts(db)
	if err != nil {
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.GetContacts")
		return
	}

	js, err := json.Marshal(contacts)
	if err != nil {
		utils.DefineReturnRequestError(w, err, "Erro ao fazer o marshal dos contatos",)		
		return
	}	

	utils.WriteJson(w, js)
}

func GetContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !utils.APIKeyIsValid(w, r) {
		return
	}

	id, err := strconv.ParseInt(p.ByName("id"), 0, 64)
	if err != nil {
		utils.DefineReturnRequestIdInvalid(w, err)	
		return
	}

	contact, err := controllers.GetContact(db, id)
	if err != nil {		
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.GetContact")
		return
	}

	js, err := json.Marshal(contact)
	if err != nil {
		utils.DefineReturnRequestError(w, err, "Erro ao fazer o marshal do contato")					
		return
	}

	utils.WriteJson(w, js)
}


func DeleteContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	
	if !utils.APIKeyIsValid(w, r) {
		return
	}	

	id, err := strconv.ParseInt(p.ByName("id"), 0, 64)
	if err != nil {
		utils.DefineReturnRequestIdInvalid(w, err)
		return
	}

	err = controllers.DeleteContact(db, id)
	if err != nil {
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.DeleteContact")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetCollectlinks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if !utils.APIKeyIsValid(w, r) {
		return
	}

	collectlinks, err := controllers.GetCollectlinks(db)
	if err != nil {
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.GetCollectlinks")
		return
	}

	js, err := json.Marshal(collectlinks)
	if err != nil {
		utils.DefineReturnRequestError(w, err, "Erro ao fazer o marshal das sugestões de link",)		
		return
	}	

	utils.WriteJson(w, js)

}

func DeleteCollectlink(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if !utils.APIKeyIsValid(w, r) {
		return
	}

	id, err := strconv.ParseInt(p.ByName("id"), 0, 64)
	if err != nil {
		utils.DefineReturnRequestIdInvalid(w, err)
		return
	}

	err = controllers.DeleteCollectlink(db, id)
	if err != nil {
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.DeleteCollectlink")
		return
	}

	w.WriteHeader(http.StatusOK)

}


func GetNewslatters(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if !utils.APIKeyIsValid(w, r) {
		return
	}

	newslatters, err := controllers.GetNewslatters(db)
	if err != nil {
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.GetNewslatters")
		return
	}

	js, err := json.Marshal(newslatters)
	if err != nil {
		utils.DefineReturnRequestError(w, err, "Erro ao fazer o marshal das assinaturas de newslatters",)		
		return
	}	

	utils.WriteJson(w, js)
}

func DeleteNewslatter(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if !utils.APIKeyIsValid(w, r) {
		return
	}

	id, err := strconv.ParseInt(p.ByName("id"), 0, 64)
	if err != nil {
		utils.DefineReturnRequestIdInvalid(w, err)
		return
	}

	err = controllers.DeleteNewslatter(db, id)
	if err != nil {
		utils.DefineReturnRequestFailExecFunc(w, err, "controllers.DeleteNewslatter")
		return
	}

	w.WriteHeader(http.StatusOK)

}


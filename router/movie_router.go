package movierouter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/programadriano/go-restapi/config/dao"
	helper "github.com/programadriano/go-restapi/config/helper"
	redis "github.com/programadriano/go-restapi/config/redis"
	. "github.com/programadriano/go-restapi/models"
	"gopkg.in/mgo.v2/bson"
)

var dao = MoviesDAO{}

func GetAll(w http.ResponseWriter, r *http.Request) {
	movies := []Movie{}
	reply, err := redis.Get("movies")

	if err != nil {
		fmt.Println("Buscando no mongoDB")
		movies, err := dao.GetAll()
		helper.HandleError(err)
		m, err := json.Marshal(movies)
		helper.HandleError(err)
		redis.Set("movies", []byte(m))
		helper.RespondWithJson(w, http.StatusOK, movies)

	} else {
		fmt.Println("Buscando no redis")
		json.Unmarshal(reply, &movies)
		helper.RespondWithJson(w, http.StatusOK, movies)
	}

}

func GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := dao.GetByID(params["id"])
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Movie ID")
		return
	}
	helper.RespondWithJson(w, http.StatusOK, movie)
}

func Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	movie.ID = bson.NewObjectId()
	if err := dao.Create(movie); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusCreated, movie)
}

func Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(params["id"], movie); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": movie.Name + " atualizado com sucesso!"})
}

func Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	if err := dao.Delete(params["id"]); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

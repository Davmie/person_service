package delivery

import (
	"encoding/json"
	"fmt"
	personUseCase "github.com/Davmie/person_service/internal/person/usecase"
	"io"
	"net/http"
	"strconv"

	"github.com/Davmie/person_service/models"
	"github.com/Davmie/person_service/pkg/logger"
	//"github.com/asaskevich/govalidator"
)

type PersonHandler struct {
	PersonUseCase personUseCase.PersonUseCaseI
	Logger        logger.Logger
}

func (ah *PersonHandler) Create(w http.ResponseWriter, r *http.Request) {
	person := models.Person{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ah.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ah.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &person)
	if err != nil {
		ah.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	//_, err = govalidator.ValidateStruct(person)
	//if err != nil {
	//	ah.Logger.Infow("can`t validate form",
	//		"err:", err.Error())
	//	http.Error(w, "bad data", http.StatusBadRequest)
	//	return
	//}

	err = ah.PersonUseCase.Create(&person)
	if err != nil {
		ah.Logger.Infow("can`t create person",
			"err:", err.Error())
		http.Error(w, "can`t create person", http.StatusBadRequest)
		return
	}

	//resp, err := json.Marshal(person)
	//
	//if err != nil {
	//	ah.Logger.Errorw("can`t marshal person",
	//		"err:", err.Error())
	//	http.Error(w, "can`t make person", http.StatusInternalServerError)
	//	return
	//}

	w.Header().Set("Location", fmt.Sprintf("/api/v1/persons/%d", person.ID))
	w.WriteHeader(http.StatusCreated)

	//_, err = w.Write(resp)
	//if err != nil {
	//	ah.Logger.Errorw("can`t write response",
	//		"err:", err.Error())
	//	http.Error(w, "can`t write response", http.StatusInternalServerError)
	//	return
	//}
}

func (ah *PersonHandler) Get(w http.ResponseWriter, r *http.Request) {
	personIdString := r.PathValue("personId")
	if personIdString == "" {
		ah.Logger.Errorw("no personId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	personId, err := strconv.Atoi(personIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	person, err := ah.PersonUseCase.Get(personId)
	if err != nil {
		ah.Logger.Infow("can`t get person",
			"err:", err.Error())
		http.Error(w, "can`t get person", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(person)

	if err != nil {
		ah.Logger.Errorw("can`t marshal person",
			"err:", err.Error())
		http.Error(w, "can`t make person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (ah *PersonHandler) Update(w http.ResponseWriter, r *http.Request) {
	personIdString := r.PathValue("personId")
	if personIdString == "" {
		ah.Logger.Errorw("no personId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	personId, err := strconv.Atoi(personIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	person := &models.Person{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ah.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ah.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, person)
	if err != nil {
		ah.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	//_, err = govalidator.ValidateStruct(person)
	//if err != nil {
	//	ah.Logger.Infow("can`t validate form",
	//		"err:", err.Error())
	//	http.Error(w, "bad data", http.StatusBadRequest)
	//	return
	//}

	person.ID = personId
	err = ah.PersonUseCase.Update(person)
	if err != nil {
		ah.Logger.Infow("can`t update person",
			"err:", err.Error())
		http.Error(w, "can`t update person", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(person)

	if err != nil {
		ah.Logger.Errorw("can`t marshal person",
			"err:", err.Error())
		http.Error(w, "can`t make person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (ah *PersonHandler) Delete(w http.ResponseWriter, r *http.Request) {
	personIdString := r.PathValue("personId")
	if personIdString == "" {
		ah.Logger.Errorw("no personId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	personId, err := strconv.Atoi(personIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	err = ah.PersonUseCase.Delete(personId)
	if err != nil {
		ah.Logger.Infow("can`t delete person",
			"err:", err.Error())
		http.Error(w, "can`t delete person", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ah *PersonHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	persons, err := ah.PersonUseCase.GetAll()
	if err != nil {
		ah.Logger.Infow("can`t get all persons",
			"err:", err.Error())
		http.Error(w, "can`t get all persons", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(persons)
	if err != nil {
		ah.Logger.Errorw("can`t marshal person",
			"err:", err.Error())
		http.Error(w, "can`t make person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

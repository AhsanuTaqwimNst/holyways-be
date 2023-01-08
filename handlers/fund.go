package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	funddto "holyways/dto/fund"
	dto "holyways/dto/result"
	"holyways/models"
	"holyways/repositories"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerFund struct {
	FundRepository repositories.FundRepository
}

func HandlerFund(FundRepository repositories.FundRepository) *handlerFund {
	return &handlerFund{FundRepository}
}

func (h *handlerFund) FindFunds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fund, err := h.FundRepository.FindFunds()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: fund}
	json.NewEncoder(w).Encode(response)
}
func (h *handlerFund) GetFund(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var fund models.Fund
	fund, err := h.FundRepository.GetFund(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: fund}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFund) GetFundByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var fund []models.Fund
	fund, err := h.FundRepository.GetFundByUser(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: fund}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFund) CreateFund(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64)) // ini apa ya

	dataContex := r.Context().Value("dataFile")
	filepath := dataContex.(string)

	goal, _ := strconv.Atoi(r.FormValue("goal"))

	request := funddto.CreateFundRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Goal:        goal,
		UserID:      userId,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	//Add Cloudinary
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	//Upload file to Cloudinary oks
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "holy_ways"})

	if err != nil {
		fmt.Println(err.Error())
	}

	status := "Running"

	fund := models.Fund{
		Title:       request.Title,
		Image:       resp.SecureURL,
		Goal:        request.Goal,
		Description: request.Description,
		Status:      status,
		UserID:      userId,
	}

	fund, err = h.FundRepository.CreateFund(fund)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "error"}

		json.NewEncoder(w).Encode(response)
		return
	}

	fund, _ = h.FundRepository.GetFund(fund.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: fund}
	json.NewEncoder(w).Encode(response)
} //holy_ways

func (h *handlerFund) UpdateFund(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(funddto.UpdateFundRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "err.Error()"}
		json.NewEncoder(w).Encode(response)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	fund, err := h.FundRepository.GetFund(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "aa.Error()"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Status != "" {
		fund.Status = request.Status
	}

	data, err := h.FundRepository.UpdateFund(fund, id)
	fmt.Println("ini fund", fund.Status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "err.sError()"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFund) DeleteFund(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fund, err := h.FundRepository.GetFund(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.FundRepository.DeleteFund(id, fund)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}

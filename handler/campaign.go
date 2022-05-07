package handler
import (
	"bwastartup/user"
	"bwastartup/helper"
	"bwastartup/campaign"
	"strconv"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)
// tangkap parameter di handler
// handler ke service
// service menentukan repository mana yang di panggil
// repository : GetAll, GetByUserId
// db

type campaignHandler struct{
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context)  {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err:= h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) GetCampaign(c *gin.Context)  {
	// handler : mapping id yang di url ke struct input => service, call formatter
	// service : inputnya struct input => menangkap id di url, manggil repo
	// repo : get campaign by id
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err:= h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign Detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) CreateCampaign(c *gin.Context)   {
	// tangkap parameter dari user ke input struct 
	// ambil current user dari jwt
	// panggil service, parameternya input struct (dan juga untuk slug)
	// panggil repository untuk simpan data campaign baru
	var input campaign.CreateCampaignInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil{

		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context)   {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil{

		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}
func (h *campaignHandler) UploadImage(c *gin.Context)   {
// handler tangkap input dan ubah ke struct input dan save image campaign ke suatu folder
// service (kondisi manggil point 2 di repo, panggil repo point 1)
// repository : 
// 1. create images/data image ke dalam tabel campaign_images
// 2. ubah is_primary true ke false ketika ada yang ubah is_primary
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID 

	file, err := c.FormFile("file")

		if err != nil{
			data := gin.H{
				"is_uploaded": false,
			}

			response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		
		path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	
		err = c.SaveUploadedFile(file, path)
		if err != nil{
			data := gin.H{
				"is_uploaded": false,
			}
	
			response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		_, err = h.service.SaveCampaignImage(input, path)
		if err != nil{
			data := gin.H{
				"is_uploaded": false,
			}
	
			response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	
		data := gin.H{
			"is_uploaded": true,
		}
	
		response := helper.APIResponse("Campaign Image successfully uploaded", http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
}
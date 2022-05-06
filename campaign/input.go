package campaign

import (
	"bwastartup/user"
)

type GetCampaignDetailInput struct{
	ID int `uri:"id" binding:"required"`
}

type CreateCampaignInput struct{
	// ID 				 int `json:"id" binding:"required"` 
	Name 			 string `json:"name" binding:"required"` 
	ShortDescription string `json:"short_description" binding:"required"` 
	Description string `json:"description" binding:"required"` 
	GoalAmount 		 int `json:"goal_amount" binding:"required"` 
	Perks 			 string `json:"perks" binding:"required"` 
	User user.User
}
package campaign
import "gorm.io/gorm"

// membuat kontrak
type Repository interface{
	FindAll() ([]Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
	FindByID(ID int) (Campaign, error)
}

//mendefinisikan struct(private) yang punya akses ke semua database, yaitu db 
type repository struct{ 
	db *gorm.DB
}

// buat func agar bisa diakses diluar file kita buat instancenya
func NewRepository(db *gorm.DB) *repository  {
	// kita return pembuatan instance baru si repository di mana kita parsing nilai db
	return &repository{db}
}

// definisikan masing2 kontrak yang sudah dibuat di atas
func (r *repository) FindAll() ([]Campaign, error)  {
	var campaigns []Campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns,err
	}
	
	return campaigns, nil
}

func (r *repository) FindByUserID(userID int) ([]Campaign, error)  {
	var campaigns []Campaign
	err := r.db.Where("user_id =?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns,err
	}
	
	return campaigns, nil
}

func (r *repository) FindByID(ID int) (Campaign, error)  {
	var campaign Campaign
	err := r.db.Preload("User").Preload("CampaignImages").Where("id =?", ID).Find(&campaign).Error
	if err != nil {
		return campaign,err
	}
	
	return campaign, nil
}




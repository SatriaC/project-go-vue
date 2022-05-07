package transaction
import "gorm.io/gorm"

// membuat kontrak
type Repository interface{
	GetByCampaignID(campaignID int) ([]Transaction, error)
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

func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error){
	var transactions []Transaction
	err := r.db.Preload("User").Where("campaign_id =?", campaignID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions,err
	}
	
	return transactions, nil
}
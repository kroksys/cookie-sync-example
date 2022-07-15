package service

import "fmt"

type MatchingTable struct {
	PartnersUserID string `gorm:"primaryKey"`
	PartnersName   string
	LocalUserID    string
}

func SaveMatchingTable(partnersName, partnerUser, localUser string) error {
	if partnersName == "" {
		return fmt.Errorf("CreateMatchingTable error: empty partnersName provided")
	}
	return DB.Save(&MatchingTable{
		PartnersName:   partnersName,
		PartnersUserID: partnerUser,
		LocalUserID:    localUser,
	}).Error
}

// Finds matching table record for given cookie name
func ReadMatchingTable(partnersName string) (MatchingTable, error) {
	var mtRecord MatchingTable
	err := DB.First(&mtRecord, "partners_name = ?", partnersName).Error
	return mtRecord, err
}

func ReadMatchingTableByLocalID(localUserID string) (MatchingTable, error) {
	var mtRecord MatchingTable
	err := DB.First(&mtRecord, "local_user_id = ?", localUserID).Error
	return mtRecord, err
}

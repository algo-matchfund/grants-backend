package database

func (db *GrantsDatabase) GetDonationCount() (int64, error) {
	var count int64

	err := db.QueryRow("SELECT COUNT(*) FROM contributions").Scan(&count)

	return count, err
}

func (db *GrantsDatabase) GetTotalDonationAmount() (int64, error) {
	var sum int64

	err := db.QueryRow("SELECT SUM(amount) FROM contributions").Scan(&sum)

	return sum, err
}

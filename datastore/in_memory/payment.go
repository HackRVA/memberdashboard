package in_memory

import "memberserver/api/models"

func (i In_memory) GetPayments() ([]models.Payment, error) {
	return []models.Payment{}, nil
}
func (i In_memory) AddPayment(payment models.Payment) error {
	return nil
}
func (i In_memory) AddPayments(payments []models.Payment) error {
	return nil
}
func (i In_memory) SetMemberLevel(memberId string, level models.MemberLevel) error {
	return nil
}
func (i In_memory) ApplyMemberCredits() {}
func (i In_memory) UpdateMemberTiers()  {}
func (i In_memory) GetPastDueAccounts() []models.PastDueAccount {
	return []models.PastDueAccount{}
}

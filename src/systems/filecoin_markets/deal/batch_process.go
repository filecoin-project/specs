package deal

type StorageDealPaymentAction = int

const (
	ExpireStorageDeals StorageDealPaymentAction = 0
	CreditStorageDeals StorageDealPaymentAction = 1
)

type StorageDealSlashAction = int

const (
	SlashDeclaredFaults   StorageDealSlashAction = 0
	SlashDetectedFaults   StorageDealSlashAction = 1
	SlashTerminatedFaults StorageDealSlashAction = 2
)

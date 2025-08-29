package types

type RecordStatus string

const (
    StatusSent      RecordStatus = "تم الإرسال"
    StatusReviewing RecordStatus = "جاري المراجعة"
    StatusApproved  RecordStatus = "تمت الموافقة"
)

type InformationStatus string
const (
    CurrentReviewing      InformationStatus = "inProgress"
    InfoStatusOK          InformationStatus = "accepted"
    InfoStatusCancel      InformationStatus = "rejected"
)

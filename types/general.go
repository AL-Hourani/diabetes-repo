package types

type RecordStatus string

const (
    StatusSent      RecordStatus = "تم الإرسال"
    StatusReviewing RecordStatus = "جاري المراجعة"
    StatusApproved  RecordStatus = "تمت الموافقة"
)



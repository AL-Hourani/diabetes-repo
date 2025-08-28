package types

type RecordStatus string

const (
    StatusSent      RecordStatus = "تم الإرسال"
    StatusReviewing RecordStatus = "جاري المراجعة"
    StatusApproved  RecordStatus = "تمت الموافقة"
)

type InformationStatus string
const (
    CurrentReviewing      InformationStatus = "جاري المعالجة "
    InfoStatusOK          InformationStatus = "مقبول"
    InfoStatusCancel      InformationStatus = "مرفوض"
)

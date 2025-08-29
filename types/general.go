package types

type RecordStatus string

const (
    StatusReviewing RecordStatus = "inProgress"
    StatusApproved  RecordStatus = "accepted"
    StatusRejected  RecordStatus = "rejected"
)

type InformationStatus string
const (
    CurrentReviewing      InformationStatus = "inProgress"
    InfoStatusOK          InformationStatus = "accepted"
    InfoStatusCancel      InformationStatus = "rejected"
)

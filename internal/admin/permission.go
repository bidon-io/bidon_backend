package admin

type ResourcePermissions struct {
	Read   bool `json:"read"`
	Create bool `json:"create"`
}

type ResourceInstancePermissions struct {
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

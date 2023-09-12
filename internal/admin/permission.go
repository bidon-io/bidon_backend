package admin

type ResourcePermission struct {
	ResourceName string  `json:"resource_name"`
	Path         string  `json:"path"`
	Actions      Actions `json:"actions"`
}

type Actions struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

func GetPermissions(authCtx AuthContext) []ResourcePermission {
	if authCtx.IsAdmin() {
		return getAdminPermissions()
	}
	return getUserPermissions()
}

func getAdminPermissions() []ResourcePermission {
	return []ResourcePermission{
		{
			ResourceName: "App",
			Path:         "/apps",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "App Demand Profile",
			Path:         "/app_demand_profiles",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "Auction Configuration",
			Path:         "/auction_configurations",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "Demand Source Account",
			Path:         "/demand_source_accounts",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "Line Item",
			Path:         "/line_items",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "Segment",
			Path:         "/segments",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "Country",
			Path:         "/countries",
			Actions: Actions{
				Read: true,
			},
		},
		{
			ResourceName: "User",
			Path:         "/users",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "Demand Source",
			Path:         "/demand_sources",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
	}
}

func getUserPermissions() []ResourcePermission {
	return []ResourcePermission{
		{
			ResourceName: "App",
			Path:         "/apps",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "App Demand Profile",
			Path:         "/app_demand_profiles",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "Auction Configuration",
			Path:         "/auction_configurations",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "Demand Source Account",
			Path:         "/demand_source_accounts",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "Line Item",
			Path:         "/line_items",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "Segment",
			Path:         "/segments",
			Actions: Actions{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		{
			ResourceName: "Country",
			Path:         "/countries",
			Actions: Actions{
				Read: true,
			},
		},
		{
			ResourceName: "Demand Source",
			Path:         "/demand_sources",
			Actions: Actions{
				Read: true,
			},
		},
	}
}

//go:generate ../../../bin/threeport-codegen api-model --filename $GOFILE --package $GOPACKAGE
package v0

// Profile is a named standard configuration for a definition object.
type Profile struct {
	Common `swaggerignore:"true" mapstructure:",squash"`

	// The name of a profile
	Name *string `json:"Name,omitempty" query:"name" gorm:"not null" validate:"required"`

	// Required if no CompanyID.  The user that owns the object.
	UserID *uint `json:"UserID,omitempty" query:"userid" validate:"optional"`

	// Required if no UserID.  The company that owns the object.
	CompanyID *uint `json:"CompanyID,omitempty" query:"companyid" validate:"optional"`
}

// Tier is a level of criticality for access control.  Common tiers would be
// "development" and "production" whereby typically many users will have access
// to manage development tiers while only leads and managers have access to
// manage production tier resources.
type Tier struct {
	Common `swaggerignore:"true" mapstructure:",squash"`

	Name *string `json:"Name,omitempty" query:"name" gorm:"not null" validate:"required"`

	// Required if no CompanyID.  The user that owns the object.
	UserID *uint `json:"UserID,omitempty" query:"userid" validate:"optional"`

	// Required if no UserID.  The company that owns the object.
	CompanyID *uint `json:"CompanyID,omitempty" query:"companyid" validate:"optional"`

	// The relative rank of criticality between tiers.  The higher the number,
	// the greater the criticality.  For example, a development tier could have
	// a criticality value of 10 while production could be 100.  Access control
	// can then use this criticality value to determine user access.
	Criticality *int32 `json:"Criticality,omitempty" query:"criticality" gorm:"not null" validate:"required"`
}

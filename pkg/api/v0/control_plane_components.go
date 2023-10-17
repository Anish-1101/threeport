// +threeport-codegen route-exclude
package v0

type ControlPlaneComponent struct {
	Common `swaggerignore:"true" mapstructure:",squash"`

	// The control plane instance ID that this component belongs to
	ControlPlaneInstanceID *uint `json:"ControlPlaneInstanceID,omitempty" gorm:"not null" validate:"required"`

	// Indicate whether the component is enabled to be deployed. Currently only respected by controllers
	Enabled *bool `json:"Enabled,omitempty" yaml:"Enabled" validate:"optional"`

	// The name of the component
	Name string `json:"Name,omitempty" yaml:"Name" query:"name" gorm:"not null" validate:"required"`

	// The image name of the component
	ImageName string `json:"ImageName,omitempty" yaml:"ImageName" query:"imagename" validate:"optional"`

	// The image repo of the component
	ImageRepo string `json:"ImageRepo,omitempty" yaml:"ImageRepo" query:"imagerepo" validate:"optional"`

	// The image tag of the component
	ImageTag string `json:"ImageTag,omitempty" yaml:"ImageTag" query:"imagetag" validate:"optional"`

	// The service account name to use when deploying
	ServiceAccountName string `json:"ServiceAccountName,omitempty" yaml:"ServiceAccountName" query:"serviceaccountname" validate:"optional"`

	// The service resource name to use when deploying
	ServiceResourceName string `json:"ServiceResourceName,omitempty" yaml:"ServiceResourceName" query:"serviceresourcename" validate:"optional"`
}
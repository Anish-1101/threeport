//go:generate ../../../bin/threeport-codegen api-model --filename $GOFILE --package $GOPACKAGE
package v0

// EthereumNode is a dependency for implementing Ethereum RPC support.
type EthereumNodeDefinition struct {
	Common     `swaggerignore:"true" mapstructure:",squash"`
	Definition `mapstructure:",squash"`

	// Network to join (e.g. mainnet, ropsten, rinkeby, goerli, kovan, etc.).
	Network *string `json:"Network,omitempty" query:"network" validate:"optional"`
}

type EthereumNodeInstance struct {
	Common   `swaggerignore:"true" mapstructure:",squash"`
	Instance `mapstructure:",squash"`

	//Name *string `json:"Name,omitempty" query:"name" gorm:"not null" validate:"required"`

	EthereumNodeDefinitionID *uint `json:"EthereumNodeDefinitionID,omitempty" validate:"optional,association"`

	ClusterInstanceID *uint `json:"ClusterInstanceID,omitempty" validate:"optional,association"`
}

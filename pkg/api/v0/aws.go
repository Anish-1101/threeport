//go:generate ../../../bin/threeport-codegen api-model --filename $GOFILE --package $GOPACKAGE
//go:generate ../../../bin/threeport-codegen controller --filename $GOFILE
package v0

import (
	"gorm.io/datatypes"
)

// AwsAccount is a user account with the AWS service provider.
type AwsAccount struct {
	Common `swaggerignore:"true" mapstructure:",squash"`

	// The unique name of an AWS account.
	Name *string `json:"Name,omitempty" query:"name" gorm:"not null" validate:"required"`

	// The account ID for the AWS account.
	AccountID *string `json:"AccountID,omitempty" query:"accountid" gorm:"not null" validate:"required"`

	// If true is the AWS Account used if none specified in a definition.
	DefaultAccount *bool `json:"DefaultAccount,omitempty" query:"defaultaccount" gorm:"default:false" validate:"optional"`

	// The region to use for AWS managed services if not specified.
	DefaultRegion *string `json:"DefaultRegion,omitempty" query:"defaultregion" gorm:"not null" validate:"required"`

	// The access key ID credentials for the AWS account.
	AccessKeyID *string `json:"AccessKeyID,omitempty" gorm:"not null" validate:"required" encrypt:"true"`

	// The secret key credentials for the AWS account.
	SecretAccessKey *string `json:"SecretAccessKey,omitempty" gorm:"not null" validate:"required" encrypt:"true"`

	// The cluster instances deployed in this AWS account.
	AwsEksKubernetesRuntimeDefinitions []*AwsEksKubernetesRuntimeDefinition `json:"AwsEksKubernetesRuntimeDefinitions,omitempty" validate:"optional,association"`
}

// AwsEksKubernetesRuntimeDefinition provides the configuration for EKS cluster instances.
type AwsEksKubernetesRuntimeDefinition struct {
	Common     `swaggerignore:"true" mapstructure:",squash"`
	Definition `mapstructure:",squash"`

	// The AWS account in which the EKS cluster is provisioned.
	AwsAccountID *uint `json:"AWSAccountID,omitempty" query:"awsaccountid" gorm:"not null" validate:"required"`

	// TODO: add fields for region limitations
	// RegionsAllowed
	// RegionsForbidden

	// The number of zones the cluster should span for availability.
	ZoneCount *int `json:"ZoneCount,omitempty" query:"zonecount" gorm:"not null" validate:"required"`

	// The AWS instance type for the default initial node group.
	DefaultNodeGroupInstanceType *string `json:"DefaultNodeGroupInstanceType,omitempty" query:"defaultnodegroupinstancetype" gorm:"not null" validate:"required"`

	// The number of nodes in the default initial node group.
	DefaultNodeGroupInitialSize *int `json:"DefaultNodeGroupInitialSize,omitempty" query:"defaultnodegroupinitialsize" gorm:"not null" validate:"required"`

	// The minimum number of nodes the default initial node group should have.
	DefaultNodeGroupMinimumSize *int `json:"DefaultNodeGroupMinimumSize,omitempty" query:"defaultnodegroupminimumsize" gorm:"not null" validate:"required"`

	// The maximum number of nodes the default initial node group should have.
	DefaultNodeGroupMaximumSize *int `json:"DefaultNodeGroupMaximumSize,omitempty" query:"defaultnodegroupmaximumsize" gorm:"not null" validate:"required"`

	// The AWS EKS kubernetes runtime instances derived from this definition.
	AwsEksKubernetesRuntimeInstances []*AwsEksKubernetesRuntimeInstance `json:"AwsEksKubernetesRuntimeInstances,omitempty" validate:"optional,association"`

	// The kubernetes runtime definition for an EKS cluster in AWS.
	KubernetesRuntimeDefinitionID *uint `json:"KubernetesRuntimeDefinitionID,omitempty" query:"kubernetesruntimedefinitionid" gorm:"not null" validate:"required"`
}

// +threeport-codegen:reconciler
// AwsEksKubernetesRuntimeInstance is a deployed instance of an EKS cluster.
type AwsEksKubernetesRuntimeInstance struct {
	Common         `swaggerignore:"true" mapstructure:",squash"`
	Instance       `mapstructure:",squash"`
	Reconciliation `mapstructure:",squash"`

	// The AWS Region in which the cluster is provisioned.  This field is
	// stored in the instance (as well as definition) since a change to the
	// definition will not move a cluster.
	Region *string `json:"Region,omitempty" query:"region" validate:"optional"`

	// The definition that configures this instance.
	AwsEksKubernetesRuntimeDefinitionID *uint `json:"AwsEksKubernetesRuntimeDefinitionID,omitempty" query:"awsekskubernetesruntimedefinitionid" gorm:"not null" validate:"required"`

	// An inventory of all AWS resources for the EKS cluster.
	ResourceInventory *datatypes.JSON `json:"ResourceInventory,omitempty" validate:"optional"`

	// The kubernetes runtime instance associated with the AWS EKS cluster.
	KubernetesRuntimeInstanceID *uint `json:"KubernetesRuntimeInstanceID,omitempty" query:"kubernetesruntimeinstanceid" gorm:"not null" validate:"required"`
}

// AwsRelationalDatabaseDefinition is the configuration for an RDS instance
// provided by AWS that is used by a workload.
type AwsRelationalDatabaseDefinition struct {
	Common     `swaggerignore:"true" mapstructure:",squash"`
	Definition `mapstructure:",squash"`

	// The database engine for the instance.  One of:
	// * mysql
	// * postgres
	// * mariadb
	Engine *string `json:"Engine,omitempty" query:"engine" gorm:"not null" validate:"required"`

	// The version of the database engine for the instance.
	EngineVersion *string `json:"EngineVersion,omitempty" query:"engineversion" gorm:"not null" validate:"required"`

	// The name of the database that will be used by the client workload.
	DatabaseName *string `json:"DatabaseName,omitempty" query:"databasename" gorm:"not null" validate:"required"`

	// The port to use to connect to the database.
	DatabasePort *int `json:"DatabasePort,omitempty" query:"databaseport" gorm:"not null" validate:"required"`

	// The number of days to retain database backups for.
	BackupDays *int `json:"BackupDays,omitempty" query:"BackupDays" gorm:"default: 0" validate:"optional"`

	// The amount of compute capacity to use for the database virtual machine.
	MachineSize *string `json:"MachineSize,omitempty" query:"machinesize" gorm:"not null" validate:"required"`

	// The amount of storage in Gb to allocate for the database.
	StorageGb *int `json:"StorageGb,omitempty" query:"storagegb" gorm:"not null" validate:"required"`

	// The name of the Kubernetes secret that will be attached to the
	// running workload from which database connection configuration will be
	// supplied.  This secret name must be referred to in the Kubernetes
	// manifest, .e.g Deployment, for the workload.
	WorkloadSecretName *string `json:"WorkloadSecretName,omitempty" query:"WorkloadSecretName" gorm:"not null" validate:"required"`

	// The AWS account in which the RDS instance will be provisioned.
	AwsAccountID *uint `json:"AwsAccountID,omitempty" query:"awsaccountid" gorm:"not null" validate:"required"`
}

// +threeport-codegen:reconciler
// AwsRelationalDatabaseInstance is a deployed instance of an RDS instance.
type AwsRelationalDatabaseInstance struct {
	Common         `swaggerignore:"true" mapstructure:",squash"`
	Instance       `mapstructure:",squash"`
	Reconciliation `mapstructure:",squash"`

	// The definition that configures this instance.
	AwsRelationalDatabaseDefinitionID *uint `json:"AwsRelationalDatabaseDefinitionID,omitempty" query:"awsrelationaldatabasedefinitionid" gorm:"not null" validate:"required"`

	// An inventory of all AWS resources for the EKS cluster.
	ResourceInventory *datatypes.JSON `json:"ResourceInventory,omitempty" validate:"optional"`

	// The ID of the workload instance that the database instance serves.
	WorkloadInstanceID *uint `json:"WorkloadInstanceID,omitempty" query:"workloadinstanceid" gorm:"not null" validate:"required"`
}

// AwsObjectStorageBucketDefinition is the configuration for an S3 bucket
// provided by AWS that is used for object storage by a workload.
type AwsObjectStorageBucketDefinition struct {
	Common     `swaggerignore:"true" mapstructure:",squash"`
	Definition `mapstructure:",squash"`

	// When true, objects in the bucket are publicly readable by anyone - for use
	// cases such as storing static assets for public websites.  When false,
	// only the workload attached to an AWSObjectStorageBucketInstance and the AWS users
	// on the account may access the bucket for read or write.
	PublicReadAccess *bool `json:"PublicReadAccess,omitempty" query:"publicreadaccess" gorm:"default:false" validate:"optional"`

	// The name of the Kubernetes service account for the workload that will
	// access the S3 bucket.  Used to provide secure access using IAM roles for
	// service accounts (IRSA).
	WorkloadServiceAccountName *string `json:"WorkloadServiceAccountName,omitempty" query:"workloadserviceaccountname" gorm:"not null" validate:"required"`

	// The environment variable key that the workload is expecting to reference
	// for the name of the S3 bucket managed by threeport.
	WorkloadBucketEnvVar *string `json:"WorkloadBucketEnvVar,omitempty" query:"workloadbucketenvvar" gorm:"not null" validate:"required"`

	// The AWS account in which the RDS instance will be provisioned.
	AwsAccountID *uint `json:"AwsAccountID,omitempty" query:"awsaccountid" gorm:"not null" validate:"required"`
}

// +threeport-codegen:reconciler
// AwsObjectStorageBucketInstance is a deployed instance of an S3 bucket.
type AwsObjectStorageBucketInstance struct {
	Common         `swaggerignore:"true" mapstructure:",squash"`
	Instance       `mapstructure:",squash"`
	Reconciliation `mapstructure:",squash"`

	// An inventory of all AWS resources for the S3 bucket.
	ResourceInventory *datatypes.JSON `json:"ResourceInventory,omitempty" validate:"optional"`

	// The definition that configures this instance.
	AwsObjectStorageBucketDefinitionID *uint `json:"AwsObjectStorageBucketDefinitionID,omitempty" query:"awsobjectstoragebucketdefinitionid" gorm:"not null" validate:"required"`

	// The ID of the workload instance that uses the S3 bucket.
	WorkloadInstanceID *uint `json:"WorkloadInstanceID,omitempty" query:"workloadinstanceid" gorm:"not null" validate:"required"`
}

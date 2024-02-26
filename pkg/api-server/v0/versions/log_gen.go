// generated by 'threeport-sdk gen api-model' - do not edit

package versions

import (
	api "github.com/threeport/threeport/pkg/api"
	iapi "github.com/threeport/threeport/pkg/api-server/v0"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	"reflect"
)

// AddLogBackendVersions adds field validation info and adds it
// to the REST API versions.
func AddLogBackendVersions() {
	iapi.LogBackendTaggedFields[iapi.TagNameValidate] = &iapi.FieldsByTag{
		Optional:             []string{},
		OptionalAssociations: []string{},
		Required:             []string{},
		TagName:              iapi.TagNameValidate,
	}

	// parse struct and populate the FieldsByTag object
	iapi.ParseStruct(
		iapi.TagNameValidate,
		reflect.ValueOf(new(v0.LogBackend)),
		"",
		iapi.Translate,
		iapi.LogBackendTaggedFields,
	)

	// create a version object which contains the object name and versions
	versionObj := iapi.VersionObject{
		Object:  string(v0.ObjectTypeLogBackend),
		Version: iapi.V0,
	}

	// add the object tagged fields to the global tagged fields map
	iapi.ObjectTaggedFields[versionObj] = iapi.LogBackendTaggedFields[iapi.TagNameValidate]

	// add the object tagged fields to the rest API version
	api.AddRestApiVersion(versionObj)
}

// AddLogStorageDefinitionVersions adds field validation info and adds it
// to the REST API versions.
func AddLogStorageDefinitionVersions() {
	iapi.LogStorageDefinitionTaggedFields[iapi.TagNameValidate] = &iapi.FieldsByTag{
		Optional:             []string{},
		OptionalAssociations: []string{},
		Required:             []string{},
		TagName:              iapi.TagNameValidate,
	}

	// parse struct and populate the FieldsByTag object
	iapi.ParseStruct(
		iapi.TagNameValidate,
		reflect.ValueOf(new(v0.LogStorageDefinition)),
		"",
		iapi.Translate,
		iapi.LogStorageDefinitionTaggedFields,
	)

	// create a version object which contains the object name and versions
	versionObj := iapi.VersionObject{
		Object:  string(v0.ObjectTypeLogStorageDefinition),
		Version: iapi.V0,
	}

	// add the object tagged fields to the global tagged fields map
	iapi.ObjectTaggedFields[versionObj] = iapi.LogStorageDefinitionTaggedFields[iapi.TagNameValidate]

	// add the object tagged fields to the rest API version
	api.AddRestApiVersion(versionObj)
}

// AddLogStorageInstanceVersions adds field validation info and adds it
// to the REST API versions.
func AddLogStorageInstanceVersions() {
	iapi.LogStorageInstanceTaggedFields[iapi.TagNameValidate] = &iapi.FieldsByTag{
		Optional:             []string{},
		OptionalAssociations: []string{},
		Required:             []string{},
		TagName:              iapi.TagNameValidate,
	}

	// parse struct and populate the FieldsByTag object
	iapi.ParseStruct(
		iapi.TagNameValidate,
		reflect.ValueOf(new(v0.LogStorageInstance)),
		"",
		iapi.Translate,
		iapi.LogStorageInstanceTaggedFields,
	)

	// create a version object which contains the object name and versions
	versionObj := iapi.VersionObject{
		Object:  string(v0.ObjectTypeLogStorageInstance),
		Version: iapi.V0,
	}

	// add the object tagged fields to the global tagged fields map
	iapi.ObjectTaggedFields[versionObj] = iapi.LogStorageInstanceTaggedFields[iapi.TagNameValidate]

	// add the object tagged fields to the rest API version
	api.AddRestApiVersion(versionObj)
}

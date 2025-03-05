package testtypes

type TestExportedWithUnexportedField struct {
	ExportedField string `default:"foo"`
	//nolint:unused // we want to see if the defaulting process tries to set this field
	unexportedField string `default:"bar"`
}

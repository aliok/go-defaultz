package testtypes

type TestExportedWithUnexportedField struct {
	ExportedField string `default:"foo"`
	//nolint:unused
	unexportedField string `default:"bar"`
}

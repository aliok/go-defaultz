package testtypes

type TestExportedWithUnexportedField struct {
	ExportedField   string `default:"foo"`
	unexportedField string `default:"bar"`
}

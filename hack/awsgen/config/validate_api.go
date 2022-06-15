package config

func validateAPICall(call string) []error {
	return validateExportedIdentifier("call", call)
}

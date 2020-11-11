package schema8

func init() {
	gAssertionConstructorMap = map[string]assertionConstructorFunc{
		"type":  newAssertionType,
		"enum":  newAssertionEnum,
		"const": newAssertionConst,

		"multipleOf":       newAssertionMultipleOf,
		"maximum":          newAssertionMaximum,
		"minimum":          newAssertionMinimum,
		"exclusiveMaximum": newAssertionExMaximum,
		"exclusiveMinimum": newAssertionExMinimum,

		"maxLength": newAssertionMaxLength,
		"minLength": newAssertionMinLength,
		"pattern":   newAssertionPattern,

		"maxItems":    newAssertionMaxItems,
		"minItems":    newAssertionMinItems,
		"maxContains": newAssertionMaxContains,
		"minContains": newAssertionMinContains,
		"uniqueItems": newAssertionUniqueItems,

		"maxProperties":     newAssertionMaxProperties,
		"minProperties":     newAssertionMinProperties,
		"required":          newAssertionRequired,
		"dependentRequired": newAssertionDependentRequired,

		"format": newAssertionFormat,
	}

	// applicator 'items' will be handled specially.

	gApplicatorArraySchemaMap = map[string]applicatorConstructorFunc{
		"allOf": newApplicatorAllOf,
		"anyOf": newApplicatorAnyOf,
		"oneOf": newApplicatorOneOf,
	}

	gApplicatorObjectSchemaMap = map[string]applicatorConstructorFunc{
		"dependentSchemas":  newApplicatorDependentSchemas,
		"properties":        newApplicatorProperties,
		"patternProperties": newApplicatorPatternProperties,
	}

	gApplicatorSchemaMap = map[string]applicatorConstructorFunc{
		"not": newApplicatorNot,

		"if":   newApplicatorIf,
		"then": newApplicatorThen,
		"else": newApplicatorElse,

		"additionalItems":  newApplicatorAdditionalItems,
		"unevaluatedItems": newApplicatorUnevaluatedItems,
		"contains":         newApplicatorContains,

		"additionalProperties":  newApplicatorAdditionalProperties,
		"unevaluatedProperties": newApplicatorUnevaluatedProperties,
		"propertyNames":         newApplicatorPropertyNames,
	}

	gFormatFuncMap = map[string]FormatFunc{
		"date-time":             nil,
		"date":                  nil,
		"time":                  nil,
		"duration":              nil,
		"email":                 formatEmail, // foo@protonmail.com
		"idn-email":             nil,
		"hostname":              nil,
		"idn-hostname":          nil,
		"ipv4":                  formatIPV4,   // 10.0.0.1
		"ipv6":                  formatIPV6,   // 2001:cdba:0000:0000:0000:0000:3257:9652
		"uri":                   formatURI,    // absolute URI, e.g.: https://www.duckduckgo.com
		"uri-reference":         formatURIRef, // absolute URI or uri-reference
		"iri":                   nil,
		"iri-reference":         nil,
		"uuid":                  formatUUID,
		"json-pointer":          nil,
		"relative-json-pointer": nil,
		"regex":                 formatRegex,

		// Custom
		"mysql-datetime":   formatMySQLDateTime,   // 2006-01-02 15:04:05
		"mongodb-datetime": formatMongoDBDateTime, // 2018-11-16T06:16:36.156Z
		"timestamp":        formatTimestamp,       // 1604906268, now +/- 10 min is valid
		"timestamp-ms":     formatTimestampMS,     // 1604906268123, now +/- 10 min is valid
	}
}

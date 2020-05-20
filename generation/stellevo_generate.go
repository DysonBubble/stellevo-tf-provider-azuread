package main
import (
	"os"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
  "fmt"
)

func main() {
	args := os.Args[1:]
	provider := Provider().(*schema.Provider)
	DoMain(args[0], args[1], provider.ResourcesMap, provider.DataSourcesMap)
}

func ProviderSpecificInputPrelude() string {
	return `export const ` + TS_SCHEMA_TARGET_REF + ` = common.SchemaUnion( [
	common.SchemaString,
	common.SchemaRecord( {
		resource: common.SchemaKeyOf(` + TS_PROVIDER_RESOURCES_IMPORT_NAME + `.` + TS_ALL_RESOURCES_VAR_NAME + `),
		id: common.SchemaString
	} )
	], "` + TS_SCHEMA_TARGET_REF + `" );
export type ` + TS_TYPE_TARGET_REF + ` = common.TypeOf<typeof ` + TS_SCHEMA_TARGET_REF + `>;`
}

func ProviderSpecificSensitiveStringSchemaDefinition(schemaNameVariableName string) string {
	return `common.SchemaUnion( [
  common.SchemaRecord( {
    ss: common.SchemaLiteral( "managed" ),
    id: common.SchemaString
  } ),
  common.SchemaRecord( {
    ss: common.SchemaLiteral( "imported" ),
    vault_id: common.SchemaString,
    secret_name: common.SchemaString
  } )
], ` + schemaNameVariableName + " )"
}

var seenIDNames = map[string]string{
	"key_id": "common.SchemaString",
	"immutable_id": "common.SchemaString",
}
func TryGetProviderSpecificStringSchema(isInput bool, resource_id string, path []string) string {
	retVal := ""
	if (isInput) {
		property_name := path[len(path) - 1]
		has_id := strings.HasSuffix(property_name, "_id")
		has_ids := strings.HasSuffix(property_name, "_ids")
		if ( has_id || has_ids ) {
			if existingValue, ok := seenIDNames[property_name]; ok {
				retVal = existingValue
			} else {
				retVal = TS_SCHEMA_TARGET_REF
				if has_ids {
					retVal = fmt.Sprintf("common.SchemaMap(%s)", retVal)
				}
				seenIDNames[property_name] = retVal
			}
		}
	}
	return retVal
}

func GetProviderSpecificUnionSchemaCodeGenHandling() string {
	return `case "` + TS_SCHEMA_TARGET_REF + `":
			const targetRef = ` + TS_CODEGEN_CALLBACK_PARAMETER_VALUE + ` as ` + TS_SCHEMA_INPUT_IMPORT_NAME + `.` + TS_TYPE_TARGET_REF + `;
      return typeof targetRef === 'string' ? { kind: "` + TS_CODEGEN_RETVAL_KIND_INDIRECT + `", value: ` + TS_CODEGEN_CALLBACK_PARAMETER_VALUE + `, schema: ` + TS_CODEGEN_CALLBACK_SCHEMA_VALUE + `.types[0] as ` + TS_CODEGEN_COMMON_IMPORT_NAME + `.TPropertySchemas } : { kind: "` + TS_CODEGEN_RETVAL_KIND_DIRECT + `", value: ` + "`${ targetRef.resource }.${ targetRef.id }.id` };";
}

const TS_SCHEMA_TARGET_REF = "SchemaTargetReference"
const TS_TYPE_TARGET_REF = "TTargetReference"

func GetProvidersThisDependsOn() []string {
	return make([]string, 0)
}
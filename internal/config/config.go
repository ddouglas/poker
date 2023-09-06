package config

import (
	"context"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type pathConfig struct {
	name     string
	required bool
	provider string
	value    reflect.Value
	// value    *string
}

type LoadOpts struct {
	prefix string
	client *ssm.Client
}

type LoadOptFunc func(o *LoadOpts)

func WithSSMClient(c *ssm.Client) LoadOptFunc {
	return func(o *LoadOpts) {
		o.client = c
	}
}

func WithPrefix(prefix string) LoadOptFunc {
	return func(o *LoadOpts) {
		o.prefix = prefix
	}
}

func Load(ctx context.Context, out any, optFuncs ...LoadOptFunc) error {

	opts := new(LoadOpts)
	for _, f := range optFuncs {
		f(opts)
	}

	if opts.prefix == "" {
		opts.prefix = "/"
	}

	outType := reflect.TypeOf(out)
	outValue := reflect.ValueOf(out)

	if outType.Kind() != reflect.Ptr || outType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Load must be passed a pointer to a struct")
	}

	outValue = outValue.Elem()

	if opts.client == nil {
		config, err := awsConfig.LoadDefaultConfig(ctx)
		if err != nil {
			return fmt.Errorf("failed to load default config for aws: %w", err)
		}
		opts.client = ssm.NewFromConfig(config)
	}

	names := getSSMRecursiveTags(outValue, opts.prefix)
	names = append(names, getEnvRecursiveTags(outValue)...)

	envPathName := make([]string, 0, len(names))
	ssmPathNames := make([]string, 0, len(names))
	for _, nc := range names {
		switch nc.provider {
		case "env":
			envPathName = append(envPathName, nc.name)
		case "ssm":
			ssmPathNames = append(ssmPathNames, nc.name)
		default:
			return fmt.Errorf("unhandled provider %s: %s", nc.provider, nc.name)
		}
	}
	resultMap := make(map[string]string)

	if len(envPathName) > 0 {
		// os.Environ returns a slice of values
		// each value is the Name=Values
		// to get the names, we'll build a map of name to value
		envList := os.Environ()

		for _, e := range envList {
			ee := strings.Split(e, "=")
			name, value := ee[0], ee[1]
			resultMap[name] = value
		}

	}

	if len(ssmPathNames) > 0 {
		result, err := opts.client.GetParameters(ctx, &ssm.GetParametersInput{
			Names:          ssmPathNames,
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("failed to fetch parameters: %w", err)
		}

		for _, parameter := range result.Parameters {
			resultMap[aws.ToString(parameter.Name)] = aws.ToString(parameter.Value)
		}

	}

	for _, p := range names {

		result, ok := resultMap[p.name]
		if !ok {
			fmt.Printf("%s is missing, required: %t\n", p.name, p.required)
			if p.required {

				return fmt.Errorf("%s is required but no values was located with provider %s", p.name, p.provider)
			}
			continue
		}

		setFieldValue(p.value, result)
	}

	return nil
	// return setRecursiveTags(outValue, opts.prefix, resultMap)

}

func setFieldValue(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		boolValue, _ := strconv.ParseBool(value)
		field.SetBool(boolValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, _ := strconv.ParseInt(value, 10, 64)
		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, _ := strconv.ParseUint(value, 10, 64)
		field.SetUint(uintValue)
	case reflect.Float32, reflect.Float64:
		floatValue, _ := strconv.ParseFloat(value, 64)
		field.SetFloat(floatValue)
	}
}

func getSSMRecursiveTags(v reflect.Value, prefix string) []*pathConfig {

	if prefix == "" {
		prefix = "/"
	}

	t := v.Type()
	nameConfigs := make([]*pathConfig, 0)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldT := t.Field(i)
		tag := fieldT.Tag.Get("ssm")
		if fieldT.Type.Kind() == reflect.Struct {
			t := tag
			if tag != "" {
				t = path.Join(prefix, t)
			}
			nameConfigs = append(nameConfigs, getSSMRecursiveTags(field, t)...)
			continue
		}

		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		required := false
		if len(parts) == 2 && parts[1] == "required" {
			required = true
		}

		nameConfigs = append(nameConfigs, &pathConfig{
			name:     path.Join(prefix, parts[0]),
			required: required,
			provider: "ssm",
			value:    field,
		})
	}

	return nameConfigs
}

func getEnvRecursiveTags(v reflect.Value) []*pathConfig {

	t := v.Type()
	nameConfigs := make([]*pathConfig, 0)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldT := t.Field(i)
		tag := fieldT.Tag.Get("env")
		if fieldT.Type.Kind() == reflect.Struct {
			nameConfigs = append(nameConfigs, getEnvRecursiveTags(field)...)
			continue
		}

		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		required := false
		if len(parts) == 2 && parts[1] == "required" {
			required = true
		}

		nameConfigs = append(nameConfigs, &pathConfig{
			name:     parts[0],
			required: required,
			provider: "env",
			value:    field,
		})

	}

	return nameConfigs
}

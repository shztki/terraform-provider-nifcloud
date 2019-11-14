package nifcloud

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func validateDbName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if (len(value) < 1) || (len(value) > 64) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length", k))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and underscore allowed in %q", k))
	}
	return
}

func validateUserName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if (len(value) < 1) || (len(value) > 16) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 16 characters in length", k))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and underscore allowed in %q", k))
	}
	return
}

func validatePassword(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if (len(value) < 1) || (len(value) > 41) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 41 characters in length", k))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters allowed in %q", k))
	}
	return
}

func validateRdbIdentifier(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if (len(value) < 1) || (len(value) > 63) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 63 characters in length", k))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and hyphens allowed in %q", k))
	}
	if !regexp.MustCompile(`^[a-zA-Z]`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"first character of %q must be a letter", k))
	}
	if regexp.MustCompile(`--`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot contain two consecutive hyphens", k))
	}
	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a hyphen", k))
	}
	return
}

func validateOnceADayWindowFormat(v interface{}, k string) (ws []string, errors []error) {
	// valid time format is "hh24:mi"
	validTimeFormat := "([0-1][0-9]|2[0-3]):([0-5][0-9])"
	validTimeFormatConsolidated := "^(" + validTimeFormat + "-" + validTimeFormat + "|)$"

	value := v.(string)
	if !regexp.MustCompile(validTimeFormatConsolidated).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must satisfy the format of \"hh24:mi-hh24:mi\"", k))
	}
	return
}

func validateOnceAWeekWindowFormat(v interface{}, k string) (ws []string, errors []error) {
	// valid time format is "ddd:hh24:mi"
	validTimeFormat := "(sun|mon|tue|wed|thu|fri|sat):([0-1][0-9]|2[0-3]):([0-5][0-9])"
	validTimeFormatConsolidated := "^(" + validTimeFormat + "-" + validTimeFormat + "|)$"

	value := strings.ToLower(v.(string))
	if !regexp.MustCompile(validTimeFormatConsolidated).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must satisfy the format of \"ddd:hh24:mi-ddd:hh24:mi\"", k))
	}
	return
}

func validateDbParamGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if (len(value) < 1) || (len(value) > 255) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 255 characters in length", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z]`).MatchString(value) {
		errors = append(errors, fmt.Errorf("first character of %q must be a letter", k))
	}

	if regexp.MustCompile(`--`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q cannot contain two consecutive hyphens", k))
	}

	if strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("%q cannot end with a - character", k))
	}

	if !regexp.MustCompile(`^[0-9a-zA-Z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q can only contain alphanumeric and %q characters", k, "-"))
	}

	return
}

func validateLbName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if (len(value) < 1) || (len(value) > 15) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 15 characters in length", k))
	}
	if !regexp.MustCompile(`^[0-9A-Za-z]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alphanumeric characters allowed in %q: %q",
			k, value))
	}
	return
}

func validateListenerProtocol() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"HTTP",
		"HTTPS",
		"FTP",
		"",
	}, true)
}

func validateVpnConnectionTunnelPreSharedKey(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if (len(value) < 1) || (len(value) > 64) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length", k))
	}

	if strings.HasPrefix(value, "0") {
		errors = append(errors, fmt.Errorf("%q cannot start with zero character", k))
	}

	if !regexp.MustCompile(`^[0-9a-zA-Z-+&!@#$%^*(),.:_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q can only contain alphanumeric and %q characters", k, "-+&!@#$%^*(),.:_"))
	}

	return
}

func validateHeathCheckTarget(v interface{}, k string) (ws []string, errors []error) {
	// TCP:port | ICMP
	validFormat := "(ICMP|TCP:[0-9]{1,5})"
	validFormatConsolidated := "^(" + validFormat + "|)$"

	value := strings.ToUpper(v.(string))
	if !regexp.MustCompile(validFormatConsolidated).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must satisfy the format of \"TCP:port | ICMP\"", k))
	}
	return
}

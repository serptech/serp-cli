package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/serptech/serp-go/api/const/conf"
	"github.com/serptech/serp-go/api/const/liveness"
	"github.com/serptech/serp-go/utils"
)

func writeOutput(data interface{}) {
	if outputPath != "" {
		preOut, err := utils.GetPretty(data)
		ifErrorExit(err)
		ifErrorExit(utils.WriteToFile(outputPath, preOut))
		return
	}
	ifErrorExit(utils.PrettyPrint(data))
}

func stringPtr(v string) *string { return &v }

func boolPtr(v bool) *bool { return &v }

func intPtr(v int) *int { return &v }

func resolveConf(value string) (conf.Conf, error) {
	lower := strings.ToLower(strings.TrimSpace(value))
	switch lower {
	case "nm", "no-match", "no-matches":
		return conf.Nm, nil
	case "new":
		return conf.New, nil
	case "exact":
		return conf.Exact, nil
	case "junk":
		return conf.Junk, nil
	case "ha", "high-accuracy":
		return conf.Ha, nil
	case "det":
		return conf.Det, nil
	case "reinit":
		return conf.Reinit, nil
	case "nf", "no-face":
		return conf.Nf, nil
	case "":
		return 0, fmt.Errorf("conf value is required")
	default:
		val, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("unknown conf value %q", value)
		}
		parsed := conf.Conf(val)
		if err := parsed.Validate(); err != nil {
			return 0, err
		}
		return parsed, nil
	}
}

func resolveLiveness(value string) (liveness.Liveness, error) {
	lower := strings.ToLower(strings.TrimSpace(value))
	switch lower {
	case "passed":
		return liveness.Passed, nil
	case "failed":
		return liveness.Failed, nil
	case "undetermined":
		return liveness.Undetermined, nil
	case "":
		return "", fmt.Errorf("liveness value is required")
	default:
		return "", fmt.Errorf("unknown liveness %q", value)
	}
}

func parseDate(value string) (time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, nil
	}
	layouts := []string{time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, trimmed); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date %q", value)
}

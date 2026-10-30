package actioninfo

import (
	"fmt"
	"log/slog"
)

type DataParser interface {
	Parse(datastring string) error
	ActionInfo() (string, error)
}

func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		if err := dp.Parse(data); err != nil {
			slog.Error("Failed to parse data", "data", data, "error", err)
			continue
		}

		info, err := dp.ActionInfo()
		if err != nil {
			slog.Error("Failed to generate info", "data", data, "error", err)
			continue
		}

		fmt.Println(info)
	}
}

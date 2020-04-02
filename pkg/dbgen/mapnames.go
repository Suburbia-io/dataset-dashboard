package dbgen

import "strings"

func SQLNameToGo(name string) string {
	parts := strings.Split(name, "_")
	for i, part := range parts {
		switch part {
		case "id", "api", "sftp":
			parts[i] = strings.ToUpper(part)
		default:
			parts[i] = strings.ToUpper(string(part[0])) + part[1:]
		}
	}

	return strings.Join(parts, "")
}

func SQLNameToJSON(name string) string {
	parts := strings.Split(name, "_")
	for i, part := range parts {
		if i > 0 {
			switch part {
			case "id", "api", "sftp":
				parts[i] = strings.ToUpper(part)
			default:
				parts[i] = strings.ToUpper(string(part[0])) + part[1:]
			}
		}
	}
	return strings.Join(parts, "")
}

func TableNameToRowType(name string) string {
	name = SQLNameToGo(name)
	if name[len(name)-1] == 's' {
		name = name[:len(name)-1]
	}
	return name
}

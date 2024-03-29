package object

import "fmt"

func FormatObject(object Object) string {
	return fmt.Sprintf(
		"Key: %s\nLast Modified: %s\nSize: %d KB\n %s\nContent:\n\n%s",
		object.Key,
		object.LastModified.Format("2006-01-02"),
		object.Size,
		"---",
		object.Content,
	)
}

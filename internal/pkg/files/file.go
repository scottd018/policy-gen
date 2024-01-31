package files

import "fmt"

func Path(dir *Directory, file string, extension string) string {
	return fmt.Sprintf("%s/%s.%s", dir.Path, file, extension)
}

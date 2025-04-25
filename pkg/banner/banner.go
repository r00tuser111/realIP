package banner

import (
	"fmt"
)

// ASCII艺术字体 - realIP
const asciiArt = `
 ____            _  ___ ____  
|  _ \ ___  __ _| |/_ /|  _ \ 
| |_) / _ \/ _` + "`" + `| | | || |_) |
|  _ <  __/ (_| | | | ||  __/ 
|_| \_\___|\__,_|_||___|_|    

Domain to IP Resolution Tool
`

// 版本信息
const (
	Version = "1.0.0"
	Author  = "kking"
	GitHub  = "github.com/r00tuser111/realip"
)

// Show 显示banner
func Show() {
	fmt.Println(asciiArt)
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Author: %s\n", Author)
	fmt.Printf("GitHub: %s\n", GitHub)
	fmt.Println()
} 
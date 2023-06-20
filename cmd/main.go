package main

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

var (
	MainRoot string
)

func main() {
	args := os.Args
	var dirName []string

	tags := ""
	if len(args) < 2 || len(args) > 3 {
		fmt.Println("invalid number of arguments. \n <fisherman -h> for help")
		return
	} else if len(args) == 3 {
		tags = args[2]
	}

	path := strings.Split(args[1], "/")

	MainRoot = path[len(path)-1]

	dirList, err := os.ReadDir(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(MainRoot)

	for _, dir := range dirList {
		switch tags {
		case "-h":
			fmt.Println("Fisherman is a tool to visualize the directories and files in a neat tree")
			fmt.Println()
			fmt.Println("usage:: fisherman <tags>")
			fmt.Println("._______________________________________.")
			fmt.Println("|tags          |             description|")
			fmt.Println("|-h            |                    help|")
			fmt.Println("|-l            |                list all|")
			fmt.Println("|-e            |      extend directories|")
			fmt.Println("|-a            |  extend all directories|")
			fmt.Println("'.-------------------------------------.'")
			fmt.Println("By 'all directories', hidden directories also")
			return

		case "-l":
			PrintDirs(dir.Name(), dirName, false)

		case "-a":
			if dir.IsDir() {
				nestedDirList, err := os.ReadDir(args[1] + "/" + dir.Name())
				if err != nil {
					return
				}
				for _, nDir := range nestedDirList {
					dirName = append(dirName, nDir.Name())

				}
				PrintDirs(dir.Name(), dirName, false)
				dirName = dirName[:0]
			} else {
				PrintDirs(dir.Name(), dirName, false)
			}
		case "-e":
			if dir.Name()[0] != '.' {
				if dir.IsDir() {
					nestedDirList, err := os.ReadDir(args[1] + "/" + dir.Name())
					if err != nil {
						return
					}
					for _, nDir := range nestedDirList {
						if nDir.Name()[0] != '.' {
							dirName = append(dirName, nDir.Name())
						}
					}
					lastDir := false
					if dir.Name() == dirList[len(dirList)-1].Name() {
						lastDir = true
					}
					PrintDirs(dir.Name(), dirName, lastDir)
					dirName = dirName[:0]
				} else {
					if dir.Name()[0] != '.' {
						PrintDirs(dir.Name(), dirName, false)
					}
				}
			}
		case "":
			if dir.Name()[0] != '.' {
				PrintDirs(dir.Name(), dirName, false)
			}
		case "-p":
			PrintDirStats(args[1], dirList)
			return
		default:
			fmt.Println("invalid tag")
			fmt.Println("<fisherman -h> for help")
			return
		}

	}

}

func PrintDirStats(root string, dirList []fs.DirEntry) {
	fmt.Println("file information of", MainRoot)
	fmt.Println("name                    |size                    |mode                   |isDir                    |")
	fmt.Println("|-----------------------|------------------------|-----------------------|-------------------------|")

	for _, dir := range dirList {
		stats, err := os.Stat(root + "/" + dir.Name())
		if err != nil {
			fmt.Println(err)
			return
		}
		length := 0
		name := stats.Name()

		if len(name) < 24-1 {
			length = len(name) + 1
			for i := 0; i < 24-length; i++ {
				name += " "
			}
		}

		size := strconv.Itoa(int(stats.Size()))
		if len(size) < 24 {
			length = len(size) + 2
			for i := 0; i < 24-length; i++ {
				size += " "
			}
		}

		mode := stats.Mode().String()
		if len(mode) < 24 {
			length = len(mode) + 3
			for i := 0; i < 24-length; i++ {
				mode += " "
			}
		}

		yorno := "no"
		isDir := stats.IsDir()
		if isDir {
			yorno = "yes"
		}

		if len(yorno) < 24 {
			length = len(yorno) + 1
			for i := 0; i < 24-length; i++ {
				yorno += " "
			}
		}
		fmt.Println(name, "|", size, "|", mode, "|", yorno, "|")
		fmt.Println("'-----------------------|------------------------'-----------------------|-------------------------|")
	}
}

func PrintDirs(root string, dirName []string, lastDir bool) {
	PrintSpace(len(MainRoot), false)
	fmt.Print("--")
	for len(root) < 16 {
		root = root + "_"
	}

	fmt.Println(root)
	for _, name := range dirName {
		PrettyPrintDir(len(MainRoot), name, lastDir)
	}
}

func PrintSpace(length int, lastDir bool) {
	for i := 0; i <= length; i++ {
		if i == len(MainRoot) && !lastDir {
			fmt.Print("|")
		} else if i == len(MainRoot) && lastDir {
			fmt.Print(" ")
		}
		fmt.Print(" ")
	}
}

func PrettyPrintDir(rootLen int, name string, lastDir bool) {
	PrintSpace(rootLen+18, lastDir)
	fmt.Print("|")
	fmt.Print("__ ")
	fmt.Print(name)
	fmt.Println()
}

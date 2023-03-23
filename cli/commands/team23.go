/*
Root command is created as team23. All other commands are built on
top of this command. Creation of new commands requires an init
function per command with rootCmd.AddCommand(<newCmd>)
*/

package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/19chonm/461_1_23/fileio"
	"github.com/19chonm/461_1_23/logger"
	"github.com/19chonm/461_1_23/worker"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "team23",
	Short: "team23 - root command for app",
	Long:  "team23 is the root command to navigate through Team 23's CLI",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// test URL_FILE with: /Users/emile/461_1_23/test/urls_file.txt

// First function to be ran on main. Will check if second argument is either
// an absolute filepath, one of the recognized commands or neither. If neither,
// program will throw error. If argument is an absolute filepath, a direct call
// to functions are executed. No cobra command is created because name varies.
func Execute() {

	if len(os.Args) != 2 {
		logger.DebugMsg(`CLI: Please use one of the recognized commands: 'build', 
		'install', 'test', or 'URL_FILE' where URL_FILE is an absolute path 
		to a file`)
	} else if filepath.IsAbs(os.Args[1]) {
		//read the urls into an array
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("error opening url file: %v", err)
			return
		}
		defer file.Close()
		var lines []string
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := strings.TrimRight(scanner.Text(), "\n")
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("error reading URL file: %v", err)
			return
		}
		//TODO
		//store and sort the packages in the for loop above and print the results
		var packageRatings []*fileio.Rating
		// Start output
		for _, url := range lines {
			r := worker.RunTask(url)
			if r == nil {
				fmt.Printf("error running tasks for the url: %v", err)
				continue
			}
			packageRatings = append(packageRatings, r)
		}

		//TODO: sort and output the results
		packageRatings = SortPackages(packageRatings)
		for _, obj := range packageRatings {
			b, err := json.Marshal(obj)
			if err != nil {
				fmt.Println("error", err)
				return
			}
			fmt.Println(string(b))
		}

		// ratings, errors := fileio.ReadWorkerResults(worker_output_ch)
		// if len(errors) > 0 {
		// 	fileio.PrintErrors(errors)
		// 	os.Exit(1)
		// } else {
		// 	fileio.Print_sorted_output(ratings)
		// }

	} else if os.Args[1] == "build" || os.Args[1] == "install" ||
		os.Args[1] == "test" {

		if err := rootCmd.Execute(); err != nil {
			logger.DebugMsg("CLI: Error using CLI ", err.Error())
			os.Exit(1)
		}
	} else {
		logger.DebugMsg("CLI: Not a recognized command")
		os.Exit(1)
	}

	os.Exit(0)
}

func SortPackages(packages []*fileio.Rating) []*fileio.Rating {
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].NetScore > packages[j].NetScore
	})

	return packages
}

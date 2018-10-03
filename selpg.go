package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

//  selpg argument struct
type selpg_args struct {
	start_page  uint
	end_page    uint
	in_filename string
	page_len    uint
	page_type   uint
	print_dest  string
}

var sp_args = selpg_args{0, 0, "", 72, 'l', ""}

// flag
var (
	h           bool
	s           uint
	e           uint
	l           uint
	f           bool
	d           string
	inputStream = os.Stdin     //default
	stdout      io.WriteCloser //default
	err         error
)

// programme name
var progname string = "selpg"

// init the flag
func Init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.UintVar(&s, "s", 0, "start page")
	flag.UintVar(&e, "e", 0, "end page")
	flag.UintVar(&l, "l", 0, "the number of line in one page")
	flag.BoolVar(&f, "f", false, "the mode that a page depends on the \f")
	flag.StringVar(&d, "d", "", "the destination")
	flag.Usage = usage
}

// selpg command usage
func usage() {
	fmt.Fprintf(os.Stderr, "\nUSAGE: %s -s start_page -e end_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n",
		progname)

}

// print the help list
func helpList() {
	fmt.Println("Options: ")
	flag.PrintDefaults()
}

// main function
func main() {
	Init()
	//process the input arguments
	process_args()
	//accoring to arguments to execute the command
	process_input()

}

//process the arguments
func process_args() {
	//parse the flag
	flag.Parse()

	//if -h then print the help list
	if h {
		flag.Usage()
		helpList()
		os.Exit(1)
	}
	// judge if there are enough arguments
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
		os.Exit(6)
	}
	//-s format check
	if os.Args[1] != "-s" {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -s start_page\n", progname)
		flag.Usage()
		os.Exit(2)
	} else {
		if s <= 0 {
			fmt.Fprintf(os.Stderr, "%s: start_page must be greater than zero\n", progname)
			os.Exit(3)
		}
		sp_args.start_page = s

	}
	//-e format check
	if os.Args[3] != "-e" {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -e end_page\n", progname)
		flag.Usage()
		os.Exit(4)
	} else {
		if e <= 0 {
			fmt.Fprintf(os.Stderr, "%s: end_page must be greater than zero\n", progname)
			os.Exit(5)
		}
		if e < s {
			fmt.Fprintf(os.Stderr, "%s: end_page must be greater or equal to start_page\n", progname)
			os.Exit(7)
		}

		sp_args.end_page = e

	}
	// -s start_page -e end_page [-f | -l page_line][-d destination][ filename ]
	// f mode check
	if f {
		//-f and -l can not exit at the same time
		if l != 0 {
			fmt.Fprintf(os.Stderr, "%s: -f and -l can not exist at the same time\n", progname)
			os.Exit(8)
		}
		//change the page_type
		sp_args.page_type = 'f'
	} else {
		//the dafault number of line in one page is 72 or can be set
		if l == 0 {
			l = 72
		}
		sp_args.page_len = l
	}
	// print destination
	if d != "" {
		sp_args.print_dest = d
	}

	//get input file name
	if flag.NArg() != 0 {
		sp_args.in_filename = flag.Args()[0]
	}
}

//execute the command according to the arguments
func process_input() {
	//if the inputfile name is nil then inputstream is stdin defaultly
	//or open the file
	if sp_args.in_filename != "" {
		inputStream, err = os.Open(sp_args.in_filename)
		//print error
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s : ", progname)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(9)
		}
	}
	//get reader from File
	myReader := bufio.NewReader(inputStream)
	//stdout defaultly is os stdout or output to the printer
	stdout = os.Stdout
	if sp_args.print_dest != "" {
		var printErr error
		printer := exec.Command("lp", "-d", sp_args.print_dest)
		//make the outputstream to the printer's input stream
		stdout, printErr = printer.StdinPipe()

		if printErr != nil {
			fmt.Fprintf(os.Stderr, "%s : ", progname)
			fmt.Fprintln(os.Stderr, printErr)
			os.Exit(16)
		}
		//output according to the page type
		if sp_args.page_type == 'l' {
			outputWithModeL(myReader)
		} else {
			outputWithModeF(myReader)
		}
		stdout.Close()
		printer.Stdout = os.Stdout
		printer.Stderr = os.Stderr
		err := printer.Start()

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s : ", progname)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(17)
		}

		printer.Wait()

	}

	//output according to the page_type
	if sp_args.page_type == 'l' {
		outputWithModeL(myReader)
	} else {
		outputWithModeF(myReader)
	}
}

func outputWithModeL(reader *bufio.Reader) {
	var line_cur uint = 0 //record the line number which is being read
	var page_cur uint = 1 //record the page number which is being read
	for {
		//page_number is greater to the end _page then break
		if page_cur > sp_args.end_page {
			break
		}
		//read the File with delim "'\n'" | read one line
		stringBuf, err := reader.ReadString('\n')
		if page_cur < sp_args.start_page {
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s :", progname)
				fmt.Fprintln(os.Stderr, err)
				os.Exit(10)
			} else {
				//refresh the line_cur and page_cur
				line_cur++
				if line_cur == sp_args.page_len {
					line_cur = 0
					page_cur++
				}
				continue
			}
		}
		if err == nil {
			//output
			stdout.Write([]byte(stringBuf))
			line_cur++
			if line_cur == sp_args.page_len {
				line_cur = 0
				page_cur++
			}
			continue
			//to the EOF
		} else if err == io.EOF {
			stdout.Write([]byte(stringBuf))
			if page_cur < sp_args.end_page {
				fmt.Fprintf(os.Stderr, "%s: end_page is greater than total pages ", progname)
				os.Exit(11)
			}
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "%s :", progname)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(12)
		}

	}
}

func outputWithModeF(reader *bufio.Reader) {
	var page_cur uint = 1

	for {
		if page_cur > sp_args.end_page {
			break
		}

		stringBuf, err := reader.ReadString('\f')
		if page_cur < sp_args.start_page {
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s :", progname)
				fmt.Fprintln(os.Stderr, err)
				os.Exit(13)
			} else {
				//refresh the line_cur and page_cur
				page_cur++
				continue
			}
		}

		if err == nil {
			stdout.Write([]byte(stringBuf))
			page_cur++

		} else if err == io.EOF {
			if page_cur < sp_args.end_page {
				fmt.Fprintf(os.Stderr, "%s: end_page is greater than total pages ", progname)
				os.Exit(14)
			}
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "%s :", progname)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(15)
		}

	}

}

package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"github.com/fatih/color"
)

type Config struct{
	Db_user string 
	Db_pass string
	Wp_user string
	Wp_pass string
	Wp_email string
	Xampp string
}
func GetConf() Config{
	file,_:= os.ReadFile("xampress/config.json")
	var dec_Conf Config
	json.Unmarshal(file, &dec_Conf)
	return dec_Conf
}
//Removes spaces from project name if have any
func GenProj(projName string)(string, string){
	var base_fldr = strings.ToLower(strings.ReplaceAll(projName, " ", "-"))
	var base_db = strings.ToLower(strings.ReplaceAll(projName, " ", "_"))
	return base_db,base_fldr
}
//Check internet connection
func CheckConnection() bool {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
//Download wordpress
func WPDown(projName string) bool{
	var base_fldr = strings.ToLower(strings.ReplaceAll(projName, " ", "-"))
	if ProjExist(base_fldr){
		PrintScrn("'"+projName+"' site already exists please use another name...", "red", 0, false)
		os.Exit(0)
	}
	ferr := os.Mkdir("htdocs/"+base_fldr, 0775)
	Chk_error(ferr, "Encountering error during creating project folder "+base_fldr)
	downCore := exec.Command("php/php", "php/wp-cli.phar", "core", "download", "--path=htdocs/"+base_fldr)
	_, err := downCore.CombinedOutput()
	Chk_error(err, "Encountering error while executing WP-Cli command")
	return true
}

//Print on terminal
func PrintScrn(usr_text string, tcolor string, prev_len int, clr_ovr bool) {
	var main_clr color.Attribute

	switch tcolor {
	case "green":
		main_clr = color.FgGreen
	case "red":
		main_clr = color.FgRed
	case "blue":
		main_clr = color.FgBlue
	case "yellow":
		main_clr = color.FgYellow
	case "cyan":
		main_clr = color.FgCyan
	case "magenta":
		main_clr = color.FgMagenta
	default:
		main_clr = color.FgWhite
	}

	if clr_ovr {
		clearScreen()
	}

	c := color.New(main_clr)
	c.Print(usr_text, strings.Repeat(" ", prev_len))
}

// To Clear Screen
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
// Check error and panic
func Chk_error(stus error, msg string) {
	if stus != nil {
		PrintScrn("\n"+msg+"\n", "red", 0, false)
		os.Exit(0)
	}
}
// To migrate files
func FileMgrt(sourceDir, targetDir string) bool {
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return false
	}

	// Iterate over the contents
	for _, file := range files {
		sourcePath := filepath.Join(sourceDir, file.Name())
		targetPath := filepath.Join(targetDir, file.Name())

		// Move the file or folder
		err := os.Rename(sourcePath, targetPath)
		if err != nil {
			return false
		}
	}

	return true
}

// To verify xampress inside Xampp or not
func Chk_dir() bool {
	check_path := filepath.Join("htdocs")
	if _, err := os.Stat(check_path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}

}
// Check Projecy exist or not
func ProjExist(ProjName string) bool {
	checkPath := filepath.Join("htdocs", ProjName)
	if _, err := os.Stat(checkPath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func Down_wpCli() bool{
	var wp_cli string = filepath.Join("php", "wp-cli.phar")
		resp,_ :=http.Get("https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar")
		if !CheckConnection(){
			PrintScrn("Your device isn't connected to internet, Please connect and try again\n","red", 0,false)
			os.Exit(0)
		}
		defer resp.Body.Close()
		wp_cli_file, cli_err := os.Create(wp_cli)
		Chk_error(cli_err, "Encountering error during downloading WP-Cli")	
		defer wp_cli_file.Close()
		if _, err := io.Copy(wp_cli_file, resp.Body); err != nil {
			return false
		}

		return true
}

func Chk_wpcli() bool {
	var wp_cli string = filepath.Join("php", "wp-cli.phar")
	if _, err := os.Stat(wp_cli); os.IsNotExist(err) {
		return false
	}else{
		return true
	}

}

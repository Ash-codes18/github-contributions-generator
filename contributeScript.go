package main

import (
    "flag"
    "fmt"
    "os"
    "log"
    "time"
    "strings"
    "strconv"
    "math/rand"
    "math"
    "os/exec"
)

const FILE_NAME string = "data.txt"
const DATE_FORMAT string = "2006-01-02 15:04:05"

const ROWS int = 7
const COLUMNS int = 5

func writeErrorMessage(err error) {
    if err != nil {
        log.Println("There has been an error: ",err)
        return
    }
}


func contribute(commit_date string) {
    // usefull links: https://golangbot.com/write-files/

    // file mode is in octal notation for the user, group and other (indicated by the 0, then we have 6 for 110, and 4 for 100. rwx)
    // https://stackoverflow.com/a/18415935/7973144 and https://ss64.com/bash/chmod.html and https://docs.nersc.gov/filesystems/unix-file-permissions/
    file, err := os.OpenFile(FILE_NAME, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        log.Println("OpenFile err: ",err)
        return
    }
    defer file.Close()

    _, err2 := file.WriteString(commit_date+"\n\n")

    if err2 != nil {
        log.Fatal("WriteString err: ", err2)
    }
    //https://pkg.go.dev/os#File.Sync just in case we flush.
    file.Sync()
    //randomNum := rand.Intn(100 - 1) + 3
    exec.Command("git", "add", ".").Run()
    exec.Command("git", "commit", "-m", "Commit date was: "+ commit_date, "--date", commit_date).Run()

}

func generateDate(commit_limit, frequency int, time_period [2]int) {
    currentTime := time.Now()

    //https://pkg.go.dev/time#Time.AddDate -> AddDate(year,month,day)
    //startCommitDate := currentTime.AddDate(-1,0,0)
    var date string
    //rand.Intn(max - min) + min
    for i := 0; i < time_period[1]; i++ {
        rndNumOfCommits := rand.Intn(commit_limit - 1) + 1
        n := 0
        if rand.Intn(100 + 1) <= frequency {
            for n < rndNumOfCommits {
                date = currentTime.AddDate(-1,time_period[0],i).Format(DATE_FORMAT)
                n++
                contribute(date)
            }
        }
    }

}
func Date(year, month, day int) time.Time {
    return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
func getNumberOfDaysBetweenMonths(startMonth, endMonth int) int {
    currentYear, _, _ := time.Now().Date()
    return int(math.Round(Date(currentYear, endMonth, 0).Sub(Date(currentYear, startMonth, 0)).Hours() / 24))
}

func contributionsPerDay(num int) int {
    if num >= 15 {
        num = 15
    }
    if num < 2 {
        num = 2
    }
    return num
}

func contributins_specific_months(specified_months string) [2]int {
    //allright, this isn't the best solution but, its currently the only one I can come up with

    //an array that stores 2 values. The firs value defines the start month and the second is the number of days that the 2 month have between them
    //so when we read the [1] value we know how many times we have to contribute, untill we reach the desired time period, default walue is 365 days.
    months := [2]int {0,12}

    _, currentMonth, _ := time.Now().Date() //a hack to get the month type in int. The Month() function returns the name of the month in a string.

    if strings.Contains(specified_months,"-") {
        //if its the default value we just commit the entire year.
        if specified_months == "1-12" {
            months[0] = 0
            months[1] = getNumberOfDaysBetweenMonths(0, 12)
        } else {
           
        sliceOfMonths := strings.Split(specified_months, "-")

        startMonth, err := strconv.Atoi(sliceOfMonths[0])
        writeErrorMessage(err)

        endMonth, err1 := strconv.Atoi(sliceOfMonths[1])
        writeErrorMessage(err1)
        
        if startMonth - int(currentMonth) < 0 {
            months[0] = 0
        } else {
            months[0] = startMonth - int(currentMonth) //if we want to commit on specific months we have to subtract the month we are currently in.
        }
        
        months[1] = getNumberOfDaysBetweenMonths(startMonth, endMonth  + 1)
        fmt.Println("Number of days with possible commits:", months[1]) 
        }
    } else {
       startMonth, err := strconv.Atoi(specified_months)
       writeErrorMessage(err)
   
       if startMonth - int(currentMonth) < 0 {
            months[0] = 0
        } else {
            months[0] = startMonth - int(currentMonth) //if we want to commit on specific months we have to subtract the month we are currently in.
        }
        // get the number of days in one specific month. we use startMonth +1 because the functions calculates the number of days between 2 months
       months[1] = getNumberOfDaysBetweenMonths(startMonth,startMonth + 1)
       fmt.Println("Number of days with possible commits:", months[1]) 
    }

    fmt.Println("the functions contribute_specific_months_returned: ", months)
    return months
}

func runScript(repository, timePeriod string, commit_limit, frequency int) {

    out := os.MkdirAll("randomContributions", os.ModeDir)

    of := os.Chdir("randomContributions")
    os.RemoveAll(".git")
    os.RemoveAll("data.txt")
    exec.Command("git", "init").Run()

    if of != nil {
        log.Fatal(of)
    }

    if repository != ""{
        generateDate(contributionsPerDay(commit_limit), frequency, contributins_specific_months(timePeriod))
    } else {
        fmt.Println("Holdup. You just wanted to run a github contrubutins script without entering a github repo? Try again.")
    }
    
    //TODO: error handling of the git commands -> https://pkg.go.dev/os/exec#Command
    exec.Command("git", "branch", "-M", "main").Run()
    exec.Command("git", "remote","add", "origin", repository).Run()

    exec.Command("git", "push", "-u", "origin", "main").Run()

    fmt.Println("Done!")

    if out == nil  {
        fmt.Println("Os command Successfully Executed")
        
    } else {
        fmt.Printf("%s", out)
    }
}

var bMatrix = [][]int{   {0,1,0,0,0}, 
{0,1,0,0,0}, 
{0,1,1,1,0},
{0,1,0,0,1}, 
{0,1,0,0,1}, 
{0,1,1,1,0}, 
{0,0,0,0,0}}

func runNonRandomScript(message string) {
    os.MkdirAll("nonrandomContributions", os.ModeDir)

    of := os.Chdir("nonrandomContributions")
    os.RemoveAll(".git")
    os.RemoveAll("data.txt")
    exec.Command("git", "init").Run()

    if of != nil {
        log.Fatal(of)
    }

   // i := 0
   // j := 0
   //TODO flag to enter what date is in the corner!
    counter := -1
    var date string
    currentTime := time.Now()
    fmt.Println("here")
    for i := 0; i < COLUMNS; i++ {
        for j := 0; j < ROWS; j++ {
            fmt.Println("here1")
            if bMatrix[j][i] == 1{
                 date = currentTime.AddDate(-1, 0, counter).Format(DATE_FORMAT)
                           
                 contributeTmp(date)
                }       
            counter++
            fmt.Println("conunter: ", counter)

            if counter >= 35 {
                // reset the variables to 0 
                //bMatrix = assignMatrix(message)
            }

        }
    }


}

func assignMatrix(message string) [][]int {

    var charOfAlphabet = make([][]int, 2)

    for i := 0; i < len(message); i++ {

   // make a map of alphabeet so you can search for the letter
    charOfAlphabet = bMatrix
    }

    return charOfAlphabet
}

func contributeTmp (date string) {

    file, err := os.OpenFile(FILE_NAME, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        log.Println("OpenFile err: ",err)
        return
    }
    defer file.Close()

    _, err2 := file.WriteString(date+"\n\n")

    if err2 != nil {
        log.Fatal("WriteString err: ", err2)
    }
    //https://pkg.go.dev/os#File.Sync just in case we flush.
    file.Sync()

}

func main() {

    randomFlag := flag.NewFlagSet("random", flag.ExitOnError)

    repository := randomFlag.String("repository","","Enter a link to an empty non-initialized GitHub repository to which you want to push the generated file. The link can be an SSH (assuming you have an ssh key) or the HTTPS format. (e.g., git@github.com:yourusername/yourrepo.git or https://github.com/yourusername/yourrepo.git) ")
    commitLimmit := randomFlag.Int("commit_limit", 7, "Set the limit(n) of commits per single day. The script will randomly commit from 1 to n times a day. The maximum is 15 and the minimum is 1.")
    frequency := randomFlag.Int("frequency", 85, "The procentage of days out of 365 you would like to contribute. E.g., if you enter 20, you will contribute 73 days out of 1 year.")
    timePeriod := randomFlag.String("month","1-12", "Contribute only in a specific period. If you enter 3-5, you will only commit from march(starting the same day of month as today) to may. Entering only one number like 8(october) will prompt the script to commit only on the specified month.")


    nonrandomFlag := flag.NewFlagSet("nonrandom", flag.ExitOnError)
    message := nonrandomFlag.String("message", "hello", "Enter the message you would like to be displayed on the contribution graph. The maximum ammount of characters is 10.")

    if len(os.Args) < 2 {
        fmt.Println("expected 'foo' or 'bar' subcommands")
        os.Exit(1)
    }

    switch os.Args[1] {

    case "random":
        randomFlag.Parse(os.Args[2:])
        fmt.Println("subcommand 'random'")
        fmt.Println("  repository:", *repository)
        fmt.Println("  commit limit:", *commitLimmit)
        fmt.Println("  time period (months):", *timePeriod)
        fmt.Println("  frequency of commits is:", *frequency, " % of the year")
        //fmt.Println("  random values you entered that are up to no good:", randomFlag.Args())
        runScript(*repository, *timePeriod, *commitLimmit, *frequency)
   
    case "nonrandom":
        nonrandomFlag.Parse(os.Args[2:])
        fmt.Println("subcommand 'bar'")
        fmt.Println("  message:", *message)
        fmt.Println("  tail:", nonrandomFlag.Args())
        runNonRandomScript(*message)

    default:
        fmt.Println("Expected 'random' or 'nonrandom' subcommands!\n")
        fmt.Println("Arguments for random flag: \n")
        randomFlag.PrintDefaults()
        fmt.Println("\n\n")
        fmt.Println("Arguments for nonrandom flag: \n")
        nonrandomFlag.PrintDefaults()
        os.Exit(1)
    }
}
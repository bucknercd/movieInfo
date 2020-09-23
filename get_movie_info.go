package main
import (
    "fmt"
    "regexp"
    "strconv"
    "os"
    "net/http"
    "net/url"
    "bufio"
    "io/ioutil"
    "encoding/json"
)

const (
        APIKey = "b43b5bf0"
        BaseURL = "http://www.omdbapi.com"
)

type Movie struct {
    Title string
    Year string
    Type string
    Poster string
    Plot string
}

type Search struct {
    Key []Movie `json:"Search"`
    Error string
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Println("\nEnter movie to search for:")
        scanner.Scan()
        searchTerm := scanner.Text()

        body, err := SearchMovies(searchTerm)
        if err != nil {
            fmt.Println(err)
            continue
        }
        searchResult, err := OutputSearch(body)
        if err != nil {
            fmt.Println(err)
            continue
        }
        if searchResult.Key == nil {
            continue
        }
        fmt.Println("\nEnter the number to download the desired movie infomation:")
        scanner.Scan()
        movieDesired := scanner.Text()
        chosenNum, _ := strconv.Atoi(movieDesired)
        if chosenNum < 1 || chosenNum > 10 {
            fmt.Println("\nInvalid choice.")
            continue
        }
        chosenNum--
        title := searchResult.Key[chosenNum].Title
        fmt.Printf("\nDownloading information for '%s' ...\n", title)
        DownloadAndSave(title)
        fmt.Println("\nDo you wish to quit the program? ['y' or 'yes' to exit]")
        scanner.Scan()
        quit := scanner.Text()
        match, _ := regexp.MatchString(`(?i)^\s*y(?:es)?\s*$`, quit)
        if match {
            break
        }
    }
}

func SearchMovies(searchTerm string) ([]byte, error) {
    query := fmt.Sprintf("%s/?apikey=%s&s=%s&plot=full", BaseURL, APIKey, url.QueryEscape(searchTerm))
    resp, err := http.Get(query)
    defer resp.Body.Close()
    if err != nil {
        return nil, err
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return body, nil
}

func OutputSearch(data []byte) (Search, error) {
    var searchResult Search
    if err := json.Unmarshal(data, &searchResult); err != nil {
        return Search{}, err
    }
    if searchResult.Error != "" {
        fmt.Println("\n", searchResult.Error)
    } else {
        for i, v := range searchResult.Key {
            fmt.Printf("%d Title: '%s'  Year: %s  Type: %s\n", i+1, v.Title, v.Year, v.Type)
        }
    }
    return searchResult, nil
}

func DownloadAndSave(title string) error {
    query := fmt.Sprintf("%s/?apikey=%s&t=%s&plot=full", BaseURL, APIKey, url.QueryEscape(title))
    resp, err := http.Get(query)
    defer resp.Body.Close()
    if err != nil {
        return err
    }
    var movie Movie
    if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
        return err
    }
    movieInfo := fmt.Sprintf("Title: %s\nYear: %s\nType: %s\nPlot: %s\n", movie.Title, movie.Year, movie.Type, movie.Plot)
    fmt.Printf("\n%s", movieInfo)
    os.Mkdir("movies", 0755)
    os.Mkdir("movies/" + title, 0755)
    f, _ := os.Create("movies/" + title + "/" + title + ".txt")
    defer f.Close()
    f.WriteString(movieInfo)
    pwd, _ := os.Getwd()
    res, err := http.Get(movie.Poster)
    if err != nil {
        return err
    }
    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return err
    }
    f_img, _  := os.Create("movies/" + title + "/" + title + ".jpg")
    defer f_img.Close()
    f_img.Write(body)
    fmt.Printf("\nThe movie info was saved to '%s/info/%s'\n", pwd, movie.Title)
    return nil
}

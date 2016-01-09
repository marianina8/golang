package main

import (
    "encoding/xml"
    "fmt"
    "net/http"
    "io/ioutil"
    "os"
    "github.com/tealeg/xlsx"
    "time"
    "strings"
)

type Program struct {
    category string
    title string
    epnum string
    start time.Time
    end time.Time
}

type Rss struct {
    XMLName   xml.Name `xml:"rss"`
    Chann     Channel `xml:"channel"`
}

type Channel struct {
    XMLName   xml.Name `xml:"channel"`
    Title       string `xml:"title"`
    Description string `xml:"description"`
    It []Item   `xml:"item"`
}

type Item struct {
    XMLName   xml.Name `xml:"item"`
    Title  string `xml:"title"`
    PubDate string `xml:"pubDate"`
    Category string `xml:"media:category"`
}

func inTimeSpan(start, end, check time.Time) bool {
    return check.After(start) && check.Before(end)
}

func main() {
    excelFileName := os.Args[1]
    xlFile, err := xlsx.OpenFile(excelFileName)
    if err != nil {
        fmt.Println(err)
    }
    var programs []Program
    for _, sheet := range xlFile.Sheets {
        for i, row := range sheet.Rows {
            if i > 3 {
                title := row.Cells[0].String()
                if len(strings.TrimSpace(title)) != 0 {
                    episodeNum := row.Cells[1].String()
                    temp := strings.Replace(strings.Replace(row.Cells[3].String(),";","",-1),"@","",-1)
                    start, _ := time.Parse("1/2/06", temp)
                    temp = strings.Replace(strings.Replace(row.Cells[4].String(),";","",-1),"@","",-1)
                    end, _ := time.Parse("1/2/06", temp)
                    if inTimeSpan(start,end,time.Now()) {
//                      fmt.Printf("%s\t%s\t%s\t%s\n",title,episodeNum,start.Format("1/2/2006"),end.Format("1/2/2006"))
                        var p Program
                        if len(strings.TrimSpace(episodeNum)) == 0 {
                            p = Program{"Feature", title, "", start, end }
                        } else {
                            p = Program{"Series", title, episodeNum, start, end }
                        }
                        programs = append(programs,p)
                    }
                }                
            }
        }
    }
    for _, prog := range programs {
        fmt.Printf("%s,%s,%s,%s,%s\n",prog.category, prog.title, prog.epnum, prog.start.Format("1/2/2006"),prog.end.Format("1/2/2006"))
    }
    ymdTime := time.Now().Format("2006/01/02")
    uri := "http://smb.cdn.neulion.com/u/pivot/" + ymdTime + "/video.xml"
    response, err := http.Get(uri)
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }
        //fmt.Printf("%s\n", string(contents))

        decoder := xml.NewDecoder(strings.NewReader(string(contents))) 

        for { 
            // Read tokens from the XML document in a stream. 
            t, _ := decoder.Token() 
            if t == nil { 
                break 
            } 
            // Inspect the type of the token just read. 
            switch se := t.(type) { 
            case xml.StartElement: 
                // If we just read a StartElement token 
                // ...and its name is "page" 
                if se.Name.Local == "rss" { 
                    var r Rss 
                    // decode a whole chunk of following XML into the
                    // variable p which is a Page (se above) 
                    decoder.DecodeElement(&r, &se) 
                    // Do some stuff with the page. 
                   fmt.Println(r.Chann.It[0].Title)
                } 
            }
        }
    }
}

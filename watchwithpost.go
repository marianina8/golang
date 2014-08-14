package main

import (
	"fmt"
	"path/filepath"
	"io"
	"os"
	"io/ioutil"
	"encoding/xml"
	"bytes"
	"net/http"
)

import "time"

const (
    // A generic XML header suitable for use with the output of Marshal.
    // This is not automatically added to any output of this package,
    // it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

type XMLBroadcastSched struct {
	XMLName xml.Name `xml:"broadcastSchedule"`
    LastUpdated string `xml:"lastUpdated"`
	ProgramSched XMLNewProgramSched `xml:"ProgramSched"`
}

type XMLNewProgramSched struct {
    XMLName  xml.Name `xml:"ProgramSched"`
    Week     string   `xml:"week,attr"`
    CallSign    string   `xml:"callSign,attr"`
    LogDate   *XMLLogDate `xml:"LogDate"`
}

type XMLProgramSched struct {
    XMLName  xml.Name `xml:"ProgramSched"`
    Week     string   `xml:"week,attr"`
    CallSign    string   `xml:"callSign,attr"`
    LogDates   []*XMLLogDate `xml:"LogDate"`
}

type XMLLogDate struct {
    XMLName  xml.Name    `xml:"LogDate"`
    LogDate  string      `xml:"logDate,attr"`
    Programs   []*XMLProgram `xml:"Program"`
}

func (sched XMLProgramSched) GetLogDate(i int) *XMLLogDate {
	return sched.LogDates[i]
}

type XMLProgram struct {
    XMLName  xml.Name    `xml:"Program"`
    StartTime  string      `xml:"startTime,attr"`
    Id  string      `xml:"id,attr"`
    Name  string      `xml:"name,attr"`
    Length  string      `xml:"length,attr"`
    Episode XMLEpisode `xml:"Episode"`
}


type XMLEpisode struct {
    XMLName  xml.Name    `xml:"Episode"`
    ReferenceNumber string      `xml:"referenceNumber,attr"`
    StartTime  string      `xml:"startTime,attr"`
    Id  string      `xml:"id,attr"`
    Name  string      `xml:"name,attr"`
    Length  string      `xml:"length,attr"`
    Eprating string      `xml:"Eprating,attr"`
    Segments []*XMLSegment `xml:"Segment"`
}

type XMLSegment struct {
    XMLName  xml.Name    `xml:"Segment"`
    StartTime  string      `xml:"startTime,attr"`
    Id  string      `xml:"id,attr"`
    Name  string      `xml:"name,attr"`
    Length  string      `xml:"length,attr"`
}

func dir(thepath string) {
  
	filepath.Walk(thepath,VisitFile)
	return
	
}

func VisitFile(fp string, fi os.FileInfo, err error) error {
  
	t := time.Now()
	
    if err != nil {
        fmt.Println(err) // can't walk here,
        return nil       // but continue walking elsewhere
    }
    if !!fi.IsDir() {
        return nil // not a file.  ignore.
    }
    matched, err := filepath.Match("*.xml", fi.Name())
	
    if err != nil {
        fmt.Println(err) // malformed pattern
        return err       // this is fatal.
		write(os.O_APPEND|os.O_RDWR, t.Local().String() + ": No XML files found.\n") 
    }
    if matched {
	    fmt.Println(fp)
		write(os.O_APPEND|os.O_RDWR, t.Local().String() + ": Processing file: " + fp + "\n") 
		method(fp)
    }
    return nil
}

func ReadProgramSchedule(reader io.Reader) (*XMLProgramSched, error) {
    xmlProgramSched := &XMLProgramSched{}
    decoder := xml.NewDecoder(reader)
	t := time.Now()

    if err := decoder.Decode(xmlProgramSched); err != nil {
		fmt.Printf("Error: %s", err)
		write(os.O_APPEND|os.O_RDWR, t.Local().String() + ": Error decoding file.\n") 
        return nil, err
    }

    return xmlProgramSched, nil
}


func write(flag int, text string) { 
	file := "E:\\WWLog\\log_" + time.Now().Format("010206") + ".txt" 
        f, err:=os.OpenFile(file, flag, 0666) 
        if err != nil { fmt.Println(err); return } 
        n, err := io.WriteString(f, text) 
        if err != nil { fmt.Println(n, err); return } 
        f.Close() 
        data, err := ioutil.ReadFile(file) 
        if err != nil { fmt.Println(err); return } 
        fmt.Println(string(data)) 
} 


func method(fpath string) {

	var xmlProgramSched *XMLProgramSched
    var file *os.File
	t := time.Now()

    defer func() {
        if file != nil {
            file.Close()
        }
    }()

    // Build the location of the xml file
    // filepath.Abs appends the file name to the default working directly
    programsFilePath, err := filepath.Abs(fpath)

    if err != nil {
        panic(err.Error())
    }

    // Open the xml file
    file, err = os.Open(programsFilePath)

    if err != nil {
        panic(err.Error())
    }

	xmlProgramSched, err = ReadProgramSchedule(file)
	
    if err != nil {
        panic(err.Error())
    }
	
    // Generate XML files for Daily Program Schedules and post to HTML
	
	for i := 0; i < len(xmlProgramSched.LogDates); i++ {
	
	    
		fmt.Printf("Posting XML Segment times for date: %s\n",xmlProgramSched.LogDates[i].LogDate)
		
		v := &XMLBroadcastSched{}
		v.LastUpdated = time.Now().UTC().Format(time.RFC3339Nano)
		v.ProgramSched.Week = xmlProgramSched.GetLogDate(i).LogDate
		v.ProgramSched.CallSign = "PIVT"
		v.ProgramSched.LogDate = xmlProgramSched.GetLogDate(i)

		output, err := xml.MarshalIndent(v, "  ", "    ")
		if err != nil {
		fmt.Printf("error: %v\n", err)
		}

		output = []byte(xml.Header + string(output))

		buf := bytes.NewBuffer(output)
		resp, err := http.Post("http://ws-pivot.watchwith.com/pivot/import", "text/xml", buf)
		if err != nil {
		  fmt.Printf("err: %s\n", err)
		}
		defer resp.Body.Close()
		if(resp.StatusCode == 200) {
			fmt.Println("RESPONSE - 200 OK")
		    write(os.O_APPEND|os.O_RDWR, t.Local().String() + ": Post XML Segment times for date: " + xmlProgramSched.LogDates[i].LogDate + " RESPONSE - 200 OK\n") 
			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			if(err2 != nil) {
				write(os.O_APPEND|os.O_RDWR, "ERROR\n") 
			}
			bodyString := string (bodyBytes)
			write(os.O_APPEND|os.O_RDWR, t.Local().String() + ": Full Message: " + bodyString + "\n") 
		}
		
		if(resp.StatusCode == 400) {
			fmt.Println("RESPONSE - 400 ERROR")
			write(os.O_APPEND|os.O_RDWR, t.Local().String() + ": Post XML Segment times for date: " + xmlProgramSched.LogDates[i].LogDate + " RESPONSE - 400 ERROR\n") 
			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			if(err2 != nil) {
				write(os.O_APPEND|os.O_RDWR, "ERROR\n") 
			}
			bodyString := string (bodyBytes)
			write(os.O_APPEND|os.O_RDWR, t.Local().String() + ": Full Message: " + bodyString + "\n") 
		}
		
	}
}

func main() {
  write(os.O_CREATE|os.O_TRUNC|os.O_RDWR, "Log file for " + time.Now().Format("01/02/06") + "\n")
  path := "E:\\WWGripitExports\\"
  dir(path)

}


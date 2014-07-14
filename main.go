package main
import (
    "fmt"
    "os"
    //"text/scanner"
    "strings"
    "flag"
    "regexp"
    "bufio"
    "bytes"
    "strconv"
    "sort"
    //"runtime"
)
//import "sync"

var (
    regexPattern *regexp.Regexp
    tokenType string
)

const (
    DEBUG bool = true
)

type Partial struct {
    key string
    value string
}

type Result struct {
    token string
    counts map[string] int
}

type sortedMap struct {
    m map[string]int
    s []string
}

func (sm *sortedMap) Len() int {
    return len(sm.m)
}
 
func (sm *sortedMap) Less(i, j int) bool {
    return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}
 
func (sm *sortedMap) Swap(i, j int) {
    sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}
 
func sortedKeys(m map[string]int) []string {
    sm := new(sortedMap)
    sm.m = m
    sm.s = make([]string, len(m))
    i := 0
    for key, _ := range m {
        sm.s[i] = key
        i++
    }
    sort.Sort(sm)
    return sm.s
}

func Readln(r *bufio.Reader) (string, error) {
  var (isPrefix bool = true
       err error = nil
       line, ln []byte
      )
  for isPrefix && err == nil {
      line, isPrefix, err = r.ReadLine()
      ln = append(ln, line...)
  }
  return string(ln),err
}

func writeToFile(filePath string, dataToDump string) int{

      fh, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0640)
      if err != nil {
        //panic(err)
    fh, _ = os.Create(filePath)      
    if DEBUG { fmt.Println("Notice: File doesn't exist. Creating it.") }
      }
      defer fh.Close()
      nb,_ := fh.WriteString(string(dataToDump))
      fh.Sync()
      if DEBUG { fmt.Println("Notice: Wrote "+strconv.Itoa(nb)+" bytes to "+filePath) }
      return nb
}

func GetUserAgentDetails(ua string) map[string]string{

     //ua = strings.ToLower(ua);
     //matches = regexp.MustCompile(`(?i)(Windows NT\s+[0-9]\.[0-9]|Android|iOS|FirefoxOS|Windows\s*Phone OS [0-9]\.[0-9]|BlackBerry [0-9]{4,4}|BB10)`).FindStringSubmatch(ua)
     matches := regexp.MustCompile(`(?i)(Windows NT|Android|iOS|Firefox|Windows\s*Phone OS|BlackBerry|BB10|iphone os|ipad|ipod|Macintosh|SymbianOS|Series60)`).FindStringSubmatch(ua)
     deviceData := map[string]string{}

     if len(matches) >= 2 {

         switch strings.ToLower(matches[1]) {

                 case "windows nt":

             matches = regexp.MustCompile(`(?i)Windows NT\s+([0-9]+\.[0-9]+)`).FindStringSubmatch(ua)       
             if len(matches) >= 2 {
                deviceData["platform"] = "Windows"
             }

             if matches[1] == "5.1" || matches[1] == "5.2" {
                    deviceData["os_version"] = "XP"
                 } else if matches[1] == "6.0" {
                    deviceData["os_version"] = "Vista"
                 } else if matches[1] == "6.1" {
                deviceData["os_version"] = "7"
             } else if matches[1] == "6.2" {
                deviceData["os_version"] = "8"
             } else if matches[1] == "6.3" {
                    deviceData["os_version"] = "8.1"
                 }

             matches = regexp.MustCompile(`(?i)(ARM|Touch|Tablet)`).FindStringSubmatch(ua)                  
             if len(matches) >= 2 {
                deviceData["ua_type"] = "Mobile"
             } else {
               deviceData["ua_type"] = "Desktop"
             }
          
            case "windows phone":
                 deviceData["platform"] = "Windows Phone"
                 matches = regexp.MustCompile(`(?i)Windows Phone OS\s+([0-9]+\.[0-9]+);`).FindStringSubmatch(ua)
             if len(matches) >= 2 {
                    deviceData["os_version"] = matches[1]
             }

             matches = regexp.MustCompile(`(?i)IEMobile\/([0-9]+\.[0-9]+);`).FindStringSubmatch(ua)
             if len(matches) >= 2 {
                deviceData["rendering_engine"] = "Trident"
                deviceData["browser"] = "Internet Explorer Mobile"
                deviceData["browser_version"] = matches[1]
             }
             deviceData["ua_type"] = "Mobile"

            case "android":
                 deviceData["platform"] = "Android"
             matches = regexp.MustCompile(`(?i)Android\s+([0-9]+\.[0-9]+(\.[0-9]+)*)`).FindStringSubmatch(ua)         
             if len(matches) >= 3 {
                deviceData["os_version"] = matches[1]
             }
             deviceData["ua_type"] = "Mobile"

             matches = regexp.MustCompile(`(?i)(Chrome|Firefox|UCWeb|Maxthon|Opera Mini|Opera|Skyfire|Netfront)`).FindStringSubmatch(ua)
             if len(matches) >= 1 {
                deviceData["browser"] = matches[1]
             }

            case "ios", "iphone os", "ipad", "ipod":
                 deviceData["platform"] = "iOS"
                 matches = regexp.MustCompile(`(?i)OS\s+([0-9]+_[0-9](_[0-9]+)?)`).FindStringSubmatch(ua)
             if len(matches) >= 2 {
                    deviceData["os_version"] = strings.Replace(matches[1], "_", ".", -1)
             }
             deviceData["manufacturer"] = "Apple"
             deviceData["ua_type"] = "Mobile"

              case "macintosh":
                 deviceData["platform"] = "Mac OSX"
                 matches = regexp.MustCompile(`(?i)Version\/([0-9]+\.[0-9]+\.[0-9]+)\s+Safari\/([0-9]+\.[0-9]+\.[0-9]+)`).FindStringSubmatch(ua)
             if len(matches) >= 2 {
                deviceData["browser"] = "Safari"
                deviceData["rendering_engine"] = "WebKit"
             }

             matches = regexp.MustCompile(`(?i)OS\s+X\s+([0-9]+_[0-9]+_[0-9]+)`).FindStringSubmatch(ua)
                 if len(matches) >= 2 {
                    deviceData["os_version"] = strings.Replace(matches[1], "_", ".", -1)
                 }
             deviceData["ua_type"] = "Desktop"

            case "firefox":
                 // Should match:  mozilla/5.0 (mobile; rv:18.0) gecko/18.0 firefox/18.0
             matches = regexp.MustCompile(`(?i)mozilla\/5\.0\s+\(([^;]+;)+\s+rv:[0-9]+\.[0-9]+\)\s+gecko\/[0-9]+\.[0-9]+\s+firefox\/([0-9]+\.[0-9]+)`).FindStringSubmatch(ua)
                 //matches = regexp.MustCompile(`(?i)Android\s+([0-9]+\.[0-9]+(\.[0-9]+)*)`).FindStringSubmatch(ua)
                 if len(matches) >= 3 {
                if matches[1] == "mobile" {
                       deviceData["platform"] = "FirefoxOS"
                   deviceData["os_version"] = matches[2]
                   deviceData["ua_type"] = "Mobile"
                }
                 }

            case "blackberry", "bb10":
                 deviceData["platform"] = "BlackBerry"
             deviceData["manufacturer"] = "RIM"
                 matches = regexp.MustCompile(`(?i)(Version/([0-9]+\.[0-9]+(\.[0-9]+)*))`).FindStringSubmatch(ua)
             if len(matches) >= 2 {
                    deviceData["os_version"] = matches[2]
             }
             matches = regexp.MustCompile(`BlackBerry ([0-9]{4,4});`).FindStringSubmatch(ua)
                 if len(matches) >= 2 {
                    deviceData["model"] = matches[1]
                deviceData["rendering_engine"] = "Mango"
                 }
             deviceData["ua_type"] = "Mobile"


            case "symbianos","series60":
                  deviceData["platform"] = "SymbianOS"
             deviceData["manufacturer"] = "Nokia"         
                 matches = regexp.MustCompile(`(?i)(Series60|SymbianOS)\/([0-9]+\.[0-9]+)`).FindStringSubmatch(ua)
                 if len(matches) >= 3 {
                    deviceData["os_version"] = matches[2]
                 }
             deviceData["ua_type"] = "Mobile"


         } // End switch statement

     } // End outer regex if


     // Try to determine the device manufacturer
     _,ok := deviceData["manufacturer"]
     if !ok {
           matches = regexp.MustCompile(`(?i)(Acer|Archos|benQ| SIE|GeeksPhone|HTC|Huawei|INQ|Kyocera|Lenovo| LG|Meizu|NEC|Nokia|Palm|Pantech|Samsung|Sanyo|Sharp|ZTE)`).FindStringSubmatch(ua)
      if len(matches) >= 2 {
             deviceData["manufacturer"] = strings.Trim(matches[1], " ")
         _,ok := deviceData["ua_type"]
         if !ok {
             deviceData["ua_type"] = "Mobile"
             }
          }   
     }

     // Try one final attempt to detect the rendering engine
     _,ok = deviceData["rendering_engine"]
     if strings.Contains(strings.ToLower(ua), "webkit") && !ok {
         deviceData["rendering_engine"] = "WebKit"
     } else if strings.Contains(strings.ToLower(ua), "gecko") && !ok {
           deviceData["rendering_engine"] = "Gecko"
     } else if strings.Contains(strings.ToLower(ua), "trident") && !ok {
        deviceData["rendering_engine"] = "Trident"
     } else if strings.Contains(strings.ToLower(ua), "presto") && !ok {
        deviceData["rendering_engine"] = "Presto"
     } else if strings.Contains(strings.ToLower(ua), "netfront") && !ok {
        deviceData["rendering_engine"] = "NetFront"
     } else if strings.Contains(strings.ToLower(ua), "obigo") && !ok {
        deviceData["rendering_engine"] = "Obigo"
     }


     // Try one final attempt to detect the browser name
     _,ok = deviceData["browser"]

     if deviceData["platform"] == "iOS" && strings.Contains(strings.ToLower(ua), "safari") && !ok {
        deviceData["browser"] = "Safari Mobile"
     } else if strings.Contains(strings.ToLower(ua), "MSIE") && !ok {
           deviceData["browser"] = "Internet Explorer"
     } else if !ok {
           matches = regexp.MustCompile(`(?i)(Opera Mini|Opera|Skyfire|Chrome|Bolt|Blazer|Series60|UCBrowser)`).FindStringSubmatch(ua)
        if len(matches) >= 2 {
           deviceData["browser"] = matches[1]
        }
     }


     // Now set the default values if fields are empty
     if _,ok := deviceData["platform"]; !ok {
         deviceData["platform"] = "Unknown"
     }
     if _,ok := deviceData["os_version"]; !ok {
          deviceData["os_version"] = "Unknown"
     }
     if _,ok := deviceData["model"]; !ok {
          deviceData["model"] = "Unknown"
     }
     if _,ok := deviceData["rendering_engine"]; !ok {
           deviceData["rendering_engine"] = "Unknown"
     }
     if _,ok := deviceData["browser"]; !ok {
          deviceData["browser"] = "Unknown"
     }
     if _,ok := deviceData["manufacturer"]; !ok {
          deviceData["manufacturer"] = "Unknown"
     }


     // Attempt to confirm it's a bot
     matches = regexp.MustCompile(`(?i)(Googlebot|Baiduspider|YandexBot|YandexWebmaster|Bingbot|MSNbot|NaverBot|Yeti|Exabot|AhrefsBot|cURL)`).FindStringSubmatch(ua)
     fmt.Println("Bot matches:", matches)
     if len(matches) >= 2 {
        deviceData["ua_type"] = "Bot"
        deviceData["bot_type"] = matches[1]
     } else if _,ok := deviceData["ua_type"]; !ok {
          deviceData["ua_type"] = "Desktop"
     }



    

     return deviceData

}


func Map(fileName string, intermediate chan Partial) {

    fmt.Println("\tDEBUG: Mapping file:", fileName)
    // 37.58.100.171 - - [30/Jun/2014:13:03:41 -0400] "GET /default/device_details/name/Apple-iPhone/id/6eeb0707bcd564f39c91cc669df5dd60/dvid/169 HTTP/1.1" 200 9479 "-" "Mozilla/5.0 (compatible; AhrefsBot/5.0; +http://ahrefs.com/robot/)"
    f, err := os.Open(fileName)
    if err != nil {
        fmt.Printf("error opening file: %v\n",err)
        os.Exit(1)
    }
    r := bufio.NewReader(f)

    s, e := Readln(r)
    for e == nil {

        //fmt.Println(s)
        matches := regexPattern.FindStringSubmatch(s)

        if tokenType == "ip" && len(matches) >= 1 {
            intermediate <- Partial{matches[1], fileName}
        } else if tokenType == "uri" && len(matches) >= 5 {
            intermediate <- Partial{matches[5], fileName}
        } else if tokenType == "referer" && len(matches) >= 8 {
            intermediate <- Partial{matches[8], fileName}
        } else if tokenType == "useragent" && len(matches) >= 9 {
            intermediate <- Partial{matches[9], fileName}
        }

        s,e = Readln(r)
    }

    intermediate <- Partial{"", ""}
}

func Reduce(token string, files []string, final chan Result) {
    counts := make(map[string] int)
    for _, file := range files {
        counts[file]++
    }
    final <- Result{token, counts}
}

func collectPartials(intermediate chan Partial, count int, final chan map[string] map[string] int) {

    intermediates := make(map[string] []string)
    for count > 0 {
        res := <- intermediate
        if res.value == "" && res.key == "" {
            count--
        } else {
            v := intermediates[res.key]
            if v == nil {
                v = make([]string, 0, 10)
            }
            v = append(v, res.value)
            intermediates[res.key] = v
        }
    }

    collect := make(chan Result)
    for token, files := range intermediates {
        go Reduce(token, files, collect)
    }

    results := make(map[string] map[string] int)

    // Collect one result for each goroutine we spawned
    for _, _ = range intermediates {
        r := <- collect
        results[r.token] = r.counts
    }
    final <- results
}

func main() {

    //runtime.GOMAXPROCS(runtime.NumCPU())

    var fileDir, outFile string
    flag.StringVar(&fileDir, "d", ".", "Directory to scan for files")
    flag.StringVar(&outFile, "o", "map-reduce-results.txt", "File to which results will be saved.")
    flag.StringVar(&tokenType, "t", "ip", "Type of token to analyze (ip/useragent/uri/referrer)")
    flag.Parse()

    intermediate := make(chan Partial)
    final := make(chan map[string] map[string] int)
    dir, _ := os.Open(fileDir)
    names, _ := dir.Readdirnames(-1)


    validTokenTypes := map[string]int{"ip":1,"uri":1, "useragent":1, "referrer":1}
    if _,ok := validTokenTypes[tokenType]; !ok {
        fmt.Println("ERROR: Must provide valid token type: ip, useragent, uri or referrer")
        os.Exit(0)
    }

    // Initialize the regex pattern
    // 37.58.100.171 - - [30/Jun/2014:13:03:41 -0400] "GET /default/device_details/name/Apple-iPhone/id/6eeb0707bcd564f39c91cc669df5dd60/dvid/169 HTTP/1.1" 200 9479 "-" "Mozilla/5.0 (compatible; AhrefsBot/5.0; +http://ahrefs.com/robot/)"
    regexPattern = regexp.MustCompile(`(?i)([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)\s+-(.+)-\s+\[(.+)\]\s+"(GET|POST)\s+(.*)\s+HTTP/1\.[0-2]"\s+([0-9]{3})\s([0-9]+)\s+"(.+)"\s+"(.+)"`)

    fmt.Println("\tDEBUG: Names:", names)

    go collectPartials(intermediate, len(names), final)

    for _, file := range names {
        if (strings.HasSuffix(file, ".log") || strings.HasSuffix(file, ".txt") ) {
            go Map(fileDir+"/"+file, intermediate)
        } else {
            intermediate <- Partial{"", ""}
        }
    }

    result := <- final
    var buffer bytes.Buffer

    buffer.WriteString("token,total occurences\n")

    flattendMap := map[string]int{}

    // Add up all the counts for the occurences of the token in each file
    for token, counts := range result {

        //fmt.Printf("\n\nToken: %v\n", token)
        total := 0
        //for file, count := range counts {
        for _,count := range counts {   
            //fmt.Printf("\t%s:%d\n", fileDir+"/"+file, count)
            total += count
        }
        flattendMap[token] = total
        //fmt.Printf("Total: %d\n", total)
    }
   

    // Sort the map
    sortedFlattendMap := sortedKeys(flattendMap)


    // Buffer the content, then write the result to file
    for _,token := range sortedFlattendMap {
        buffer.WriteString(token)
        buffer.WriteString(",")
        buffer.WriteString(strconv.Itoa(flattendMap[token]))
        buffer.WriteString("\n")
    }
    writeToFile(outFile, buffer.String())


}
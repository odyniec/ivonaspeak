package main

import (
    "bufio"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
    "unicode"
    "unicode/utf8"

    ivona "github.com/jpadilla/ivona-go"
)

var client *ivona.Ivona

func ucFirst(s string) string {
    if s == "" {
        return ""
    }
    r, n := utf8.DecodeRuneInString(s)
    return string(unicode.ToUpper(r)) + s[n:]
}

func listVoices(filter string) {
    voiceOptions := ivona.Voice{}

    if filter != "all" {
        for _, s := range strings.Split(filter, ",") {
            if strings.ToLower(s) == "male" {
                voiceOptions.Gender = "Male"
            } else if strings.ToLower(s) == "female" {
                voiceOptions.Gender = "Female"
            } else {
                voiceOptions.Language = s
            }
        }
    }

    response, _ := client.ListVoices(voiceOptions)

    for _, voice := range response.Voices {
        fmt.Printf("%s, %s, %s\t(-name=%s -gender=%s -language=%s)\n",
            voice.Name, voice.Gender, voice.Language,
            voice.Name, voice.Gender, voice.Language)
    }
}

func main() {
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, `Usage: %s [options] [text] output-file

  -access-key KEY            Ivona Speech Cloud access key
  -secret-key KEY            Ivona Speech Cloud secret key
  -name NAME                 Voice name (e.g., "Eric")
  -gender GENDER             Voice gender (e.g., "male")
  -language LANGUAGE         Voice language (e.g., "en-US")
  -list-voices FILTER        List available voices and exit
                             (FILTER examples: "all", "male", "male,en-US")
  -text-file                 File containing text to be synthesized
                             (e.g., "speech.txt", "anecdote.xml")
  -format FORMAT             Format of text data
                             ("text" or "ssml", default: "text")
`,
            os.Args[0])
    }

    var voicesFilter string
    var accessKey, secretKey string
    var voiceName, voiceGender, voiceLanguage string
    var textFile, textFormat string

    flag.StringVar(&accessKey, "access-key", "",
        "Ivona Speech Cloud access key")
    flag.StringVar(&secretKey, "secret-key", "",
        "Ivona Speech Cloud secret key")
    flag.StringVar(&voicesFilter, "list-voices", "",
        "List available voices and exit")
    flag.StringVar(&voiceName, "name", "", "Voice name")
    flag.StringVar(&voiceGender, "gender", "", "Voice gender")
    flag.StringVar(&voiceLanguage, "language", "", "Voice language")
    flag.StringVar(&textFile, "text-file", "",
        "File containing text to be synthesized")
    flag.StringVar(&textFormat, "format", "", "Format of text data")

    flag.Parse()

    if accessKey == "" || secretKey == "" {
        // Keys not set -- try getting them from a config file
        configFile := filepath.Join(os.Getenv("HOME"), ".ivonaapi")
        file, err := os.Open(configFile)

        if err != nil {
            log.Fatal(err)
        }

        scanner := bufio.NewScanner(file)
        if scanner.Scan() { accessKey = scanner.Text() }
        if scanner.Scan() { secretKey = scanner.Text() }
    }

    client = ivona.New(accessKey, secretKey)

    if voicesFilter != "" {
        // List voices and exit
        listVoices(voicesFilter)
        return
    }

    var text string
    
    if textFile != "" {
        contents, err := ioutil.ReadFile(textFile)

        if err != nil {
            log.Fatal(err)
        }

        text = string(contents)
    }
    
    if text == "" {
        // Make sure we have at least two arguments (text and output file)
        if len(flag.Args()) < 2 {
            flag.Usage()
            return
        }

        text = flag.Args()[0]
    } else {
        // Make sure we have at least one argument (output file)
        if len(flag.Args()) < 1 {
            flag.Usage()
            return
        }
    }

    outputFile := flag.Args()[len(flag.Args())-1]

    speechOptions := ivona.NewSpeechOptions(text)

    if len(voiceName) > 0 { speechOptions.Voice.Name = voiceName }
    if len(voiceGender) > 0 {
        // Make sure the value is capitalized properly
        speechOptions.Voice.Gender = ucFirst(strings.ToLower(voiceGender))
    }
    if len(voiceLanguage) > 0 { speechOptions.Voice.Language = voiceLanguage }

    if strings.ToLower(textFormat) == "ssml" {
        speechOptions.Input.Type = "application/ssml+xml"
    } else {
        speechOptions.Input.Type = "text/plain"
    }
    
    r, err := client.CreateSpeech(speechOptions)

    if err != nil {
        log.Fatal(err)
    }

    ioutil.WriteFile(outputFile, r.Audio, 0644)
}

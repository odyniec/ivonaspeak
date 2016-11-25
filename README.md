IvonaSpeak
==========

Generate synthesized speech using the Ivona Speech Cloud API.

### Usage

    ivonaspeak [options] [text] output-file

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

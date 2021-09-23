# gslides - tool to work with the Google Slides API

## install

    go install github.com/emicklei/gslides@latest

## requirements

Get your OAuth client ID credentials:

* Create (or reuse) a developer project at <https://console.developers.google.com>
* Enable Google Slides API at [API library page](https://console.developers.google.com/apis/library)
* Go to [Credentials page](https://console.developers.google.com/apis/credentials) and click "+ Create credentials" at the top
* Select "OAuth client ID" authorization credentials
* Choose type "Computer Application" and give it some name.
* Download client credentials file.
* Copy it to `gslides.json` (name has to match).

## usage

Commands expect an identifier of a Google slidedeck, such as `1EA.......C6Vuc`.
Use the flag "-v" for verbose logging.

### thumbnails

Create PNG file for each slide in a presentation.
 
    gslides export thumbnails <source-presentation-id>

### notes

Create TXT file with notes from each slide in a presentation.
    
    gslides export notes <source-presentation-id>

### list

Print the list of presentations with `<document-id>` and `name`.

    gslides list
    gslides list -owner you@company.com

### pdf

Export a presentation (or any document) to a PDF formatted file.
This requires activation of the Drive API in the Google [API library page](https://console.developers.google.com/apis/library).

    gslides export pdf -o mydoc.pdf <document-id>

## append (Work in progress)

Add slide numbers *1* and *2* from the source presentation to the end of target presentation. Use *all* to copy every slide. 
Slides appended will use the layout styling of the target master.

    gslides append <target-presentation-id> <source-presentation-id> 1,2


&copy; 2021+, ernestmicklei.com. MIT License. Contributions welcome.
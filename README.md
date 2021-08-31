# gslides - tool to work with the Google Slides API

## requirements

Then get your OAuth client ID credentials:

* Create (or reuse) a developer project at <https://console.developers.google.com>
* Enable Google Slides API at [API library page](https://console.developers.google.com/apis/library)
* Go to [Credentials page](https://console.developers.google.com/apis/credentials) and click "+ Create credentials" at the top
* Select "OAuth client ID" authorization credentials
* Choose type "Computer Application" and give it some name.
* Download client credentials file.
* Copy it to `gslides.json` (name has to match).

## usage

 Create PNG file for each slide in a presentation.
 
    gslides export thumbnails <source-presentation-id>

 Create TXT file with notes from each slide in a presentation.
    #
    gslides export notes <source-presentation-id>


## work in progress

Add slide numbers *index1* and *index2* from the source presentation to the end of target presentation.
Slides appended will use the layout styling of the target master.

    gslides append <target-presentation-id> <source-presentation-id> index1,index2

Commands expect an identifier of a Google slidedeck, such as `1EA.......C6Vuc`.

Use the flag "-v" for verbose logging.

&copy; 2021+, ernestmicklei.com. MIT License. Contributions welcome.
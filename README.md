# gslides - tool to work with the Google Slides API

## requirements

- GCP project with Slides API enabled
- OAuth 2.0 credential as downloaded JSON and rename it to gslides.json.

## usage

    # create PNG file for each slide in a presentation
    #
    gslides export thumbnails <source-presentation-id>

    # create TXT file with notes from each slide in a presentation
    #
    gslides export notes <source-presentation-id>


## work in progress

    # add slide numbers *index1* and *index2* from the source presentation to the end of target presentation
    # slides appended will use the layout styling of the target master.
    #
    gslides append <target-presentation-id> <source-presentation-id> index1,index2

Commands expect an identifier of a Google slidedeck, such as `1EA.......C6Vuc`.

Use the flag "-v" for verbose logging.

&copy; 2021+, ernestmicklei.com. MIT License. Contributions welcome.
# Forza Data Tools / Forza Stats
Some tools for playing with the UDP data out feature from the Forza Motorsport 7 / Forza Horizon 4 games. Built with [golang](https://golang.org/dl/).  

Forza Stats data writing currently for use in [Forza Horizon 5 Leaderboards and Stat Tools](https://docs.google.com/spreadsheets/d/1UzB2IIzqNqzs9sWWV65w0VVHUmUaeFH1eGlK4-jyNMc/edit?usp=sharing) spreadsheet.




## Forza Data Tools (fdt) Features
- Realtime telemetry output to terminal  
- Telemetry data logging to csv file  
- Serve Forza Telemetry data as JSON over HTTP
- Display race statistics from race/drive (when logging to CSV)

## Forza Stats (writestats) Features
- Calculating race telemetry statistics from csv log
- Reading/Writing to stats spreadsheet through Google Sheets API
- Remotely trigger spreadsheet scripts through Apps Script API



&nbsp;

## Setup
From your game HUD options, enable the data out feature and set it to use the IP address of your computer. Port 9999.  
Forza Motorsport 7 select the "car dash" format.  

### Writestats credentials: 
Change client id and client secret in "credentials_template" to OAuth information for your Google Cloud Project. Then rename to simply "credentials.json" and the program will be able to read credentials properly. This ensures sensitive information is not uploaded publically.  

Once credentials are properly set up, running the writestats program will prompt the user with a URL to paste into a browser. This will direct you to to authorizoation for your Google account and will send you to a URL that contains a code which can be pasted back into the command prompt, thus successfully authorizing the user to write to the target sheet and create a token. Tokens expire every 7 days and once that occurs, simply delete "token.json" that is created from aforementioned process and generate a new one in the same manner. Currently only set up with [Forza Horizon 5 Leaderboards and Stat Tools](https://docs.google.com/spreadsheets/d/1UzB2IIzqNqzs9sWWV65w0VVHUmUaeFH1eGlK4-jyNMc/edit?usp=sharing) spreadsheet, but could be modified for your own spreadsheet by changing the "spreadsheetID" variable in `writestats.go` (May need further setup with Google Cloud for this to work though)

&nbsp;

## Build
Forza Data Tools telemetry processing included as "fdt.exe" (already built)  
To build the writestats application, compile with the command: `go build -o writestats`  

&nbsp;

## Run
### Forza Data Tools command line options
Turn off csv logging: `-c`    
EV mode - enables continuous datastream even in menus (for use in collecting electric vehicle stats): `-e`    
Enable JSON server: `-j`   
Disable realtime terminal output: `-q`   
Enable debug information: `-d`


##### Example 
`fdt`  Runs program  
`fdt -e`  Runs program with EV mode  

### Writestats command line options
Default: writes stat line to sheet and triggers color script to color output data  
Ordinal Info Collection Mode: `-o` Writes ordinal numbers into Ordinal Data sheet  
Race Mode: `-r` Writes race statistics - best lap time and track top speed  

Currently for use in Forza Horizon 5 Leaderboards and Stat Tools Spreadsheet  


##### Example
`writestats`   
`writestats -o`  
`writestats -r`


&nbsp;

### JSON Data
If the `-j` flag is provided, JSON data will be available at: http://localhost:8080/forza. Could be used to make a web dashboard interface or something similar. JSON Format is an array of objects containing the various Forza data types.  

You can see a sample of the kind of data that will be returned [here](https://github.com/richstokes/Forza-data-tools/blob/master/dash/sample.json).  

There is a basic example JavaScript dashboard (with rev limiter function) in the `/dash` directory.  

&nbsp; 

## Further reading
- Forza data out format: https://forums.forzamotorsport.net/turn10_postsm926839_Forza-Motorsport-7--Data-Out--feature-details.aspx#post_926839

- Forza Horizon 4 has some mystery data in the packet, waiting on info from the developers: https://forums.forzamotorsport.net/turn10_postsm1086012_Data-Output.aspx#post_1086012

# Hue Sensor Collector
Application that connects to Hue sensors through Hue bridge and collect temperature readings.

## Application Details
- Written in Golang.
- Intended to be run in RasberryPi, batch file included to run build in Windows and target ARM/Linux.
- Data collected every 10 minutes and saved in Google Cloud Platform Datastore using Go Cloud Datastore package(https://godoc.org/cloud.google.com/go/datastore)
- Configuration files:
  - user.txt - created automatically during login to Hue bridge
  - gcpproject.txt - Create this file with your gcp project id
  - gcp.json - Create this file with your google credential file (refer to https://cloud.google.com/docs/authentication/production)

The purpose of this app is to calculate each of our company credit cardholder's transaction activity for the month, develop a journal entry to apply the payment correctly accross all cards, post it to Quickbooks Online using their API, and alerts me to credits that need to be handled differently.

Requirements:
In order to run this app, you need a few things:
1. Go installation
2. A developer.intuit.com account
3. An app on developer.intuit.com and the associated access key
4. This app reads .csv files with the cardholder's transaction history for the current statement. These must be downloaded from Chase and saved in the csv folder of the app. Before running the script, you must make sure that the card holder's name is added as a field to the csv file. Chase does not include it by default.

First Use Instructions:
1. Clone the GitHub repo to your computer and place it in the src directory of your $GOPATH
2. Set your Quickbooks Online realmID as an environment variable "QBO_realmID_production"
3. Use the OAuth2 Playground at developer.intuit.com to request an access key 
4. When you first set this up, create .csv files for each cardholder for the current statement cycle.
 - f[lower case first initial of cardholder]L[Upper case last initial of cardholder][statementEndingDate with format mmddyy].csv 

Running the code:
1. Access key must be entered as a parameter to the main function.
2. Use the command `go run main.go [access key]` to run script.
3. If the required .csv files are not available, or if the .csv files contain non-digit values, the script will panic and end.
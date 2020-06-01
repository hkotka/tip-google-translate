# tip-google-translate
Google cloud translations API integration for Universal Tip.

This is a provider script for a [Universal Tip](https://github.com/tanin47/tip) program for MacOS.
The script looks for a translation to target language using Google Translate API when Tip is invoked.

##Requirements
* Google Cloud account
* Cloud Translation API Enabled
* Service Account with priveleges to use Translations API
  * credentials.json file for the Service Account accessible for this provider script
  
##Usage
* On main.go, change const values to match your setup
* Run build-deploy.sh
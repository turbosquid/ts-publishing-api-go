# ts-publishing-api-go
Golang implementation of TurboSquid Publishing API

# Installation
You can download a compiled version of the ts-publishing-api-go app for your platform. The unzipped file can be placed anywhere on your computer that you have read/write access.

You can also build your own application as long as you have Go version 1.12 or later. Please follow general Go language instructions.

# Usage
ts-publishing-api-go app is a commandline application. To run it, go to your command prompt or terminal and change the current working directory to the folder that you installed ts-publishing-api-go. Run the following command where "product-folder" is the name of a folder in the current directory that has a product.json definition and files to publish to TurboSquid.

```bash
./ts-publishing-api-go product-folder
```

The first time you run this application, it will ask you for your TurboSquid API Key and it will save that information in settings.yml file for future use.

# TurboSquid Sample Product
We have created a sample product that shows the formatting for product.json that the publishing api app expects. You can download and unzip this sample product into the same directory as the ts-publishing-api-go application.

https://static.turbosquid.com/API/turbosquid-sample-product-1.0.zip

```bash
./ts-publishing-api-go turbosquid-sample-product
```

# TurboSquid Account and API Key
Please ensure that your TurboSquid artist account has agreed to all of the artist license agreements at https://www.turbosquid.com/Seller/.

You can setup an API Key at https://www.turbosquid.com/MemberInfo/, however you currently need to be a member of the API beta group. Please contact support for more information.

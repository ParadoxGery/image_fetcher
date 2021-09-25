# image_fetcher

simple proxy to fetch images

intended to be deployed on the same server as a FoundryVTT application.
sending requests to this server will fetch the image and store it to the servers drive.

## install

the usual `go install`  
sorry for none gophers if you need help create an issue and I will provide releases

## setup

start with `image_fetcher --path="your image base path" --port=1337`  
> you can change the port as you like and there is no need to expose this port to the internet as it should be called from your FoundryVTT server  
> but of course you can expose this to the internet but do so with caution as there is no authorization implemented

## wire

to wire this with your VTT application you need to do calls on the API endpoint `localhost:1337/token`  
add the query params `target=${url encoded uri of the image to fetch}` and `name=${the name you want to give the image}`
if the image already exists with the same name it will NOT be replaced

## contribution

if you feel like adding to this thing feel free to do so  
create issues, fork, create PRs :-)

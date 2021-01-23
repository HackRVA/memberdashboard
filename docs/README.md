# MemberServer Docs

## SwaggerUI
SwaggerUI is the static site for serving up our api documentation.

The static html/js files can be found in the SwaggerUI release (in the dist folder)
https://github.com/swagger-api/swagger-ui

## Generate the documentation
This project is documented with [go-swagger](https://github.com/go-swagger/go-swagger) 

Download the binary from [here](https://github.com/go-swagger/go-swagger/releases)

You will probably have to allow the binary to be executable

```
chmod +x swagger_linux_amd64
```

Then copy/move the binary somewhere you can access it

```
cp swagger_linux_amd64 /usr/local/bin/swagger
```

Once you have your swagger binary, you can run the `gendocs.sh` script in the root of the repo.
This will create the `swagger.json` file and place it in the `./docs/swaggerui/` directory.

You will then be able to access the swagger documentation at `/swaggerui/` in your web browser.

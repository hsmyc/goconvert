# Go File Conversion Service Documentation

## Overview

This Go application is a simple file conversion service that runs a web server to accept zip files containing documents. Users can upload zip files through an HTTP POST request to a specified endpoint. The application will then process each document within the zip file, converting it from one format to another using an external tool (`pandoc`). This service is particularly useful for batch processing of document format conversions.

## Getting Started

### Installation

#### Using Docker
- Docker installed on your system.
1. Pull the Docker image:

   ```bash
   docker pull hsmyc/goconverter
   ```

2. Run the container, mapping the ports if necessary (e.g., 8080):

   ```bash
   docker run -d -p 8080:8080 hsmyc/goconverter
   ```

#### Running Locally

1. Ensure Go is installed on your system.
2. Clone or download the application code to your local machine.
3. Install `pandoc`:
   - On Ubuntu: `sudo apt-get install pandoc`
   - On macOS: `brew install pandoc`
   - For other systems, refer to the [Pandoc installation guide](https://pandoc.org/installing.html).

### Running the Server

To start the server, navigate to the directory containing the application code and run:

```shell
go run .
```

The server will start and listen on `http://localhost:8080`. It's ready to accept file upload requests at the `/convert` endpoint.

## Usage

### Uploading Files for Conversion

To convert documents, you need to send a POST request to `http://localhost:8080/convert` with the following parameters:

- `outputFormat`: The desired output format for the documents (e.g., `markdown`).
- `file`: The zip file containing the documents to be converted.

You can use tools like `curl` or Postman to make the request. Here's an example using `curl`:

```shell
curl -X POST -F "outputFormat=markdown" -F "file=@path_to_your_file.zip" http://localhost:8080/convert
```

Replace `path_to_your_file.zip` with the actual path to your zip file.

### Server Response

After processing the uploaded file, the server responds with a message indicating the success or failure of the upload and processing steps. In case of success, it returns "File uploaded and processed successfully."

## Implementation Details

### Main Components

- **HTTP Server**: Uses the standard Go `net/http` package to listen for incoming HTTP requests.
- **Upload Handler**: A handler function that processes POST requests, extracting the zip file and the desired input and output formats.
- **File Processor**: Processes each file within the uploaded zip, performing the conversion by invoking `pandoc` with the appropriate arguments.
- **Containerization**: Encapsulated in a Docker container for easy deployment and scalability.

### Concurrency

The application utilizes Go's concurrency model (goroutines and wait groups) to process multiple documents within the zip file concurrently. This approach enhances performance, especially when dealing with large numbers of documents.

### Error Handling

The application includes basic error handling to respond appropriately to different failure scenarios, such as unsupported HTTP methods, file processing errors, and internal server errors.

## Limitations

- The application currently supports only the conversion formats available through `pandoc`.
- It's designed to process documents contained within zip files, and other types of archives are not supported.
- Error handling is basic and might need enhancements for production use, including more detailed error messages and logging.

## Future Improvements

- Extend support for other archive formats like `rar` or `tar.gz`.
- Support for additional document and archive formats.
- Enhanced error handling and logging.
- User authentication for secure file uploads.

## Conclusion

This Go application provides a basic but powerful service for converting documents from one format to another in batch mode. It leverages Go's powerful concurrency model and integrates with `pandoc`, offering a flexible solution for document conversion needs. Docker image simplifies document conversion, offering a scalable and easy-to-deploy service for handling various document formats in batch mode.

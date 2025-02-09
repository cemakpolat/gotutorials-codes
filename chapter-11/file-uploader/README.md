A simple file upload service where users can upload files to a server.

Features:

- File upload via HTTP POST
- File storage on the server

Start go app with the command below on a terminal

`go run .`

On the other terminal, use curl to upload a file

 `curl -X POST -F "file=@./test.txt" http://localhost:8080/upload`


Explanation:
- `-X POST`: Specifies that this is a POST request.
- `-F "file=@/path/to/your/file.txt"`: Indicates a form field named file, with the file to upload specified after the `@` symbol.
- Replace `/path/to/your/file.txt` with the actual path to the file you want to upload.
`http://localhost:8080/upload`: The endpoint where the file will be uploaded.

The result shound be seen as below once the file is uploaded: 

`File uploaded successfully: %stest.txt`

# Bookstore
This project provices a simple RESTful API that allows users to perform Create, Read, Update and Delete (CRUD) operations on books. Furthermores, the application downloads 20 books from GoogleBooks and OpenLibrary. The parsing operation due to the varied data structure for these libraries are ddifferent. Each library has its own public API. In this example, all books are stored in the memory, however, it can be combined with a NoSQL/SQL database.

The features of the project includes:

- Basic CRUD endpoints (/books)
- JSON data handling
- Use of mux router

To run the project, you should run the following commands:
```
go mod init bookstore
go mod tidy
go get github.com/gorilla/mux
go run .
```

## How to use the Book Store API
Let's break down the `curl` commands for each operation of your API (Get a book, Create a book, Delete a book, and Update a book)

**Assumptions**

*   Your API is running on `http://localhost:8080`
*   You've already populated some initial books using the google books api.

**1. Get All Books**

*   **HTTP Method:** `GET`
*   **Endpoint:** `/books`

```bash
curl http://localhost:8080/books
```

   This command will fetch and display all the books in JSON format.

**2. Get a Specific Book**

*   **HTTP Method:** `GET`
*   **Endpoint:** `/books/{id}` (replace `{id}` with the actual ID of the book)
    * For this, let's get a book with `id = 2`.

```bash
curl http://localhost:8080/books/2
```

   This command will fetch and display a single book with the specified ID in JSON format.

**3. Create a Book**

*   **HTTP Method:** `POST`
*   **Endpoint:** `/books`
*   **Request Body:** JSON data representing the book to be created
    * For example:
        ```json
        {
            "bookId": "unique_id",
            "title": "The Hitchhiker's Guide to the Galaxy",
            "author": "Douglas Adams"
        }
       ```
*   **`curl` Command:**

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "bookId": "unique_id",
        "title": "The Hitchhiker\'s Guide to the Galaxy",
        "author": "Douglas Adams"
    }' \
    http://localhost:8080/books
```

**Explanation:**

*   `-X POST`: Specifies the HTTP method as `POST`.
*   `-H "Content-Type: application/json"`: Sets the Content-Type header to indicate that the body contains JSON data.
*   `-d '...'`: Provides the JSON data to be sent in the request body.
    * Note that the quotes need to be escaped when using `bash`. Alternatively, you can copy the json to a file called `book.json`, and use `curl -X POST -H "Content-Type: application/json" -d @book.json http://localhost:8080/books`

**4. Delete a Book**

*   **HTTP Method:** `DELETE`
*   **Endpoint:** `/books/{id}` (replace `{id}` with the actual ID of the book to be deleted)
*   For example to delete the book with `id = 1`

```bash
curl -X DELETE http://localhost:8080/books/1
```

**Explanation:**

*   `-X DELETE`: Specifies the HTTP method as `DELETE`.

   This command will delete the book with the specified ID. If successful, it won't display any content, as it returns a 204 response.

**5. Update a Book**

*   **HTTP Method:** `PUT`
*   **Endpoint:** `/books/{id}` (replace `{id}` with the actual ID of the book to be updated)
*   **Request Body:** JSON data representing the updated book details.
     * For example:
        ```json
        {
           "bookId": "unique_id_2",
           "title": "Updated Title",
           "author": "Updated Author"
        }
       ```
*   **`curl` Command:**
   * Let's update a book with `id=2` with the json above.

```bash
curl -X PUT \
  -H "Content-Type: application/json" \
  -d '{
        "bookId": "unique_id_2",
        "title": "Updated Title",
        "author": "Updated Author"
     }' \
  http://localhost:8080/books/2
```

**Explanation:**

*   `-X PUT`: Specifies the HTTP method as `PUT`.
*   `-H "Content-Type: application/json"`: Sets the Content-Type header to indicate that the body contains JSON data.
*   `-d '...'`: Provides the JSON data to be sent in the request body.
    * Note that the quotes need to be escaped when using `bash`. Alternatively, you can copy the json to a file called `book.json`, and use `curl -X PUT -H "Content-Type: application/json" -d @book.json http://localhost:8080/books/2`

**Important Notes**

*   **Response:** The `GET` requests will return a JSON object. The `POST` will return the created book in the JSON body, `PUT` will return the updated book in the JSON body, while the `DELETE` request, when successful will return an empty body, as it sends back a 204 status code.
*   **Errors**: If there is a problem in your requests, the server will send an error code in the header.

These `curl` commands should cover all your API operations, allowing you to test and interact with your endpoints effectively. If you have any more questions, let me know!

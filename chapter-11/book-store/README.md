# Bookstore
A simple RESTful API that allows users to perform Create, Read, Update, and Delete (CRUD) operations on a resource (e.g., books).

Features:

Basic CRUD endpoints (/books)
JSON data handling
Use of mux router

From scrach if you want to implement 
```
go mod init bookstore
go mod tidy
go get github.com/gorilla/mux
go run .
```
Test the code 


GoogleBooksResponse Struct: A new struct GoogleBooksResponse is defined to parse the JSON response from Google Books. This aligns with the format of the API's search endpoint.
fetchBooksFromGoogleBooks Function: This new function is responsible for:
Fetching Data: It uses http.Get to retrieve book data from the Google Books API endpoint. Here we are searching for the the lord of the rings books and limiting the response to three books.
Error Handling: Basic error handling is included for the http get request.
JSON Decoding: The json.Unmarshal function is used to parse the JSON data into the GoogleBooksResponse struct. Error handling is added to make sure if the json is parsed correctly.
Book Creation: It loops through the decoded Google Books response and converts each entry to our internal Book struct.
ID Creation: The id for the book in this case is the id returned by google.
Author: The author is taken from the Authors list. In this case we pick the first author.




OpenLibraryResponse Struct: A new struct OpenLibraryResponse is defined to parse the JSON response from Open Library. This aligns with the format of the API's search endpoint.
fetchBooksFromOpenLibrary Function: This new function is responsible for:
Fetching Data: It uses http.Get to retrieve book data from the Open Library search API endpoint. Here we are searching for the the lord of the rings books and limiting the response to three books.
Error Handling: Basic error handling is included for the http get request.
JSON Decoding: The json.Unmarshal function is used to parse the JSON data into the OpenLibraryResponse struct. Error handling is added to make sure if the json is parsed correctly.
Book Creation: It loops through the decoded Open Library response and converts each entry to our internal Book struct.
ID Creation: The id for the book in this case is the index of the book in the response. You can generate a more unique id if you want to.
Author: The author is taken from the AuthorName list. In this case we pick the first author.
main Function Modification: The fetchBooksFromOpenLibrary() is called before the server starts, seeding the database with initial book data.


Description: A community-driven project aiming to catalog every book ever published. It's a good source for book metadata, editions, and author information.
Strengths:
Free and open, no API keys required for basic use.
Extensive catalog.
Good for data on different editions and authors.
Weaknesses:
May not have as detailed information for every book as Google Books.
Sometimes slower response times.
Does not provide book covers easily.
Example URL:
https://openlibrary.org/search.json?q=the+lord+of+the+rings&limit=3
Use code with caution.

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

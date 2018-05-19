// Requires cpprestsdk-dev (by Microsoft), boost, ssl, crypt. Standard Debian packages available.


#include <cpprest/http_client.h>
#include <cpprest/filestream.h>

using namespace utility;                    // Common utilities like string conversions
using namespace web;                        // Common features like URIs.
using namespace web::http;                  // Common HTTP functionality
using namespace web::http::client;          // HTTP client features
using namespace concurrency::streams;       // Asynchronous streams

// Retrieves a JSON value from an HTTP request.
pplx::task<void> RequestJSONValueAsync()
{
    http_client client(L"http://www.bing.com/");

    // Build request URI and start the request.
    uri_builder builder(U("/search"));
    builder.append_query(U("q"), U("cpprestsdk github"));
    return client.request(methods::GET, builder.to_string());


    return client.request(methods::GET, builder.to_string())
        .then([](http_response response) -> pplx::task<json::value>
    {
        if(response.status_code() == status_codes::OK)
        {
            return response.extract_json();
        }

        // Handle error cases, for now return empty json value... 
        return pplx::task_from_result(json::value());
    })
        .then([](pplx::task<json::value> previousTask)
    {
        try
        {
            const json::value& v = previousTask.get();
            // Perform actions here to process the JSON value...
        }
        catch (const http_exception& e)
        {
            // Print error.
            wostringstream ss;
            ss << e.what() << endl;
            wcout << ss.str();
        }
    });

    /* Output:
    Content-Type must be application/json to extract (is: text/html)
    */
}



int main(int argc, char* argv[])
{
    auto fileStream = std::make_shared<ostream>();

    // Open stream to output file.
    pplx::task<void> requestTask = stringstream::open_ostream(U("results.html")).then([=](ostream outFile)
    {
        *fileStream = outFile;

        // Create http_client to send the request.
        http_client client(U("http://www.bing.com/"));

        // Build request URI and start the request.
        uri_builder builder(U("/search"));
        builder.append_query(U("q"), U("cpprestsdk github"));
        return client.request(methods::GET, builder.to_string());
    })

    // Handle response headers arriving.
    .then([=](http_response response)
    {
        printf("Received response status code:%u\n", response.status_code());

        // Write response body into the file.
        return response.body().read_to_end(fileStream->streambuf());
    })

    // Close the file stream.
    .then([=](size_t)
    {
        return fileStream->close();
    });

    // Wait for all the outstanding I/O to complete and handle any exceptions
    try
    {
        requestTask.wait();
    }
    catch (const std::exception &e)
    {
        printf("Error exception:%s\n", e.what());
    }

    return 0;
}
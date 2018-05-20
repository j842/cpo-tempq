// Requires cpprestsdk-dev (by Microsoft), boost, ssl, crypt. Standard Debian packages available.

#include <iostream>

#include <cpprest/http_client.h>
#include <cpprest/filestream.h>

using namespace utility;                    // Common utilities like string conversions
using namespace web;                        // Common features like URIs.
using namespace web::http;                  // Common HTTP functionality
using namespace web::http::client;          // HTTP client features
using namespace concurrency::streams;       // Asynchronous streams

using namespace std;

// Retrieves a JSON value from an HTTP request.
pplx::task<void> RequestJSONValueAsync()
{

    // /2.2/questions?fromdate=1519862400&todate=1522454400&order=desc&sort=activity&tagged=python&site=stackoverflow
    http_client client(U("http://api.stackexchange.com/"));

    // Build request URI and start the request.
    uri_builder builder(U("/2.2/questions"));
    builder.append_query(U("fromdate"), U("1519862400"));
    builder.append_query(U("todate"), U("1522454400"));
    builder.append_query(U("tagged"), U("python"));
    builder.append_query(U("site"), U("stackoverflow"));


    return client.request(methods::GET, builder.to_string())
        .then([](http_response response) -> pplx::task< utility::string_t >
    {
        if(response.status_code() == status_codes::OK)
        {
            return response.extract_string();
        }

        std::cerr << "Sad status code: " << response.status_code() << std::endl;
        utility::string_t err( U("ERROR"));
        return pplx::task_from_result(err);
    })
        .then([](pplx::task<utility::string_t> previousTask)
    {
        try
        {
            const utility::string_t& v = previousTask.get();
            std::cout << "---" << std::endl << v << std::endl;
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

    std::cout << L"Calling RequestJSONValueAsync..." << endl;
    RequestJSONValueAsync().wait();
    return 0;
}
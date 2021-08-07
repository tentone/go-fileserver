# Go File Server
- Go HTTP file server for REST resource to store and load files.
- Utils to convert file formats, resize and cache images, store folders, generate zip files etc.



### Setup

- Download and install go on your computer, prepare a configuration file.
- The configuration file contains all the parameters necessary to run the server.
- Run the server code.



### Usage

- Data is organized in libraries, each library is stored in a different data folder and has different content.
- The user can create libraries to organize data the way it wants to.
- Libraries can have a type attached, some types contain additional functionality.
  - Images - Stores all types of images (png, jpeg, tiff, webp, gif, ...)
  - Files - Generic storage for files that do not require any type of processing (txt, docx, pdf, ...)
  - Folders - Stores entire folders of data (can be downloaded as zip).
- SQL Metadata database to store data about the existing resources and index them (based on library configuration).



### Data Storage

- Data is stored as files, using the UUID identifier as file name and the correct file format extension.
- Files are organized in library folders, each library should contains a different type of data (images, videos, documents, etc).
- The platform relies on the file system to index and access all data quickly.



### Performance

- Go was selected based on its performance, we compared some solutions using other languages and frameworks and decided that go with net/http was the best compromise between speed and complexity.
- The table bellow presents the performance of fasthttp, node.js express on HTTP 1.x compared against the performance of net/http running on HTTP 1.X and 2.0.

| Upload            | express  (ms) | net/http  (ms) | fasthttp (ms) |
| ----------------- | ------------- | -------------- | ------------- |
| 3000  x 0.0001 Mb | 18903.3       | 12264.40       | 6549.19       |
| 1000  x 0.001 Mb  | 6335.69       | 4825.70        | 2239.20       |
| 500  x 0.01 Mb    | 3072.0        | 2436.79        | 1167.5        |
| 500  x 0.1 Mb     | 3426.89       | 2342.19        | 1456.79       |
| 200  x 1 Mb       | 2076.49       | 1345.60        | 1060.39       |
| 20  x 100 M       | 11960.59      | 17809.69       | 8573.00       |

| Download          | express (ms) | net/http (ms) | fasthttp (ms) |
| ----------------- | ------------ | ------------- | ------------- |
| 3000  x 0.0001 Mb | 6944.1       | 5142.49       | 4237.89       |
| 1000  x 0.001 Mb  | 2651.70      | 1834.90       | 1423.09       |
| 500  x 0.01 Mb    | 1420.5       | 827.30        | 916.50        |
| 500  x 0.1 Mb     | 5201.09      | 4770.19       | 4372.99       |
| 200  x 1 Mb       | 4942.99      | 7406.30       | 6634.50       |
| 20  x 100 Mb      | 27971.20     | 28725.69      | 29549.49      |

- I have also tested the performance of net/http with HTTP 2.0 or fasthttp with HTTP 1.0, for this test I have applied network throttling to limit speed as if a 3G network was in use.

| Upload           | HTTP  2.0  net/http  (ms) | HTTP  1.1  fasthttp (ms) |
| ---------------- | ------------------------- | ------------------------ |
| 200  x 0.0001 Mb | 4178                      | 19585                    |
| 200  x 0.001 Mb  | 7507                      | 20285                    |
| 20  x 0.01 Mb    | 3346                      | 4262                     |
| 20  x 0.1 Mb     | 24169                     | 24925                    |

| Download         | HTTP  2.0  net/http  (ms) | HTTP  1.1  fasthttp (ms) |
| ---------------- | ------------------------- | ------------------------ |
| 200  x 0.0001 Mb | 2158                      | 19467                    |
| 200  x 0.001 Mb  | 2888                      | 19448                    |
| 20  x 0.01 Mb    | 1731                      | 3050                     |
| 20  x 0.1 Mb     | 11220                     | 12491                    |

### License

- This project is distributed under an MIT license available on the project repository.

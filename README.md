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

- Go was selected based on its performance, we compared multiple solutions using other languages and frameworks and decided that go with net/http was the best compromise between speed and complexity.
- The table bellow presents the performance of fasthttp, node.js express on HTTP 1.x compared against the performance of net/http running on HTTP 1.X and 2.0.



### License

- This project is distributed under an MIT license available on the project repository.

# GoDonkey
- GoDonkey resource is a Go lang based REST resource server to store and load files.
- Utils to convert file formats, resize and cache images, store folders, generate zip files etc.
- Support for booth HTTP 2.0 using the http package or HTTP 1.1 using fasthttp.



### Usage

- Data is organized in libraries, each library is stored in a different data folder and has different content.
- The user can create libraries to organize data the way it wants to.
- There are some default libraries that contain additional functionality.
  - Images - Stores all types of images (png, jpeg, tiff, webp, gif, ...)
  - Files - Generic storage for files that do not require any type of processing (txt, docx, pdf, ...)
  - Folders - Stores entire folders of data (can be downloaded as zip).
- Metadata database to store data about the existing resources.



### Data storage

- Data is stored as files, using the UUID identifier as file name and the correct file format extension.
- The platform relies on the filesystem to index and access all data quickly.



### Setup

- Install go on your machine, prepare a configuration file.
- The configuration file contains all the parameters necessary to run the server.



### Performance

- Compared the performance of fasthttp on http 1.x against the performance of net/http running on http 2.0.



### License

- This project is distributed under an MIT license available on the project repository.
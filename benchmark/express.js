var express = require("express");
var upload = require("express-fileupload");
var http = require("http");
var fs = require("fs");

var app = express();
var dataPath = "data";

//If path doesnt exist create new path for data
if(!fs.existsSync(dataPath))
{
	fs.mkdirSync(dataPath);
}

//Server listen on port 80
http.Server(app).listen(9090); 

//Configure middleware
app.use(upload());

app.use(function(req, res, next)
{
	res.header("Access-Control-Allow-Origin", req.get("Origin") || "*");
	res.header("Access-Control-Allow-Credentials", "true");
	res.header("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE");
	res.header("Access-Control-Expose-Headers", "Content-Length");
	res.header("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Requested-With, Range");
	
	if(req.method === "OPTIONS")
	{
		return res.send(200);
	}
	else
	{
		return next();
	}
});


console.log("Server started.");

app.use("/v1/resource/get", express.static(dataPath));
app.post("/v1/resource/upload", function(req, res)
{
	var body = req.body;

	if(req.files.file)
	{
		var file = req.files.file;
		var path = __dirname + "/" + dataPath + "/" + body.library;

		if(!fs.existsSync(path))
		{
			fs.mkdirSync(path);
		}

		path += "/" + body.uuid + "." + body.format;

		file.mv(path, function(error)
		{
			if(error)
			{
				res.send("Error occured writting file");
				res.end();
			}
			else 
			{
				res.send("OK");
				res.end();
			}
		});
	}
	else
	{
		res.send("No file provided");
		res.end();
	}
});

const express = require("express");
const pathToSwaggerUi = require("swagger-ui-dist").getAbsoluteFSPath();
const path = require("path");
const fs = require("fs");

const port = 8888;
const baseUrl = `http://localhost:${port}`;

const SWAGGER_FILE_NAME = "swagger.yml";
const SWAGGER_FILE_PATH = path.join(__dirname, "static", SWAGGER_FILE_NAME);
const DEFAULT_CONFIG_FILE = "https://petstore.swagger.io/v2/swagger.json";
const SWAGGER_FILE_URL = `${baseUrl}/swagger.yml`;

// A workaround for swagger-ui-dist not being able to set custom swagger URL
const indexContent = fs
  .readFileSync(path.join(pathToSwaggerUi, "swagger-initializer.js"))
  .toString()
  .replace(DEFAULT_CONFIG_FILE, SWAGGER_FILE_URL);
const app = express();

app.get(`/swagger-initializer.js`, (_req, res) => {
  res.send(indexContent);
});

app.use(express.static(pathToSwaggerUi));

app.get("/", (_, res) => {
  res.sendFile(__dirname + "/index.html");
});

app.get("/swagger.yml", (_, res) => {
  res.sendFile(SWAGGER_FILE_PATH);
});

console.log(`Swagger UI is available at ${baseUrl}`);
app.listen(8888);

const path = require('path');
const express = require('express');
const app = express();
const fixture_dir = __dirname + '/test/fixtures/';

app.use(express.static('static'));

app.get('/api/todoist/configuration', (req, res) => {
  res.contentType('application/json');
  res.sendFile(path.normalize(fixture_dir + 'todoist_configuration.json'));
});

app.listen(3000, (err) => {
  if (err) {
    console.log(err);
  }
  console.log('server start at port 3000')
});

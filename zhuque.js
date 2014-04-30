var vash = require('./vash/build/vash.js');
var fs = require('fs');

fs.readFile("test.gohtml", 'utf8', function(err, data) {
  if (err) throw err;

  console.log(vash.compile(data, {debug:false, package: "tpl", name: "Home"}).toString());
});

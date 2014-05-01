var vash = require('./vash.js');
var fs = require('fs');
var path = require('path');

var gz_extension = ".gohtml";
var go_extension = ".go";

if (process.argv.length != 4) {
	console.log("Usage: \nnode gorazor.js template_folder_path output_path");
	return;
}

var tpl_folder = process.argv[2];
var output_folder = process.argv[3];

function get_names(tpl_path) {
	var parts = tpl_path.split(path.sep);
	var filename = parts[parts.length - 1];
	var names = {};

	names.package_name = parts[parts.length - 2];
	var func_name = filename.substr(0, filename.length - gz_extension.length);

	//Capitalize first character
	names.func_name = func_name.substr(0, 1).toUpperCase() + func_name.substr(1);
	return names;
}

function normalize_sep(folder_path) {
	if (folder_path[folder_path.length - 1] != path.sep) {
		folder_path = folder_path + path.sep;
	}

	return folder_path;
}

function process_file(tpl_path, output_path) {
	fs.readFile(tpl_path, 'utf8', function(err, data) {
		if (err) throw err;
		var options = {debug:false};

		var names = get_names(tpl_path);
		options["package"] = names.package_name;
		options["name"] = names.func_name;

		var code = vash.compile(data, options).toString();
		fs.writeFileSync(output_path, code);
	});
}

function process_folder(tpl_path, output_path) {
	tpl_path = normalize_sep(tpl_path);
	output_path = normalize_sep(output_path);

	fs.readdir(tpl_path, function(err, files) {
		if (err) throw err;
		for(var i=0; i < files.length; i++) {
			var input = tpl_path + files[i];
			var output = output_path + files[i].replace(gz_extension, go_extension);
			process_file(input, output);
		}
	});
}

process_folder(tpl_folder, output_folder);
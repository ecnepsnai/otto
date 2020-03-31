// jshint esversion:6
const fs = require('fs');
const path = require('path');

const pageName = process.argv[2];
const dirName = 'js/ng/pages/' + pageName;
const moduleName = 'otto';

var output = "angular.module('" + moduleName + "').component('" + pageName + "', {\n" +
    "    templateUrl: '/ottodev/static/html/" + pageName + ".html',\n" +
    "    bindings: {},\n" +
    "    controller: '" + pageName + "',\n" +
    "    controllerAs: ''\n" +
    "});\n";

fs.mkdirSync(dirName);
fs.writeFileSync(path.join(dirName, pageName) + 'Component.js', output);

output = "angular.module('" + moduleName + "').controller('" + pageName + "', function($scope) {\n" +
    "    var $ctrl = this;\n" +
    "});\n";

fs.writeFileSync(path.join(dirName, pageName) + 'Controller.js', output);
fs.writeFileSync(path.join(dirName, pageName) + '.html', '<page-wrapper page-title="' + pageName + '"></page-wrapper>');
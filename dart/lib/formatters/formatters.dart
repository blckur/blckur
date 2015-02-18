library formatters;

import 'package:blckur/formatters/capitalize.dart' as capitalize;

import 'package:angular/angular.dart' as ng;

class FormattersMod extends ng.Module {
  FormattersMod() {
    this.bind(capitalize.CapitalizeForm);
  }
}

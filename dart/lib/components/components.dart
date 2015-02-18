library components;

import 'package:blckur/components/test/test.dart'
  as test;

import 'package:angular/angular.dart' as ng;

class ComponentsMod extends ng.Module {
  ComponentsMod() {
    this.bind(test.TestComp);
  }
}

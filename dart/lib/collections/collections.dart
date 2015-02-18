library collections;

import 'package:blckur/collections/tests.dart' as tsts;

import 'package:angular/angular.dart' as ng;

class CollectionsMod extends ng.Module {
  CollectionsMod() {
    this.bind(tsts.Tests);
  }
}

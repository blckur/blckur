library models;

import 'package:blckur/models/user.dart' as usr;

import 'package:angular/angular.dart' as ng;

class ModelsMod extends ng.Module {
  ModelsMod() {
    this.bind(usr.User);
  }
}

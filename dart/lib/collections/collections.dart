library collections;

import 'package:blckur/collections/accounts.dart' as accts;
import 'package:blckur/collections/events.dart' as evnts;

import 'package:angular/angular.dart' as ng;

class CollectionsMod extends ng.Module {
  CollectionsMod() {
    this.bind(accts.Accounts);
    this.bind(evnts.Events);
  }
}

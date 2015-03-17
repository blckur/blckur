library collections;

import 'package:blckur/collections/account_types.dart' as account_types;
export 'package:blckur/collections/account_types.dart' show AccountTypes;
import 'package:blckur/collections/accounts.dart' as accounts;
export 'package:blckur/collections/accounts.dart' show Accounts;
import 'package:blckur/collections/alert_types.dart' as alert_types;
export 'package:blckur/collections/alert_types.dart' show AlertTypes;
import 'package:blckur/collections/events.dart' as events;
export 'package:blckur/collections/events.dart' show Events;

import 'package:angular/angular.dart' as ng;

class CollectionsMod extends ng.Module {
  CollectionsMod() {
    this.bind(account_types.AccountTypes);
    this.bind(accounts.Accounts);
    this.bind(alert_types.AlertTypes);
    this.bind(events.Events);
  }
}

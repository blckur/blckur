library models;

import 'package:blckur/models/account.dart' as account;
export 'package:blckur/models/account.dart' show Account;
import 'package:blckur/models/account_add.dart' as account_add;
export 'package:blckur/models/account_add.dart' show AccountAdd;
import 'package:blckur/models/account_type.dart' as account_type;
export 'package:blckur/models/account_type.dart' show AccountType;
import 'package:blckur/models/auth.dart' as auth;
export 'package:blckur/models/auth.dart' show Auth;
import 'package:blckur/models/event.dart' as event;
export 'package:blckur/models/event.dart' show Event;
import 'package:blckur/models/filter_type.dart' as filter_type;
export 'package:blckur/models/filter_type.dart' show FilterType;
import 'package:blckur/models/user.dart' as user;
export 'package:blckur/models/user.dart' show User;

import 'package:angular/angular.dart' as ng;

class ModelsMod extends ng.Module {
  ModelsMod() {
    this.bind(account.Account);
    this.bind(account_add.AccountAdd);
    this.bind(account_type.AccountType);
    this.bind(auth.Auth);
    this.bind(event.Event);
    this.bind(filter_type.FilterType);
    this.bind(user.User);
  }
}

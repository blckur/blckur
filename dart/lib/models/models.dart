library models;

import 'package:blckur/models/account.dart' as acct;
import 'package:blckur/models/account_add.dart' as acct_ad;
import 'package:blckur/models/account_type.dart' as acct_typ;
import 'package:blckur/models/auth.dart' as ath;
import 'package:blckur/models/event.dart' as evnt;
import 'package:blckur/models/user.dart' as usr;

import 'package:angular/angular.dart' as ng;

class ModelsMod extends ng.Module {
  ModelsMod() {
    this.bind(acct.Account);
    this.bind(acct_ad.AccountAdd);
    this.bind(acct_typ.AccountType);
    this.bind(ath.Auth);
    this.bind(evnt.Event);
    this.bind(usr.User);
  }
}

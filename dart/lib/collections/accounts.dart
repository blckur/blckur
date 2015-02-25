library accounts_col;

import 'package:blckur/collection.dart' as collec;
import 'package:blckur/models/account.dart' as acct;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class Accounts extends collec.Collection {
  Type model = acct.Account;

  Accounts(ng.Http http) : super(http);

  String url = '/accounts';
}

library account_types_col;

import 'package:blckur/collection.dart' as collec;
import 'package:blckur/models/account_type.dart' as acct_typ;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class AccountTypes extends collec.Collection {
  Type model = acct_typ.AccountType;

  AccountTypes(ng.Http http) : super(http);

  String url = 'http://localhost:3000/account_types';
}

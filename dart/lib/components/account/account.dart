library account_comp;

import 'package:blckur/models/account.dart' as acct;

import 'package:angular/angular.dart' show Component, NgOneWay;

@Component(
  selector: 'x-account',
  templateUrl: 'packages/blckur/components/account/account.html',
  cssUrl: 'packages/blckur/components/account/account.css'
)
class AccountComp {
  @NgOneWay('model')
  acct.Account model;
}

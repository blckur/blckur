library account_comp;

import 'package:blckur/models/account.dart' as acct;

import 'package:angular/angular.dart' show Component, NgOneWay;

// Add attr new if new highlight account for init setup

@Component(
  selector: 'x-account',
  templateUrl: 'packages/blckur/components/account/account.html',
  cssUrl: 'packages/blckur/components/account/account.css'
)
class AccountComp {
  bool settings;
  bool addingAlert;

  @NgOneWay('model')
  acct.Account model;

  void onSettings() {
    this.settings = this.settings != true;
  }
}

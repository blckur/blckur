library account_add_comp;

import 'package:blckur/collections/account_types.dart' as acct_typs;
import 'package:blckur/logger.dart' as logger;

import 'package:angular/angular.dart' show Component, NgTwoWay;

@Component(
  selector: 'x-account-add',
  templateUrl: 'packages/blckur/components/account_add/account_add.html',
  cssUrl: 'packages/blckur/components/account_add/account_add.css'
)
class AccountAddComp {
  acct_typs.AccountTypes accountTypes;

  @NgTwoWay('add-account')
  bool addAccount;

  AccountAddComp(this.accountTypes) {
    this.accountTypes.fetch().catchError((err) {
      logger.severe('Failed to load account types', err);
    });
  }

  void onAccountAdd() {
    this.addAccount = true;
  }
}

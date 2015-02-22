library account_add_comp;

import 'package:angular/angular.dart' show Component, NgTwoWay;

@Component(
  selector: 'x-account-add',
  templateUrl: 'packages/blckur/components/account_add/account_add.html',
  cssUrl: 'packages/blckur/components/account_add/account_add.css'
)
class AccountAddComp {
  @NgTwoWay('add-account')
  bool addAccount;

  void onAccountAdd() {
    this.addAccount = true;
  }
}

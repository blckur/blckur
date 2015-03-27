library accounts_comp;

import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/loading.dart' as loading;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgOneWay;

@Component(
  selector: 'x-accounts',
  templateUrl: 'packages/blckur/components/accounts/accounts.html',
  cssUrl: 'packages/blckur/components/accounts/accounts.css'
)
class AccountsComp extends loading.Loading {
  bool addAccount;
  collections.Accounts accounts;

  AccountsComp(this.accounts) {
    this.update();
  }

  void update() {
    this.accounts.fetch().catchError((err) {
      logger.severe('Failed to load accounts', err);
      new alert.Alert('Error loading accounts', () {
        this.update();
      });
    });
  }
}

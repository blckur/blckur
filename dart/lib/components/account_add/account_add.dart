library account_add_comp;

import 'package:blckur/models/account_add.dart' as acct_ad;
import 'package:blckur/collections/account_types.dart' as acct_typs;
import 'package:blckur/loading.dart' as lodin;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alrt;

import 'package:angular/angular.dart' show Component, NgTwoWay;
import 'dart:html' as dom;

@Component(
  selector: 'x-account-add',
  templateUrl: 'packages/blckur/components/account_add/account_add.html',
  cssUrl: 'packages/blckur/components/account_add/account_add.css'
)
class AccountAddComp extends lodin.Loading {
  acct_ad.AccountAdd model;
  acct_typs.AccountTypes accountTypes;

  @NgTwoWay('add-account')
  bool addAccount;

  AccountAddComp(this.model, this.accountTypes) {
    this.accountTypes.fetch().catchError((err) {
      logger.severe('Failed to load account types', err);
    });
  }

  void onAccountAdd() {
    this.addAccount = true;
  }

  void onAccountCancel() {
    this.addAccount = false;
    this.clearLoading();
  }

  void onAccountClick(String type) {
    if (!this.setLoading()) {
      return;
    }

    this.model.type = type;
    this.model.create().then((_) {
      if (this.addAccount != true) {
        return;
      }
      dom.window.location.replace(this.model.redirect);
    }).catchError((err) {
      if (this.addAccount != true) {
        return;
      }
      logger.severe('Failed to add account', err);
      new alrt.Alert('Unable to add account, try again later.', () {
        this.onAccountClick(type);
      });
    }).whenComplete(() {
      this.clearLoading();
    });
  }
}

library account_comp;

import 'package:blckur/models/account.dart' as acct;
import 'package:blckur/loading.dart' as lodin;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alrt;

import 'package:angular/angular.dart' show Component, NgOneWay;
import 'dart:html' as dom;

// Add attr new if new highlight account for init setup

@Component(
  selector: 'x-account',
  templateUrl: 'packages/blckur/components/account/account.html',
  cssUrl: 'packages/blckur/components/account/account.css'
)
class AccountComp extends lodin.Loading {
  bool settings;
  bool addingAlert;
  bool confirm;
  dom.ShadowRoot root;

  AccountComp(this.root);

  @NgOneWay('model')
  acct.Account model;

  void onSettings() {
    this.settings = this.settings != true;
  }

  void onDelete() {
    if (!this.setLoading()) {
      return;
    }

    this.model.destroy().catchError((err) {
      logger.severe('Failed to delete account', err);
      new alrt.Alert('Failed to delete account');
    }).whenComplete(() {
      this.root.host.remove();
    });
  }
}

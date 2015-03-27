library account_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/injectables/injectables.dart' as injectables;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgOneWay;
import 'dart:html' as dom;

// Add attr new if new highlight account for init setup

@Component(
  selector: 'x-account',
  templateUrl: 'packages/blckur/components/account/account.html',
  cssUrl: 'packages/blckur/components/account/account.css'
)
class AccountComp {
  bool settings;
  bool addingAlert;
  bool confirm;
  dom.ShadowRoot root;
  injectables.Loading loading;

  AccountComp(this.root, this.loading);

  @NgOneWay('model')
  models.Account model;

  void onSettings() {
    this.settings = this.settings != true;
  }

  void onDelete() {
    if (!this.loading.set()) {
      return;
    }

    this.model.destroy().catchError((err) {
      logger.severe('Failed to delete account', err);
      new alert.Alert('Failed to delete account');
    }).whenComplete(() {
      this.root.host.remove();
    });
  }
}

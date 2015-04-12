library account_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/utils/utils.dart' as utils;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgOneWay;
import 'dart:html' as dom;

@Component(
  selector: 'x-account',
  templateUrl: 'packages/blckur/components/account/account.html',
  cssUrl: 'packages/blckur/components/account/account.css'
)
class AccountComp {
  bool settings;
  bool addingAlert;
  bool confirm;
  utils.Loading loading;

  AccountComp(this.root) {
    this.loading = new utils.Loading();
  }

  @NgOneWay('model')
  models.Account model;

  String get itemClass {
    var cls = '';

    if (this.settings == true) {
      cls += 'settings-active';
    }

    if (this.model.newAcct == true) {
      cls += ' new';
    }

    return cls;
  }

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
      this.loading.clear();
    });
  }
}

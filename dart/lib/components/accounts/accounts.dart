library accounts_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/utils/utils.dart' as utils;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgOneWay;
import 'dart:html' as dom;

const NONE = 0;
const ACCOUNTS = 1;
const ADD_ACCOUNT = 2;

@Component(
  selector: 'x-accounts',
  templateUrl: 'packages/blckur/components/accounts/accounts.html',
  cssUrl: 'packages/blckur/components/accounts/accounts.css'
)
class AccountsComp {
  models.AccountAdd model;
  collections.Accounts accounts;
  collections.AccountTypes accountTypes;
  utils.Loading loading;
  utils.ModeSwitch modeswitch;
  Map<String, String> typeClasses;

  AccountsComp() {
    this.model = new models.AccountAdd();
    this.accounts = new collections.Accounts();
    this.accountTypes = new collections.AccountTypes();
    this.loading = new utils.Loading();
    this.modeswitch = new utils.ModeSwitch();
    this.modeswitch.mode = ACCOUNTS;
    this.update(true);
  }

  void update([bool initMode]) {
    if (!this.loading.set()) {
      return;
    }

    this.accounts.fetch().catchError((err) {
      logger.severe('Failed to load accounts', err);
      new alert.Alert('Error loading accounts', () {
        this.update();
      });
    }).whenComplete(() {
      if (initMode == true) {
        this.modeswitch.mode = ACCOUNTS;
      }
      this.loading.clear();
    });
  }

  void onAdd() {
    if (!this.loading.set()) {
      return;
    }
    this.modeswitch.mode = NONE;
    this.typeClasses = {};

    this.accountTypes.fetch().catchError((err) {
      logger.severe('Failed to load account types', err);
      new alert.Alert('Failed to load account types');
    }).whenComplete(() {
      this.loading.clear();
      this.modeswitch.mode = ADD_ACCOUNT;
    });
  }

  void onCancel() {
    this.modeswitch.mode = ACCOUNTS;
    this.loading.clear();
  }

  void onClick(String type) {
    if (!this.loading.set()) {
      return;
    }

    this.typeClasses = {
      type: 'active',
    };

    this.model.type = type;
    this.model.create().then((_) {
      if (this.modeswitch.mode != ADD_ACCOUNT) {
        return;
      }
      dom.window.location.replace(this.model.redirect);
    }).catchError((err) {
      if (this.modeswitch.mode != ADD_ACCOUNT) {
        return;
      }
      logger.severe('Failed to add account', err);
      new alert.Alert('Unable to add account, try again later.', () {
        this.onClick(type);
      });
      this.loading.clear();
    });
  }
}

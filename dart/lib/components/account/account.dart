library account_comp;

import 'package:blckur/models/account.dart' as acct;
import 'package:blckur/alert.dart' as alrt;
import 'package:blckur/logger.dart' as logger;

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
  acct.Account settingsModel;

  @NgOneWay('model')
  acct.Account model;

  void onSettings() {
    this.settingsModel = this.model.clone();
    this.settings = this.settings != true;
  }

  void onSave() {
    this.settingsModel.save(['events']).then((_) {
      this.settings = false;
      this.model = this.settingsModel;
    }).catchError((err) {
      logger.severe('Failed to save account settings', err);
      new alrt.Alert('Failed to save account settings');
    });
  }
}

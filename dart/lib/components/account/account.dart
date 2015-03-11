library account_comp;

import 'package:blckur/models/account.dart' as acct;
import 'package:blckur/loading.dart' as lodin;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alrt;

import 'package:angular/angular.dart' show Component, NgOneWay;

// Add attr new if new highlight account for init setup

@Component(
  selector: 'x-account',
  templateUrl: 'packages/blckur/components/account/account.html',
  cssUrl: 'packages/blckur/components/account/account.css'
)
class AccountComp extends lodin.Loading {
  bool settings;
  bool addingAlert;

  @NgOneWay('model')
  acct.Account model;

  void onSettings() {
    this.settings = this.settings != true;
  }

  void onAlertDel(Map<String, String> alert) {
    if (!this.setLoading()) {
      return;
    }

    this.model.alerts.remove(alert);
    this.model.save(['alerts']).catchError((err) {
      logger.severe('Failed to remove alert', err);
      new alrt.Alert('Failed to remove alert');
    }).whenComplete(() {
      this.clearLoading();
    });
  }
}

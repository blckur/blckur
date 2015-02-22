library account_comp;

import 'package:blckur/models/user.dart' as usr;
import 'package:blckur/logger.dart' as logger;

import 'package:angular/angular.dart' show Component;

@Component(
  selector: 'x-account',
  templateUrl: 'packages/blckur/components/account/account.html',
  cssUrl: 'packages/blckur/components/account/account.css'
)
class AccountComp {
  usr.User model;

  AccountComp(this.model) {
    this.update();
  }

  void update() {
    this.model.fetch().catchError((err) {
      logger.severe('Failed to load organizations', err);
    });
  }
}

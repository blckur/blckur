library feed_comp;

import 'package:blckur/collections/accounts.dart' as accts;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alrt;

import 'package:angular/angular.dart' show Component;

@Component(
  selector: 'x-feed',
  templateUrl: 'packages/blckur/components/feed/feed.html',
  cssUrl: 'packages/blckur/components/feed/feed.css'
)
class FeedComp {
  bool addAccount;
  accts.Accounts accounts;

  FeedComp(this.accounts) {
    this.update();
  }

  void update() {
    this.accounts.fetch().catchError((err) {
      logger.severe('Failed to load accounts', err);
      new alrt.Alert('Error loading accounts', () {
        this.update();
      });
    });
  }
}

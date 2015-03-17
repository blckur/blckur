library feed_comp;

import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component;

@Component(
  selector: 'x-feed',
  templateUrl: 'packages/blckur/components/feed/feed.html',
  cssUrl: 'packages/blckur/components/feed/feed.css'
)
class FeedComp {
  bool addAccount;
  collections.Accounts accounts;

  FeedComp(this.accounts) {
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

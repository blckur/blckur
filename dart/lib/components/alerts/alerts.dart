library alerts_comp;

import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component;
import 'dart:collection' as collection;

@Component(
  selector: 'x-alerts',
  template: '<paper-toast ng-repeat="alert in alerts" py-opened="alert.opened"'
    'py-text="alert.text" duration="{{alert.duration}}" '
    'autoCloseDisabled="true" ng-class="getClass(\$index)" no-select><div '
    'class="retry" ng-if="alert.retryCallback != null" '
    'ng-click="alert.retryCallback()">Retry</div></paper-toast>',
  cssUrl: 'packages/blckur/components/alerts/alerts.css'
)
class AlertsComp {
  collection.Queue<alert.Alert> alerts = alert.alerts;

  String getClass(int index) {
    return 'alert${index}';
  }
}

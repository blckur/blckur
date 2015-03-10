library alert_add_comp;

import 'package:blckur/collections/alert_types.dart' as alrt_typs;
import 'package:blckur/loading.dart' as lodin;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alrt;

import 'package:angular/angular.dart' show Component, NgTwoWay;

@Component(
  selector: 'x-alert-add',
  templateUrl: 'packages/blckur/components/alert_add/alert_add.html',
  cssUrl: 'packages/blckur/components/alert_add/alert_add.css'
)
class AlertAddComp extends lodin.Loading {
  alrt_typs.AlertTypes alertTypes;

  @NgTwoWay('active')
  bool active;

  @NgTwoWay('alerts')
  List<Map<String, String>> alerts;

  AlertAddComp(this.alertTypes);

  void onAdd() {
    if (!this.setLoading()) {
      return;
    }

    this.active = true;
    this.alertTypes.fetch().catchError((err) {
      logger.severe('Failed to load alert types', err);
      new alrt.Alert('Failed to load alert types');
    }).whenComplete(() {
      this.clearLoading();
    });
  }

  void onCancel() {
    this.active = false;
    this.clearLoading();
  }
}

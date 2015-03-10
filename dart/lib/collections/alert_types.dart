library alert_types_col;

import 'package:blckur/collection.dart' as collec;
import 'package:blckur/models/alert_type.dart' as alrt_typ;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class AlertTypes extends collec.Collection {
  Type model = alrt_typ.AlertType;

  AlertTypes(ng.Http http) : super(http);

  String url = '/alert_types';
}

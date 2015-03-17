library alert_types_col;

import 'package:blckur/collection.dart' as collection;
import 'package:blckur/models/models.dart' as models;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class AlertTypes extends collection.Collection {
  Type model = models.AlertType;

  AlertTypes(ng.Http http) : super(http);

  String url = '/alert_types';
}

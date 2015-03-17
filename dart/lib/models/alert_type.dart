library alert_type_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class AlertType extends model.Model {
  @model.Attribute('id')
  int id;

  @model.Attribute('label')
  String label;

  @model.Attribute('type')
  String type;

  @model.Attribute('value_type')
  String valueType;

  @model.Attribute('value_label')
  String valueLabel;

  @model.Attribute('value_holder')
  String valueHolder;

  AlertType(ng.Http http) : super(http);
}

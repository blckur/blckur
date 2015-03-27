library alert_types_col;

import 'package:blckur/model.dart' as model;
import 'package:blckur/collection.dart' as collection;
import 'package:blckur/models/models.dart' as models;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class AlertTypes extends collection.Collection {
  String acctType;

  AlertTypes(ng.Http http) : super(http);

  model.Model newModel() {
    return new models.AlertType(this.http);
  }

  collection.Collection newCollection() {
    return new AlertTypes(this.http);
  }

  String get url {
    return '/alert_types/${this.acctType}';
  }
}

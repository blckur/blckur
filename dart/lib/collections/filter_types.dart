library filter_types_col;

import 'package:blckur/model.dart' as model;
import 'package:blckur/collection.dart' as collection;
import 'package:blckur/models/models.dart' as models;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class FilterTypes extends collection.Collection {
  String acctType;

  FilterTypes(ng.Http http) : super(http);

  model.Model newModel() {
    return new models.FilterType(this.http);
  }

  collection.Collection newCollection() {
    return new FilterTypes(this.http);
  }

  String get url {
    return '/filter_types/${this.acctType}';
  }
}

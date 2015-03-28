library filter_types_col;

import 'package:blckur/model.dart' as model;
import 'package:blckur/collection.dart' as collection;
import 'package:blckur/models/models.dart' as models;

class FilterTypes extends collection.Collection {
  String acctType;

  model.Model newModel() {
    return new models.FilterType();
  }

  collection.Collection newCollection() {
    return new FilterTypes();
  }

  String get url {
    return '/filter_types/${this.acctType}';
  }
}

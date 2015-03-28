library accounts_col;

import 'package:blckur/model.dart' as model;
import 'package:blckur/collection.dart' as collection;
import 'package:blckur/models/models.dart' as models;

class Accounts extends collection.Collection {
  model.Model newModel() {
    return new models.Account();
  }

  collection.Collection newCollection() {
    return new Accounts();
  }

  String url = '/accounts';
}

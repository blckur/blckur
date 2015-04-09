library notifications_col;

import 'package:blckur/model.dart' as model;
import 'package:blckur/collection.dart' as collection;
import 'package:blckur/models/models.dart' as models;

class Notifications extends collection.Collection {
  model.Model newModel() {
    return new models.Notification();
  }

  collection.Collection newCollection() {
    return new Notifications();
  }

  String url = '/notifications';
}

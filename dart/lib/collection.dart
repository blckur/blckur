library collection;

import 'package:blckur/exceptions.dart';
import 'package:blckur/remote.dart' as remote;
import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' as ng;
import 'dart:collection' as collection;

abstract class Collection extends remote.Remote with collection.IterableMixin {
  List<model.Model> _collection;
  Function onAdd;
  Function onChange;
  Function onRemove;

  Collection() : super(), _collection = [];

  model.Model newModel() {
    throw new UnimplementedError('Model new not implemented.');
  }

  Collection newCollection() {
    throw new UnimplementedError('Collection new not implemented.');
  }

  Iterator get iterator {
    return this._collection.iterator;
  }

  dynamic operator [](int index) {
    return this._collection[index];
  }

  void add(Map<String, dynamic> attrs) {
    var mdl = this.newModel();
    mdl.import(attrs);
    this._collection.add(mdl);
  }

  void validate(String name) {
    for (var model in this) {
      model.validate(name);
    }
  }

  Collection clone() {
    var clone = this.newCollection();

    for (var mdl in this) {
      clone._collection.add(mdl.clone());
    }

    return clone;
  }

  void imported() {}

  void added(model.Model mdl) {}

  void changed(model.Model mdl) {}

  void removed(model.Model mdl) {}

  void import(dynamic responseData) {
    var data;

    try {
      data = this.parse(responseData);
    } on IgnoreResponse {
      return;
    }

    var curIds = [];
    var newIds = new Set();
    var models = {};
    var recModels = new collection.Queue();
    var coll = [];

    for (var model in this._collection) {
      curIds.add(model.id);
      models[model.id] = model;
    }

    for (var i = 0; i < data.length; i++) {
      newIds.add(data[i]['id']);
    }

    curIds.forEach((id) {
      if (!newIds.contains(id)) {
        var model = models.remove(id);

        this.removed(model);
        if (this.onRemove != null) {
          this.onRemove(model);
        }

        recModels.add(model);
      }
    });

    for (var i = 0; i < data.length; i++) {
      var added;
      var model = models[data[i]['id']];

      if (model == null) {
        added = true;
        if (recModels.length > 0) {
          model = recModels.removeFirst();
        }
        else {
          model = this.newModel();
        }
        model.init();
      }

      model.import(data[i]);
      coll.add(model);

      if (added == true) {
        this.added(model);
        if (this.onAdd != null) {
          this.onAdd(model);
        }
      }
    }

    this._collection = coll;

    this.imported();
    if (this.onImport != null) {
      this.onImport();
    }
  }

  void clear() {
    this._collection = [];
  }
}

library model;

import 'package:blckur/remote.dart' as remote;

import 'package:angular/angular.dart' as ng;
import 'dart:async' as async;

Map<Type, Map<String, Symbol>> _attrSymbols = {};
Map<Type, Map<String, Function>> _attrValidators = {};

class Attribute {
  final String name;
  const Attribute(this.name);
}

class Validator {
  final String name;
  const Validator(this.name);
}

class Invalid extends Error {
  String type;
  String message;

  Invalid(this.type, this.message);

  toString() => this.message;
}

abstract class Model extends remote.Remote {
  dynamic id;
  Function onLinkClear;

  Map<String, Function> get mapSet {
    throw new UnimplementedError('Setter map not implemented.');
  }
  Map<String, Function> get mapGet {
    throw new UnimplementedError('Getter map not implemented.');
  }
  Map<String, Function> get validators {
    return {};
  }
  Model modelNew() {
    throw new UnimplementedError('Model new not implemented.');
  }

  Model(ng.Http http) : super(http) {
    this.init();
  }

  void validate(String name) {
    var validator = this.validators[name];

    if (validator != null) {
      validator(this.mapGet[name]());
    }
  }

  Model clone() {
    var mdl = this.modelNew();
    var mapGet = this.mapGet;
    var mapSet = mdl.mapSet;

    mapSet.forEach((name, setter) {
      setter(mapGet[name]());
    });

    return mdl;
  }

  void import(dynamic responseData) {
    var data = this.parse(responseData);
    var mapSet = this.mapSet;

    if (data != null && data != '') {
      data.forEach((key, value) {
        var setter = mapSet[key];

        if (setter != null) {
          setter(value);
        }
      });
    }

    this.imported();
    if (this.onImport != null) {
      this.onImport();
    }
  }

  Map<String, dynamic> export([List<String> fields]) {
    var data = {};
    var mapGet = this.mapGet;

    mapGet.forEach((name, getter) {
      data[name] = getter();
    });

    return data;
  }

  async.Future destroy() {
    var loadId = this.setLoading();

    return this.http.delete(this.url).then((response) {
      this.clearLoading(loadId);
      this.clearError();
      this.import(response.data);
      return response.data;
    }).catchError((err) {
      this.clearLoading(loadId);
      return new async.Future.error(this.parseError(err));
    }, test: (e) => e is ng.HttpResponse);
  }

  async.Future send(String method, String url, List<String> fields) {
    var loadId = this.setLoading();
    var methodFunc;

    var data = this.export(fields);

    if (method == 'post') {
      methodFunc = this.http.post;
    }
    else if (method == 'put') {
      methodFunc = this.http.put;
    }
    else {
      throw new ArgumentError('Unkown method');
    }

    return methodFunc(url, data).then((response) {
      this.clearLoading(loadId);
      this.clearError();
      this.import(response.data);
      return response.data;
    }).catchError((err) {
      this.clearLoading(loadId);
      return new async.Future.error(this.parseError(err));
    }, test: (e) => e is ng.HttpResponse);
  }

  async.Future save([List<String> fields]) {
    return this.send('put', this.url, fields);
  }

  async.Future create([List<String> fields]) {
    return this.send('post', this.url, fields);
  }

  void clear() {
    var mapSet = this.mapSet;

    mapSet.forEach((_, setter) {
      setter(null);
    });
  }

  void init() {}
}
